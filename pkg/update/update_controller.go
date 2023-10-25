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

package update

import (
	"context"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	alluxiov1alpha1 "github.com/Alluxio/k8s-operator/api/v1alpha1"
	"github.com/Alluxio/k8s-operator/pkg/logger"
	"github.com/Alluxio/k8s-operator/pkg/utils"
)

type UpdateReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type UpdateReconcilerReqCtx struct {
	*alluxiov1alpha1.AlluxioCluster
	*alluxiov1alpha1.Update
	client.Client
	context.Context
	types.NamespacedName
}

func (r *UpdateReconciler) Reconcile(context context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx := UpdateReconcilerReqCtx{
		Client:         r.Client,
		Context:        context,
		NamespacedName: req.NamespacedName,
	}
	update := &alluxiov1alpha1.Update{}
	if err := r.Get(context, req.NamespacedName, update); err != nil {
		if errors.IsNotFound(err) {
			logger.Infof("update object %v in namespace %v not found. It is being deleted or already deleted.", req.Name, req.Namespace)
		} else {
			logger.Errorf("Failed to get update job %v in namespace %v: %v", req.Name, req.Namespace, err)
			return ctrl.Result{}, err
		}
	}
	ctx.Update = update

	if update.ObjectMeta.UID == "" {
		// TODO: shall we stop the load if still loading?
		return r.deleteJob(ctx)
	}

	dataset := &alluxiov1alpha1.Dataset{}
	datasetNamespacedName := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      update.Spec.Dataset,
	}
	if err := r.Get(context, datasetNamespacedName, dataset); err != nil {
		if errors.IsNotFound(err) {
			logger.Errorf("Dataset %s is not found. Please double check your configuration.", update.Spec.Dataset)
			update.Status.Phase = alluxiov1alpha1.UpdatePhaseFailed
			return r.updateUpdateStatus(ctx)
		} else {
			logger.Errorf("Error getting dataset %s: %v", update.Spec.Dataset, err)
			return ctrl.Result{}, err
		}
	}

	if dataset.Status.Phase != alluxiov1alpha1.DatasetPhaseReady {
		update.Status.Phase = alluxiov1alpha1.UpdatePhaseWaiting
		return r.updateUpdateStatus(ctx)
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

	switch update.Status.Phase {
	case alluxiov1alpha1.UpdatePhaseNone, alluxiov1alpha1.UpdatePhaseWaiting:
		return r.createUpdateJob(ctx)
	case alluxiov1alpha1.UpdatePhaseUpdating:
		return r.waitUpdateJobFinish(ctx)
	default:
		return ctrl.Result{}, nil
	}
}

func (r *UpdateReconciler) waitUpdateJobFinish(ctx UpdateReconcilerReqCtx) (ctrl.Result, error) {
	updateJob, err := r.getUpdateJob(ctx)
	if err != nil {
		return ctrl.Result{}, err
	}
	if updateJob.Status.Succeeded == 1 {
		ctx.Update.Status.Phase = alluxiov1alpha1.UpdatePhaseUpdated
		if _, err := r.updateUpdateStatus(ctx); err != nil {
			logger.Errorf("Data is updated but failed to update status. %v", err)
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{}, nil
	} else if updateJob.Status.Failed == 1 {
		ctx.Update.Status.Phase = alluxiov1alpha1.UpdatePhaseFailed
		if _, err := r.updateUpdateStatus(ctx); err != nil {
			logger.Errorf("Failed to update status. %v", err)
			return ctrl.Result{Requeue: true}, err
		}
		logger.Errorf("update data job failed. Please check the log of the pod for errors.")
		return ctrl.Result{}, nil
	} else {
		return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
	}
}

func (r *UpdateReconciler) getUpdateJob(ctx UpdateReconcilerReqCtx) (*batchv1.Job, error) {
	updateJob := &batchv1.Job{}
	updateJobNamespacedName := types.NamespacedName{
		Name:      utils.GetUpdateJobName(ctx.Name),
		Namespace: ctx.Namespace,
	}
	if err := r.Get(ctx.Context, updateJobNamespacedName, updateJob); err != nil {
		logger.Errorf("Error getting update job %s: %v", ctx.NamespacedName.String(), err)
		return nil, err
	}
	return updateJob, nil
}

func (r *UpdateReconciler) updateUpdateStatus(ctx UpdateReconcilerReqCtx) (ctrl.Result, error) {
	if err := r.Client.Status().Update(ctx.Context, ctx.Update); err != nil {
		logger.Errorf("Failed updating update job status: %v", err)
		return ctrl.Result{}, err
	}
	return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UpdateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&alluxiov1alpha1.Update{}).
		Complete(r)
}
