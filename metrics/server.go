package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// BackendCreateOpsTotal - The total number of creating backend operations
	BackendCreateOpsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "sshare",
		Subsystem: "driver",
		Name:      "backend_create_ops_total",
		Help:      "The total number of creating backend operations",
	})

	// BackendCreateSuccessOpsTotal - The total number of successfully created
	// backends with breakdown for component
	BackendCreateSuccessOpsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sshare",
			Subsystem: "driver",
			Name:      "backend_create_success_ops_total",
			Help:      "The total number of successfully created backends with breakdown for component",
		},
		[]string{"component"},
	)

	// BackendCreateErrorOpsTotal - The total number of creating backend operations
	// finished with error and with breakdown for component
	BackendCreateErrorOpsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sshare",
			Subsystem: "driver",
			Name:      "backend_create_error_ops_total",
			Help:      "The total number of creating backend operations finished with error and with breakdown for component",
		},
		[]string{"component"},
	)

	// BackendDeleteOpsTotal - The total number of deleting backend operations
	BackendDeleteOpsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "sshare",
		Subsystem: "driver",
		Name:      "backend_delete_ops_total",
		Help:      "The total number of deleting backend operations",
	})

	// BackendDeleteSuccessOpsTotal - The total number of successfully deleted
	// backends with breakdown for component
	BackendDeleteSuccessOpsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sshare",
			Subsystem: "driver",
			Name:      "backend_delete_success_ops_total",
			Help:      "The total number of successfully deleted backends with breakdown for component",
		},
		[]string{"component"},
	)

	// BackendDeleteErrorOpsTotal - The total number of deleting backend operations
	// finished with error and with breakdown for component
	BackendDeleteErrorOpsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sshare",
			Subsystem: "driver",
			Name:      "backend_delete_error_ops_total",
			Help:      "The total number of deleting backend operations finished with error and with breakdown for component",
		},
		[]string{"component"},
	)

	// BackendReadyTotal - The total number of backends reported as ready with breakdown for component
	BackendReadyTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sshare",
			Subsystem: "driver",
			Name:      "backend_ready_total",
			Help:      "The total number of backends reported as ready with breakdown for component",
		},
		[]string{"component", "status"},
	)

	// BackendReadyErrorTotal - The total number of backend reported as not ready
	// and finished with error, with breakdown for component
	BackendReadyErrorTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sshare",
			Subsystem: "driver",
			Name:      "backend_ready_error_total",
			Help:      "The total number of backend reported as not ready and finished with error, with breakdown for component",
		},
		[]string{"component"},
	)

	// BackendReadyTimeoutTotal - The total number of backend reported as not ready because of timeout
	BackendReadyTimeoutTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "sshare",
			Subsystem: "driver",
			Name:      "backend_ready_timeout_total",
			Help:      "The total number of backend reported as not ready because of timeout",
		},
	)

	// BackendIsReadyDuration - Histogram for the runtime of the IsReady function
	BackendIsReadyDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "sshare",
			Subsystem: "driver",
			Name:      "is_ready_duration_seconds",
			Help:      "Histogram for the runtime of the IsReady function",
			Buckets:   []float64{1, 2, 5, 10, 20, 60, 120},
		},
	)

	// AuthenticationTotal - The total number of authentications
	AuthenticationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sshare",
			Subsystem: "server",
			Name:      "authentication_total",
			Help:      "The total number of authentications",
		},
		[]string{"type", "interceptor"},
	)
)
