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

package unload

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"

	alluxiov1alpha1 "github.com/Alluxio/k8s-operator/api/v1alpha1"
	"github.com/Alluxio/k8s-operator/pkg/logger"
)

type UnloadReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type UnloadReconcilerReqCtx struct {
	*alluxiov1alpha1.AlluxioCluster
	*alluxiov1alpha1.Unload
	client.Client
	context.Context
	types.NamespacedName
}

func (r *UnloadReconciler) Reconcile(context context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx := UnloadReconcilerReqCtx{
		Client:         r.Client,
		Context:        context,
		NamespacedName: req.NamespacedName,
	}
	unload := &alluxiov1alpha1.Unload{}
	if err := r.Get(context, req.NamespacedName, unload); err != nil {
		if errors.IsNotFound(err) {
			logger.Infof("Unload object %v in namespace %v not found. It is being deleted or already deleted.", req.Name, req.Namespace)
		} else {
			logger.Errorf("Failed to get unload job %v in namespace %v: %v", req.Name, req.Namespace, err)
			return ctrl.Result{}, err
		}
	}
	ctx.Unload = unload

	if unload.ObjectMeta.UID == "" {
		return ctrl.Result{}, nil
	}

	dataset := &alluxiov1alpha1.Dataset{}
	datasetNamespacedName := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      unload.Spec.Dataset,
	}
	if err := r.Get(context, datasetNamespacedName, dataset); err != nil {
		if errors.IsNotFound(err) {
			logger.Errorf("Dataset %s is not found. Please double check your configuration.", unload.Spec.Dataset)
			unload.Status.Phase = alluxiov1alpha1.UnloadPhaseNotExist
			return r.updateUnloadStatus(ctx)
		} else {
			logger.Errorf("Error getting dataset %s: %v", unload.Spec.Dataset, err)
			return ctrl.Result{}, err
		}
	}

	if dataset.Status.Phase != alluxiov1alpha1.DatasetPhaseReady {
		unload.Status.Phase = alluxiov1alpha1.UnloadPhaseNotExist
		return r.updateUnloadStatus(ctx)
	}

	alluxioCluster := &alluxiov1alpha1.AlluxioCluster{}
	alluxioNamespacedName := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      dataset.Status.BoundedAlluxioCluster,
	}
	if err := r.Get(context, alluxioNamespacedName, alluxioCluster); err != nil {
		logger.Errorf("Error getting alluxio cluster %s: %v", alluxioNamespacedName.Name, err)
		return ctrl.Result{}, err
	}
	ctx.AlluxioCluster = alluxioCluster

	switch unload.Status.Phase {
	case alluxiov1alpha1.UnloadPhaseNone, alluxiov1alpha1.UnloadPhaseNotExist:
		return r.unload(ctx)
	default:
		return ctrl.Result{}, nil
	}
}

func (r *UnloadReconciler) updateUnloadStatus(ctx UnloadReconcilerReqCtx) (ctrl.Result, error) {
	if err := r.Client.Status().Update(ctx.Context, ctx.Unload); err != nil {
		logger.Errorf("Failed updating unload job status: %v", err)
		return ctrl.Result{}, err
	}
	return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UnloadReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&alluxiov1alpha1.Unload{}).
		Complete(r)
}
