package endpoints

import (
	"fmt"
	"github.com/Alluxio/k8s-operator/api/v1alpha1"
	"github.com/emicklei/go-restful"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AlluxioClusterEndpoint struct {
	client client.Client
}

func NewAlluxioClusterEndpoint(client client.Client) *AlluxioClusterEndpoint {
	return &AlluxioClusterEndpoint{client: client}
}

func (alluxioClusterEndpoint *AlluxioClusterEndpoint) SetupWithWS(ws *restful.WebService) {
	ws.Route(ws.GET("alluxio_cluster").To(alluxioClusterEndpoint.show).
		Doc("List of AlluxioClusters").
		Returns(200, "OK", &v1alpha1.AlluxioClusterList{}))

	ws.Route(ws.POST("alluxio_cluster").To(alluxioClusterEndpoint.create).
		Doc("Create AlluxioClusters").
		Returns(200, "OK", nil).
		Returns(400, "Bad Request", nil))

	ws.Route(ws.DELETE("alluxio_cluster").To(alluxioClusterEndpoint.delete).
		Doc("Delete Current AlluxioClusters").
		Returns(200, "OK", nil).
		Returns(400, "Bad Request", nil))

	ws.Route(ws.PATCH("alluxio_cluster").To(alluxioClusterEndpoint.update).
		Doc("Update Current AlluxioClusters").
		Returns(200, "OK", nil).
		Returns(400, "Bad Request", nil))

}

func (alluxioClusterEndpoint *AlluxioClusterEndpoint) show(request *restful.Request, response *restful.Response) {
	// Get AlluxioClusterList
	AlluxioClusterList := new(v1alpha1.AlluxioClusterList)
	err := alluxioClusterEndpoint.client.List(request.Request.Context(), AlluxioClusterList, &client.ListOptions{})
	// Unable to get:
	if err != nil {
		writeError(response, 404, Error{
			Title:   "Error",
			Details: fmt.Sprintf("Could not retrieve list: %s", err),
		})
		return
	}
	// Write to response:
	newAL := AlluxioClusterConverter.AlluxioClusterList(AlluxioClusterList)
	if err := response.WriteAsJson(newAL); err != nil {
		writeError(response, 404, Error{
			Title:   "Error",
			Details: "Could not list resources",
		})
	}
}

type Test struct {
	Title string `json:"title"`
}

func (alluxioClusterEndpoint *AlluxioClusterEndpoint) create(request *restful.Request, response *restful.Response) {
	test := new(Test)

	err := request.ReadEntity(test)

	if err != nil {
		writeError(response, 400, Error{
			Title:   "Bad Request",
			Details: "Could not read entity",
		})
		return
	}

	writeError(response, 400, Error{
		Title:   "Error",
		Details: "To Do, can read entity",
	})
}

func (alluxioClusterEndpoint *AlluxioClusterEndpoint) update(request *restful.Request, response *restful.Response) {
	writeError(response, 400, Error{
		Title:   "Error",
		Details: "To Do",
	})
}

func (alluxioClusterEndpoint *AlluxioClusterEndpoint) delete(request *restful.Request, response *restful.Response) {
	writeError(response, 400, Error{
		Title:   "Error",
		Details: "To Do",
	})
}
