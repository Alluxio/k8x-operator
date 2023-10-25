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
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	alluxiov1alpha1 "github.com/Alluxio/k8s-operator/api/v1alpha1"
	"github.com/Alluxio/k8s-operator/pkg/logger"
	"github.com/Alluxio/k8s-operator/pkg/utils"
)

func (r *UnloadReconciler) unload(ctx UnloadReconcilerReqCtx) (ctrl.Result, error) {
	// Update the status before starting the command instead of after, because otherwise if the status update fails,
	// the reconciler will loop again and redo the same thing, leading to a confusing state.
	ctx.Unload.Status.Phase = alluxiov1alpha1.UnloadPhaseUnLoaded
	_, err := r.updateUnloadStatus(ctx)
	if err != nil {
		logger.Infof("Unloading is pending because status was not updated successfully")
		return ctrl.Result{RequeueAfter: 15 * time.Second}, err
	}
	r.unloadInternal(ctx)
	return ctrl.Result{}, nil
}

func (r *UnloadReconciler) unloadInternal(ctx UnloadReconcilerReqCtx) {
	var workerPodsList corev1.PodList
	workerListOpts := client.MatchingLabels{
		"name": fmt.Sprintf("%s-worker", ctx.AlluxioCluster.ObjectMeta.Name),
	}
	if err := r.List(ctx.Context, &workerPodsList, client.InNamespace(ctx.Namespace), workerListOpts); err != nil {
		logger.Warnf("Failed to get the list of worker pods while trying to delete pages: %v", err)
	}
	clientSet, err := utils.GetK8sClient()
	if err != nil {
		logger.Warnf("Failed to get clientSet while trying to delete pages: %v", err)
	}
	var wg sync.WaitGroup
	for _, workerPod := range workerPodsList.Items {
		wg.Add(1)
		go func() {
			defer wg.Done()
			removeAllPagesFromOneWorkerPod(workerPod.Name, ctx, clientSet)
		}()
	}
	wg.Wait()
}

func removeAllPagesFromOneWorkerPod(podName string, ctx UnloadReconcilerReqCtx, clientSet *kubernetes.Clientset) {
	// Best effort to remove cached pages
	paths := strings.Split(ctx.AlluxioCluster.Spec.Pagestore.HostPath, ",")
	for i, path := range paths {
		paths[i] = fmt.Sprintf("%s/*", path)
	}
	cmd := []string{"sh", "-c", fmt.Sprintf("rm -r %s", strings.Join(paths, " "))}
	req := clientSet.CoreV1().RESTClient().Post().Resource("pods").Name(podName).Namespace(ctx.Namespace).SubResource("exec")
	execOpts := &corev1.PodExecOptions{
		Command: cmd,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	req.VersionedParams(execOpts, scheme.ParameterCodec)
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Warnf("Failed to get config for k8s clientSet. Remove pages is aborted. %v", err)
	}
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		logger.Warnf("Failed to create executor. Remove pages is aborted. %v")
	}
	if err := exec.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}); err != nil {
		logger.Warnf("Failed to execute removing pages command. The command is aborted. %v", err)
	}
}
