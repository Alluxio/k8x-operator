package api_server

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/Alluxio/k8s-operator/cmd/api_server/api_server/endpoints"
	"github.com/emicklei/go-restful"
	"io/fs"
	"log"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	kubelog "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	serverLogger = kubelog.Log.WithName("Alluxio Operator api-server")
)

type apiServer struct {
	server *http.Server
}

func (as *apiServer) Address() string {
	return as.server.Addr
}

func init() {
	restful.MarshalIndent = func(v interface{}, prefix, indent string) ([]byte, error) {
		var buf bytes.Buffer
		encoder := restful.NewEncoder(&buf)
		encoder.SetIndent(prefix, indent)
		if err := encoder.Encode(v); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}

//go:embed gui
var embeddedFiles embed.FS

func newApiServer(port int, allowedDomains []string, client client.Client) (*apiServer, error) {
	container := restful.NewContainer()
	// Add Web UI
	guiFS, err := fs.Sub(embeddedFiles, "gui")
	if err != nil {
		log.Fatal(err)
	}
	// Declare Main Page for webui.
	container.ServeMux.Handle("/", http.FileServer(http.FS(guiFS)))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: container.ServeMux,
	}

	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{restful.HEADER_AccessControlAllowOrigin},
		AllowedDomains: allowedDomains,
		Container:      container,
	}

	// RESTFUL API start at `/api`
	ws := new(restful.WebService)
	ws.
		Path("/api").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	addEndpoints(ws, client)
	container.Add(ws)
	container.Filter(cors.Filter)
	return &apiServer{
		server: srv,
	}, nil
}

func addEndpoints(ws *restful.WebService, client client.Client) {
	resources := []endpoints.Endpoint{
		endpoints.NewAlluxioClusterEndpoint(client),
		endpoints.NewDataSetEndpoint(client),
	}
	for _, ep := range resources {
		ep.SetupWithWS(ws)
	}
}

func (as *apiServer) Start() error {
	errChan := make(chan error)
	go func() {
		err := as.server.ListenAndServe()
		if err != nil {
			switch err {
			case http.ErrServerClosed:
				serverLogger.Info("Shutting down api-server")
			default:
				serverLogger.Error(err, "Could not start an HTTP Server")
				errChan <- err
			}
		}
	}()
	serverLogger.Info("Starting api-server", "interface", "0.0.0.0", "port", as.Address())
	select {
	case err := <-errChan:
		return err
	}
}
