package endpoints

import (
	"github.com/emicklei/go-restful/v3"
	"sigs.k8s.io/controller-runtime/pkg/client"
	kubelog "sigs.k8s.io/controller-runtime/pkg/log"
)

type AlluxioOperatorEndpoint struct {
	client client.Client
}

func NewDarkroomEndpoint(client client.Client) *AlluxioOperatorEndpoint {
	return &AlluxioOperatorEndpoint{client: client}
}

func (ep *AlluxioOperatorEndpoint) SetupWithWS(ws *restful.WebService) {
	ws.Route(ws.GET("darkrooms").To(ep.client).
		Doc("List of Darkrooms").
		Returns(200, "OK", &List{}))

	ws.Route(ws.POST("darkrooms").To(de.create).
		Doc("Create a new darkroom").
		Reads(&Darkroom{}).
		Returns(200, "OK", &AlluxioOperator{}).
		Returns(400, "Bad Request", nil))
}

func writeError(response *restful.Response, httpStatus int, err Error) {
	if err := response.WriteHeaderAndJson(httpStatus, err, "application/json"); err != nil {
		kubelog.Log.Error(err, "Could not write the error response")
	}
}
