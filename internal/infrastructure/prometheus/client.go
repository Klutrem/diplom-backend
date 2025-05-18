package prometheus

import (
	"context"
	"main/internal/config"
	"main/pkg"
	"time"

	prometheusApi "github.com/prometheus/client_golang/api"
	prometheusV1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"go.uber.org/fx"
)

type PrometheusClient struct {
	logger pkg.Logger
	client prometheusApi.Client
	api    prometheusV1.API
}

var Module = fx.Module("prometheus",
	fx.Provide(NewPrometheusClient),
)

func NewPrometheusClient(logger pkg.Logger, env config.Env) PrometheusClient {
	client, err := prometheusApi.NewClient(prometheusApi.Config{
		Address: env.PrometheusHost,
	})
	if err != nil {
		logger.Fatalf("failed to create prometheus client: %v", err)
	}
	api := prometheusV1.NewAPI(client)
	return PrometheusClient{
		logger: logger,
		client: client,
		api:    api,
	}
}

func (client PrometheusClient) GetMetricValue(query string) (model.Value, error) {
	value, warnings, err := client.api.Query(context.Background(), query, time.Now())
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		client.logger.Warnf("Prometheus query warnings: %v", warnings)
	}
	return value, nil
}
