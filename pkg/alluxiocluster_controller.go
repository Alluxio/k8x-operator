/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pkg

import (
	"context"
	alluxiocomv1alpha1 "github.com/Alluxio/k8s-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

// AlluxioClusterReconciler reconciles a AlluxioCluster object
type AlluxioClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type ReconcileRequestContext struct {
	AlluxioCluster client.Object
	client.Client
	context.Context
	Logger logr.Logger
	types.NamespacedName
}

type ReconcileResponse struct {
	Err        error
	NeedReturn bool
	Result     ctrl.Result
}

func (r *AlluxioClusterReconciler) Reconcile(context context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(context)
	logger.Info("Reconciling.", "namespace", req.Namespace, "name", req.Name)
	ctx := ReconcileRequestContext{
		Client:         r.Client,
		Context:        context,
		Logger:         logger,
		NamespacedName: req.NamespacedName,
	}

	alluxioCluster := &alluxiocomv1alpha1.AlluxioCluster{}
	ctx.AlluxioCluster = alluxioCluster

	err := r.Get(ctx, req.NamespacedName, alluxioCluster)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Alluxio cluster's resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Alluxio Cluster")
		return ctrl.Result{}, err
	}

	if alluxioCluster.Status.Phase == alluxiocomv1alpha1.ClusterPhaseNone {
		alluxioCluster.Status.Phase = alluxiocomv1alpha1.ClusterPhaseCreating
		if res := CreateAlluxioCluster(ctx); res.NeedReturn {
			return res.Result, res.Err
		}
	}

	return ctrl.Result{RequeueAfter: 2 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AlluxioClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&alluxiocomv1alpha1.AlluxioCluster{}).
		Complete(r)
}