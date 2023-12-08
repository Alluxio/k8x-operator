/*
 * The Alluxio Open Foundation licenses this work under the Apache License, version 2.0
 * (the "License"). You may not use this work except in compliance with the License, which is
 * available at www.apache.org/licenses/LICENSE-2.0
 *
 * This software is distributed on an "AS IS" basis, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied, as more fully set forth in the License.
 *
 * See the NOTICE file distributed with this work for information regarding copyright ownership.
 */

package alluxiocluster

import (
	"context"
	"github.com/Alluxio/k8s-operator/monitoring"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	alluxiov1alpha1 "github.com/Alluxio/k8s-operator/api/v1alpha1"
	"github.com/Alluxio/k8s-operator/pkg/finalizer"
	"github.com/Alluxio/k8s-operator/pkg/logger"
)

// AlluxioClusterReconciler reconciles a AlluxioCluster object
type AlluxioClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type AlluxioClusterReconcileReqCtx struct {
	*alluxiov1alpha1.AlluxioCluster
	client.Client
	context.Context
	*alluxiov1alpha1.Dataset
	types.NamespacedName
}

const ruleName = "alluxio-operator-rules"
const namespace = "alluxio-operator-system"

func (r *AlluxioClusterReconciler) Reconcile(context context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger.Infof("Reconciling AlluxioCluster %s", req.NamespacedName.String())
	ctx := AlluxioClusterReconcileReqCtx{
		Client:         r.Client,
		Context:        context,
		NamespacedName: req.NamespacedName,
	}

	// Check if prometheus rule already exists, if not create a new one
	foundRule := &monitoringv1.PrometheusRule{}
	if err := r.Get(ctx, types.NamespacedName{Name: ruleName, Namespace: namespace}, foundRule); err != nil {
		if apierrors.IsNotFound(err) {
			// Define a new prometheus rule
			prometheusRule := monitoring.NewPrometheusRule(namespace)
			if err := r.Create(ctx, prometheusRule); err != nil {
				logger.Errorf("Failed to create prometheus rule:  %v", err)
				return ctrl.Result{}, nil
			}
		}
	} else {
		// Check if prometheus rule spec was changed, if so set as desired
		desiredRuleSpec := monitoring.NewPrometheusRuleSpec()
		if !reflect.DeepEqual(foundRule.Spec.DeepCopy(), desiredRuleSpec) {
			desiredRuleSpec.DeepCopyInto(&foundRule.Spec)
			if r.Update(ctx, foundRule); err != nil {
				logger.Errorf("Failed to update prometheus rule:  %v", err)
				return ctrl.Result{}, nil
			}
		}
	}
	// END OF prometheus rule

	alluxioCluster := &alluxiov1alpha1.AlluxioCluster{}
	ctx.AlluxioCluster = alluxioCluster
	if err := r.Get(context, req.NamespacedName, alluxioCluster); err != nil {
		if errors.IsNotFound(err) {
			logger.Infof("Alluxio cluster %s not found. It is being deleted or already deleted.", req.NamespacedName.String())
		} else {
			logger.Errorf("Failed to get Alluxio cluster %s: %v", req.NamespacedName.String(), err)
			return ctrl.Result{}, err
		}
	}
	dataset := &alluxiov1alpha1.Dataset{}
	datasetNamespacedName := types.NamespacedName{
		Name:      alluxioCluster.Spec.Dataset,
		Namespace: req.Namespace,
	}
	ctx.Dataset = dataset
	if err := r.Get(context, datasetNamespacedName, dataset); err != nil {
		if errors.IsNotFound(err) {
			logger.Infof("Dataset %s not found. It is deleted or hasn't been created yet.", alluxioCluster.Spec.Dataset)
			dataset.Status.Phase = alluxiov1alpha1.DatasetPhaseNotExist
		} else {
			logger.Errorf("Failed to get Dataset %s: %v", req.NamespacedName.String(), err)
			return ctrl.Result{}, err
		}
	}
	//} else {
	//	monitoring.AlluxioClusterDatasetMountedCountTotal.Inc()
	//}

	if alluxioCluster.DeletionTimestamp != nil {
		if err := deleteConfYamlFileIfExist(ctx.NamespacedName); err != nil {
			return ctrl.Result{}, err
		}
		if err := DeleteAlluxioClusterIfExist(ctx.NamespacedName); err != nil {
			return ctrl.Result{}, err
		}
		if dataset.Status.Phase != alluxiov1alpha1.DatasetPhaseNotExist {
			ctx.Dataset.Status.Phase = alluxiov1alpha1.DatasetPhasePending
			ctx.Dataset.Status.BoundedAlluxioCluster = ""
			if err := updateDatasetStatus(ctx); err != nil {
				return ctrl.Result{}, err
			}
		}
		if err := finalizer.RemoveDummyFinalizerIfExist(r.Client, alluxioCluster, context); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if dataset.Status.Phase == alluxiov1alpha1.DatasetPhaseNotExist {
		if err := deleteConfYamlFileIfExist(ctx.NamespacedName); err != nil {
			return ctrl.Result{}, err
		}
		if err := DeleteAlluxioClusterIfExist(ctx.NamespacedName); err != nil {
			return ctrl.Result{}, err
		}
		return UpdateStatus(ctx)
	}

	if alluxioCluster.Status.Phase == alluxiov1alpha1.ClusterPhaseNone || alluxioCluster.Status.Phase == alluxiov1alpha1.ClusterPhasePending {
		if err := finalizer.AddDummyFinalizerIfNotExist(r.Client, alluxioCluster, context); err != nil {
			return ctrl.Result{}, err
		}
		if err := CreateAlluxioClusterIfNotExist(ctx); err != nil {
			return ctrl.Result{}, err
		}
		return UpdateStatus(ctx)
	}

	return UpdateStatus(ctx)
}

// SetupWithManager sets up the controller with the Manager.
func (r *AlluxioClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&alluxiov1alpha1.AlluxioCluster{}).
		Complete(r)
}
