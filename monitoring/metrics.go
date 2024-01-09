package monitoring

import (
	"github.com/Alluxio/k8s-operator/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

// MetricDescription is an exported struct that defines the metric description (Name, Help)
// as a new type named MetricDescription.
type MetricDescription struct {
	Name string
	Help string
	Type string
}

// metricDescription is a map of string keys (metrics) to MetricDescription values (Name, Help).
var metricDescription = map[string]MetricDescription{
	"AlluxioClusterAliveWorkerTotal": {
		Name: "alluxio_cluster_alive_worker_total",
		Help: "Total number of alive worker in Alluxio Cluster.",
		Type: "Gauge",
	},
	"AlluxioClusterDatasetMountedCountTotal": {
		Name: "alluxio_cluster_dataset_mounted_count_total",
		Help: "Total number of times the dataset was mounted.",
		Type: "Counter",
	},
	"DatasetAliveWorkerTotal": {
		Name: "dataset_alive_worker_total",
		Help: "Total number of alive worker in Dataset.",
		Type: "Gauge",
	},
}

// Define Metrics
var (
	AlluxioClusterAliveWorkerTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: metricDescription["AlluxioClusterAliveWorkerTotal"].Name,
			Help: metricDescription["AlluxioClusterAliveWorkerTotal"].Help,
		},
	)
	AlluxioClusterDatasetMountedCountTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: metricDescription["AlluxioClusterDatasetMountedCountTotal"].Name,
			Help: metricDescription["AlluxioClusterDatasetMountedCountTotal"].Help,
		},
	)
	DatasetAliveWorkerTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: metricDescription["DatasetAliveWorkerTotal"].Name,
			Help: metricDescription["DatasetAliveWorkerTotal"].Help,
		},
	)
)

// RegisterAlluxioControllerMetrics will register Alluxio Controller metrics with the global prometheus registry
func RegisterAlluxioControllerMetrics() {
	metrics.Registry.MustRegister(AlluxioClusterAliveWorkerTotal)
	logger.Infof("Register Metrics: AlluxioClusterAliveWorkerTotal")
	metrics.Registry.MustRegister(AlluxioClusterDatasetMountedCountTotal)
	logger.Infof("Register Metrics: AlluxioClusterDatasetMountedCountTotal")
}

// RegisterDatasetControllerMetrics will register Dataset Controller metrics with the global prometheus registry
func RegisterDatasetControllerMetrics() {
	metrics.Registry.MustRegister(DatasetAliveWorkerTotal)
	logger.Infof("Register Metrics: DatasetAliveWorkerTotal")
}

// ListMetrics will create a slice with the metrics available in metricDescription
func ListMetrics() []MetricDescription {
	v := make([]MetricDescription, 0, len(metricDescription))

	for _, value := range metricDescription {
		v = append(v, value)
	}

	return v
}
