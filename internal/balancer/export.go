package balancer

import (
	"context"
	"fmt"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

func export(project string, location string, cluster string, namespace string, hpa string, value int) error {
	value64 := int64(value)
	ctx := context.Background()
	client, err := monitoring.NewMetricClient(ctx)
	dataPoint := &monitoringpb.Point{
		Interval: &monitoringpb.TimeInterval{
			EndTime: &googlepb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
		Value: &monitoringpb.TypedValue{
			Value: &monitoringpb.TypedValue_Int64Value{
				Int64Value: value64,
			},
		},
	}
	err = client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
		Name: monitoring.MetricProjectPath(project),
		TimeSeries: []*monitoringpb.TimeSeries{
			{
				Metric: &metricpb.Metric{
					Type: "custom.googleapis.com/balancer",
				},
				Resource: &monitoredrespb.MonitoredResource{
					Type: "generic_node",
					Labels: map[string]string{
						"project_id": project,
						"location":   location,
						"namespace":  cluster,
						"node_id":    fmt.Sprintf("%s-%s", namespace, hpa),
					},
				},
				Points: []*monitoringpb.Point{
					dataPoint,
				},
			},
		},
	})
	if err != nil {
		return err
	}
	err = client.Close()
	if err != nil {
		return err
	}
	return nil
}
