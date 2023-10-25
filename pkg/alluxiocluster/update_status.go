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
	"reflect"
	"time"

	"github.com/jinzhu/copier"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	alluxiov1alpha1 "github.com/Alluxio/k8s-operator/api/v1alpha1"
	"github.com/Alluxio/k8s-operator/pkg/dataset"
	"github.com/Alluxio/k8s-operator/pkg/logger"
	"github.com/Alluxio/k8s-operator/pkg/utils"
)

func UpdateStatus(alluxioClusterCtx AlluxioClusterReconcileReqCtx) (ctrl.Result, error) {
	alluxioOriginalStatusCopy := alluxioClusterCtx.AlluxioCluster.Status.DeepCopy()
	datasetOriginalStatusCopy := alluxioClusterCtx.Dataset.Status.DeepCopy()

	alluxioClusterNewPhase := alluxioOriginalStatusCopy.Phase
	datasetNewPhase := alluxioClusterCtx.Dataset.Status.Phase

	if datasetOriginalStatusCopy.Phase == alluxiov1alpha1.DatasetPhaseNotExist {
		alluxioClusterNewPhase = alluxiov1alpha1.ClusterPhasePending
	} else if alluxioOriginalStatusCopy.Phase == alluxiov1alpha1.ClusterPhaseNone {
		alluxioClusterNewPhase = alluxiov1alpha1.ClusterPhaseCreatingOrUpdating
		datasetNewPhase = alluxiov1alpha1.DatasetPhaseBounding
	} else {
		if ClusterReady(alluxioClusterCtx) {
			alluxioClusterNewPhase = alluxiov1alpha1.ClusterPhaseReady
			datasetNewPhase = alluxiov1alpha1.DatasetPhaseReady
			alluxioClusterCtx.Dataset.Status.BoundedAlluxioCluster = alluxioClusterCtx.AlluxioCluster.Name
		} else {
			alluxioClusterNewPhase = alluxiov1alpha1.ClusterPhaseCreatingOrUpdating
			datasetNewPhase = alluxiov1alpha1.DatasetPhaseBounding
		}
	}
	alluxioClusterCtx.AlluxioCluster.Status.Phase = alluxioClusterNewPhase
	alluxioClusterCtx.Dataset.Status.Phase = datasetNewPhase

	if !reflect.DeepEqual(alluxioOriginalStatusCopy, alluxioClusterCtx.AlluxioCluster.Status) {
		if err := alluxioClusterCtx.Client.Status().Update(alluxioClusterCtx.Context, alluxioClusterCtx.AlluxioCluster); err != nil {
			logger.Errorf("Error updating cluster status: %v", err)
			return ctrl.Result{}, err
		}
	}
	if !reflect.DeepEqual(*datasetOriginalStatusCopy, alluxioClusterCtx.Dataset.Status) {
		if err := updateDatasetStatus(alluxioClusterCtx); err != nil {
			return ctrl.Result{}, err
		}
	}

	if alluxioClusterNewPhase != alluxiov1alpha1.ClusterPhaseReady {
		return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
	}
	return ctrl.Result{RequeueAfter: 2 * time.Minute}, nil
}

func ClusterReady(ctx AlluxioClusterReconcileReqCtx) bool {
	componentStatusReqCtx := utils.ComponentStatusReqCtx{}
	copier.Copy(&componentStatusReqCtx, &ctx)
	master, err := utils.GetMasterStatus(componentStatusReqCtx)
	if err != nil || master.Status.AvailableReplicas != master.Status.Replicas {
		return false
	}
	worker, err := utils.GetWorkerStatus(componentStatusReqCtx)
	if err != nil || worker.Status.AvailableReplicas != worker.Status.Replicas {
		return false
	}
	if ctx.AlluxioCluster.Spec.Fuse.Enabled != nil && *ctx.AlluxioCluster.Spec.Fuse.Enabled {
		fuse, err := utils.GetFuseStatus(componentStatusReqCtx)
		if err != nil || fuse.Status.NumberAvailable != fuse.Status.DesiredNumberScheduled {
			return false
		}
	}
	if ctx.AlluxioCluster.Spec.Proxy.Enabled != nil && *ctx.AlluxioCluster.Spec.Proxy.Enabled {
		proxy, err := utils.GetProxyStatus(componentStatusReqCtx)
		if err != nil || proxy.Status.NumberAvailable != proxy.Status.DesiredNumberScheduled {
			return false
		}
	}
	return true
}

func updateDatasetStatus(alluxioClusterCtx AlluxioClusterReconcileReqCtx) error {
	datasetCtx := dataset.DatasetReconcilerReqCtx{
		Dataset: alluxioClusterCtx.Dataset,
		Client:  alluxioClusterCtx.Client,
		Context: alluxioClusterCtx.Context,
		NamespacedName: types.NamespacedName{
			Name:      alluxioClusterCtx.Dataset.Name,
			Namespace: alluxioClusterCtx.Namespace,
		},
	}
	if _, err := dataset.UpdateDatasetStatus(datasetCtx); err != nil {
		return err
	}
	return nil
}
