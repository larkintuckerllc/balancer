package balancer

import (
	"context"
	"time"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type ClusterEvent struct {
	cluster string
	event   watch.Event
}

func Execute(project string, location string, prefix string, clusters []string, namespace string, hpa string, value int, idle int) error {
	ctx := context.Background()
	replicas := map[string]int32{}
	updater := map[string]<-chan watch.Event{}
	for _, cluster := range clusters {
		clientset, err := k8sAuth(prefix, cluster)
		if err != nil {
			return err
		}
		api := clientset.AutoscalingV2beta2().HorizontalPodAutoscalers(namespace)
		reps, resourceVersion, err := initialReplicas(ctx, api, hpa)
		if err != nil {
			return err
		}
		replicas[cluster] = reps
		listOptions := metaV1.ListOptions{
			ResourceVersion: resourceVersion,
		}
		watcher, err := api.Watch(ctx, listOptions)
		if err != nil {
			return err
		}
		defer watcher.Stop()
		uptr := watcher.ResultChan()
		updater[cluster] = uptr
	}
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	// AGGREGATE CHANNEL
	eventer := make(chan ClusterEvent)
	for _, cluster := range clusters {
		go func(cl string, c <-chan watch.Event) {
			for event := range c {
				clusterEvent := ClusterEvent{cl, event}
				eventer <- clusterEvent
			}
		}(cluster, updater[cluster])
	}

	for {
		select {
		case clusterEvent := <-eventer:
			reps := updateReplicas(clusterEvent.event, hpa, replicas[clusterEvent.cluster])
			replicas[clusterEvent.cluster] = reps
		case <-ticker.C:
			state, scaleClusters := logic(replicas, idle)
			for _, cluster := range clusters {
				val := 0
				if state == Scaling && find(scaleClusters, cluster) {
					val = value
				}
				err := export(project, location, cluster, namespace, hpa, val)
				if err != nil {
					return err
				}
			}
		}
	}
}
