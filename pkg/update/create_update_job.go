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
	"fmt"
	"os"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	alluxiov1alpha1 "github.com/Alluxio/k8s-operator/api/v1alpha1"
	"github.com/Alluxio/k8s-operator/pkg/logger"
	"github.com/Alluxio/k8s-operator/pkg/utils"
)

func (r *UpdateReconciler) createUpdateJob(ctx UpdateReconcilerReqCtx) (ctrl.Result, error) {
	// Update the status before job creation instead of after, because otherwise if the status update fails,
	// the reconciler will loop again and create another same job, leading to failure to create duplicated job which is confusing.
	ctx.Update.Status.Phase = alluxiov1alpha1.UpdatePhaseUpdating
	_, err := r.updateUpdateStatus(ctx)
	if err != nil {
		logger.Infof("Job is pending because status was not updated successfully")
		return ctrl.Result{}, err
	}
	updateJob, err := getUpdateJobFromYaml()
	if err != nil {
		return ctrl.Result{}, err
	}
	constructUpdateJob(ctx.AlluxioCluster, ctx.Update, updateJob)
	if err := r.Create(ctx.Context, updateJob); err != nil {
		logger.Errorf("Failed to update data of dataset %s: %v", ctx.NamespacedName.String(), err)
		return ctrl.Result{}, err
	}
	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func getUpdateJobFromYaml() (*batchv1.Job, error) {
	updateJobYaml, err := os.ReadFile("/opt/alluxio-jobs/update.yaml")
	if err != nil {
		logger.Errorf("Failed to read update job yaml file at /opt/alluxio-jobs/update.yaml: %v", err)
		return nil, err
	}
	udpateJob, _, err := scheme.Codecs.UniversalDeserializer().Decode(updateJobYaml, nil, nil)
	if err != nil {
		logger.Errorf("Failed to parse update job yaml file: %v", err)
	}
	return udpateJob.(*batchv1.Job), nil
}

func constructUpdateJob(alluxio *alluxiov1alpha1.AlluxioCluster, update *alluxiov1alpha1.Update, updateJob *batchv1.Job) {
	updateJob.Name = utils.GetUpdateJobName(update.Name)
	updateJob.Namespace = alluxio.Namespace
	var imagePullSecrets []corev1.LocalObjectReference
	for _, secret := range alluxio.Spec.ImagePullSecrets {
		imagePullSecrets = append(imagePullSecrets, corev1.LocalObjectReference{Name: secret})

	}
	updateJob.Spec.Template.Spec.ImagePullSecrets = imagePullSecrets
	updateJob.Spec.Template.Spec.ServiceAccountName = alluxio.Spec.ServiceAccountName
	updateJob.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("%s:%s", alluxio.Spec.Image, alluxio.Spec.ImageTag)
	updateJob.Spec.Template.Spec.Containers[0].Command = []string{"go", "run", "/update.go", update.Spec.Path}
	alluxioConfigMapName := utils.GetAlluxioConfigMapName(alluxio.Spec.NameOverride, alluxio.Name)
	updateConfigMapName := utils.GetUpdateConfigmapName(alluxio.Spec.NameOverride, alluxio.Name)
	updateJob.Spec.Template.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{
		{
			Name:      alluxioConfigMapName,
			MountPath: "/opt/alluxio/conf",
		},
		{
			Name:      updateConfigMapName,
			MountPath: "/update.go",
			SubPath:   "update.go",
		},
	}
	updateJob.Spec.Template.Spec.Volumes = []corev1.Volume{
		{
			Name: alluxioConfigMapName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: alluxioConfigMapName,
					},
				},
			},
		},
		{
			Name: updateConfigMapName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: updateConfigMapName,
					},
				},
			},
		},
	}
}
