package metrics

import "github.com/prometheus/client_golang/prometheus"

// HTTP 요청 총 수
var HttpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"method", "path", "status"},
)

// HTTP 요청 응답 시간
var HttpRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "path"},
)

// 활성 요청 수
var HttpRequestsInFlight = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Number of HTTP requests currently being processed",
	},
)

// 주문 생성 수
var OrdersCreatedTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "orders_created_total",
		Help: "Total number of orders created",
	},
	[]string{"status"},
)

// 주문 처리 시간
var OrderProcessingDuration = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "order_processing_duration_seconds",
		Help:    "Time spent processing an order",
		Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
	},
)

func init() {
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpRequestDuration)
	prometheus.MustRegister(HttpRequestsInFlight)
	prometheus.MustRegister(OrdersCreatedTotal)
	prometheus.MustRegister(OrderProcessingDuration)
}
