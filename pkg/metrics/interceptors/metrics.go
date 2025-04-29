package interceptors

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var (
	serverTotalRequestsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "server_total_requests",
			Help: "Total number of requests to the server",
		},
		[]string{"method", "status"},
	)

	serverDurationRequestsHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "server_duration_requests_seconds",
			Help:    "Duration of requests to the server in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	serverRequestsErrorsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "server_request_errors",
			Help: "Number of requests with error to the server",
		},
		[]string{"method", "status"},
	)

	serverActiveRequestsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "server_active_requests",
			Help: "Number of active requests to the server",
		},
		[]string{"method"},
	)

	serverTotalBytesCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "server_total_bytes",
			Help: "Total bytes of requests and responses to the server",
		},
		[]string{"direction"},
	)
)

func init() {
	prometheus.MustRegister(
		serverTotalRequestsCounter,
		serverDurationRequestsHistogram,
		serverRequestsErrorsCounter,
		serverActiveRequestsGauge,
		serverTotalBytesCounter,
	)
}

func ServerMetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		defer prometheus.NewTimer(serverDurationRequestsHistogram.WithLabelValues(info.FullMethod)).ObserveDuration()

		serverActiveRequestsGauge.WithLabelValues(info.FullMethod).Inc()
		defer serverActiveRequestsGauge.WithLabelValues(info.FullMethod).Dec()

		res, err := handler(ctx, req)
		code := status.Code(err)

		serverTotalRequestsCounter.WithLabelValues(info.FullMethod, code.String()).Inc()

		if err != nil {
			serverRequestsErrorsCounter.WithLabelValues(info.FullMethod, code.String()).Inc()
		}

		if m, ok := req.(proto.Message); ok {
			serverTotalBytesCounter.WithLabelValues("in").Add(float64(proto.Size(m)))
		}

		if m, ok := res.(proto.Message); ok {
			serverTotalBytesCounter.WithLabelValues("out").Add(float64(proto.Size(m)))
		}

		return res, err
	}
}
