package main

import (
	alluxiov1alpha1 "github.com/Alluxio/k8s-operator/api/v1alpha1"
	apiserver "github.com/Alluxio/k8s-operator/cmd/api_server/api_server"
	"github.com/Alluxio/k8s-operator/pkg/logger"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	kubelog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	mainLog = kubelog.Log.WithName("alluxio-operator-api-server").WithName("main")
	scheme  = runtime.NewScheme()
)

func NewAPIServerManagerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api-server-manager start",
		Short: "Start the manager of api-server",
		Run: func(cmd *cobra.Command, args []string) {
			startAPIServerManager()
		},
	}
	return cmd
}

func startAPIServerManager() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(alluxiov1alpha1.AddToScheme(scheme))

	kubelog.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := apiserver.NewManager(ctrl.GetConfigOrDie(), apiserver.Options{
		Scheme:         scheme,
		Port:           5220,
		AllowedDomains: []string{},
	})
	if err != nil {
		mainLog.Error(err, "Unable to create api-server manager")
		os.Exit(1)
	}

	mainLog.Info("Starting api-server manager.")
	if err := mgr.Start(); err != nil {
		mainLog.Error(err, "Problem running api-server manager")
		os.Exit(1)
	}
}

func main() {
	command := NewAPIServerManagerCommand()
	if err := command.Execute(); err != nil {
		logger.Fatalf("Failed to launch API Server: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
