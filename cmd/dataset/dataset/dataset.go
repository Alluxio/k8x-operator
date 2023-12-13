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

package dataset

import (
	"flag"
	"github.com/Alluxio/k8s-operator/monitoring"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	alluxiov1alpha1 "github.com/Alluxio/k8s-operator/api/v1alpha1"
	"github.com/Alluxio/k8s-operator/pkg/dataset"
	"github.com/Alluxio/k8s-operator/pkg/load"
	"github.com/Alluxio/k8s-operator/pkg/logger"
	"github.com/Alluxio/k8s-operator/pkg/unload"
	"github.com/Alluxio/k8s-operator/pkg/update"
)

var (
	scheme = runtime.NewScheme()
)

func NewDatasetManagerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dataset-manager start",
		Short: "Start the manager of dataset",
		Run: func(cmd *cobra.Command, args []string) {
			startDatasetManager()
		},
	}
	return cmd
}

func startDatasetManager() {
	monitoring.RegisterDatasetControllerMetrics()

	var metricsAddr string
	var probeAddr string

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(alluxiov1alpha1.AddToScheme(scheme))

	manager, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
	})
	if err != nil {
		logger.Fatalf("Unable to create Dataset manager: %v", err)
		os.Exit(1)
	}

	if err = (&dataset.DatasetReconciler{
		Client: manager.GetClient(),
		Scheme: manager.GetScheme(),
	}).SetupWithManager(manager); err != nil {
		logger.Fatalf("unable to create Dataset controller: %v", err)
		os.Exit(1)
	}

	if err = (&load.LoadReconciler{
		Client: manager.GetClient(),
		Scheme: manager.GetScheme(),
	}).SetupWithManager(manager); err != nil {
		logger.Fatalf("unable to create Load controller: %v", err)
		os.Exit(1)
	}

	if err = (&update.UpdateReconciler{
		Client: manager.GetClient(),
		Scheme: manager.GetScheme(),
	}).SetupWithManager(manager); err != nil {
		logger.Fatalf("unable to create Update controller: %v", err)
		os.Exit(1)
	}

	if err = (&unload.UnloadReconciler{
		Client: manager.GetClient(),
		Scheme: manager.GetScheme(),
	}).SetupWithManager(manager); err != nil {
		logger.Fatalf("unable to create Unload controller: %v", err)
		os.Exit(1)
	}

	logger.Infof("starting manager")
	if err = manager.Start(ctrl.SetupSignalHandler()); err != nil {
		logger.Fatalf("Error starting Dataset manager: %v", err)
		os.Exit(1)
	}
}
