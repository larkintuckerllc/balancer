package balancer

import (
	"k8s.io/api/autoscaling/v2beta2"
	"k8s.io/apimachinery/pkg/watch"
)

func updateReplicas(update watch.Event, hpa string, prevReplicas int32) int32 {
	hpaObject, ok := update.Object.(*v2beta2.HorizontalPodAutoscaler)
	if !ok {
		return prevReplicas
	}
	if hpaObject.ObjectMeta.Name != hpa {
		return prevReplicas
	}
	replicas := hpaObject.Status.CurrentReplicas
	return replicas
}
