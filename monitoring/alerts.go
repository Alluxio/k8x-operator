package monitoring

import (
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	ruleName            = "alluxio-operator-rules"
	alertRuleGroup      = "alluxio.rules"
	datasetMountedAlert = "AlluxioClusterDatasetMounted"
	// TODO: Define more rules here
)

// NewPrometheusRule creates new PrometheusRule(CR) for the operator to have alerts and recording rules
func NewPrometheusRule(namespace string) *monitoringv1.PrometheusRule {
	return &monitoringv1.PrometheusRule{
		TypeMeta: metav1.TypeMeta{
			APIVersion: monitoringv1.SchemeGroupVersion.String(),
			Kind:       "PrometheusRule",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ruleName,
			Namespace: namespace,
		},
		Spec: *NewPrometheusRuleSpec(),
	}
}

// NewPrometheusRuleSpec creates PrometheusRuleSpec for alerts and recording rules
func NewPrometheusRuleSpec() *monitoringv1.PrometheusRuleSpec {
	return &monitoringv1.PrometheusRuleSpec{
		Groups: []monitoringv1.RuleGroup{{
			Name: alertRuleGroup,
			Rules: []monitoringv1.Rule{
				createDatasetMountedAlertRule(),
				// TODO: Add more rules here
			},
		}},
	}
}

func createDatasetMountedAlertRule() monitoringv1.Rule {
	return monitoringv1.Rule{
		Alert: datasetMountedAlert,
		Expr:  intstr.FromString("increase(alluxio_cluster_dataset_mounted_count_total[1m]) >= 1"),
		Annotations: map[string]string{
			"description": "Total number of times the dataset was mounted. more than 1 times in the last 1 minutes.",
		},
		Labels: map[string]string{
			"severity": "warning",
		},
	}
}
