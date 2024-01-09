package api_server

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"time"
)

var (
	defaultRetryPeriod = 2 * time.Second
)

type Options struct {
	Scheme         *runtime.Scheme
	Namespace      string
	Port           int
	AllowedDomains []string
}

type Manager interface {
	Start() error
}

type manager struct {
	config         *rest.Config
	client         client.Client
	server         *apiServer
	started        bool
	errSignal      *errSignaler
	port           int
	allowedDomains []string
}

func NewManager(config *rest.Config, options Options) (Manager, error) {
	mapper, err := apiutil.NewDynamicRESTMapper(config)
	if err != nil {
		return nil, err
	}

	c, err := client.New(config, client.Options{Scheme: options.Scheme, Mapper: mapper})
	if err != nil {
		return nil, err
	}

	return &manager{
		config:         config,
		client:         c,
		port:           options.Port,
		allowedDomains: options.AllowedDomains,
	}, nil
}

func (m *manager) Start() error {
	// initialize this here so that we reset the signal channel state on every start
	m.errSignal = &errSignaler{errSignal: make(chan struct{})}

	srv, err := newApiServer(m.port, m.allowedDomains, m.client)
	if err != nil {
		return err
	}

	go func() {
		if err := srv.Start(); err != nil {
			m.errSignal.SignalError(err)
		}
	}()
	select {
	case <-m.errSignal.GotError():
		// Error starting the cache
		return m.errSignal.Error()
	}
}
