package balancer

import (
	"context"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/autoscaling/v2beta2"
)

func initialReplicas(ctx context.Context, api v2beta2.HorizontalPodAutoscalerInterface, hpa string) (int32, string, error) {
	getOptions := metaV1.GetOptions{}
	hpaObject, err := api.Get(ctx, hpa, getOptions)
	if err != nil {
		return 0, "", err
	}
	replicas := hpaObject.Status.CurrentReplicas
	hpaObjects, err := api.List(ctx, metaV1.ListOptions{})
	if err != nil {
		return 0, "", err
	}
	resourceVersion := hpaObjects.ListMeta.ResourceVersion
	return replicas, resourceVersion, nil
}
