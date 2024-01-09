package endpoints

import (
	"fmt"
	"github.com/Alluxio/k8s-operator/api/v1alpha1"
	"github.com/Alluxio/k8s-operator/pkg/logger"
	"github.com/emicklei/go-restful"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		Doc("List of Alluxio Clusters").
		Returns(200, "OK", AlluxioClusterList{}))

	ws.Route(ws.POST("alluxio_cluster").To(alluxioClusterEndpoint.create).
		Doc("Create a new Alluxio Cluster").
		Returns(200, "OK", AlluxioCluster{}).
		Returns(400, "Bad Request", nil))

	ws.Route(ws.DELETE("alluxio_cluster").To(alluxioClusterEndpoint.delete).
		Doc("Delete a Current Alluxio Cluster").
		Returns(200, "OK", nil).
		Returns(400, "Bad Request", nil))

	//ws.Route(ws.PATCH("alluxio_cluster").To(alluxioClusterEndpoint.update).
	//	Doc("Update Current AlluxioClusters").
	//	Returns(200, "OK", nil).
	//	Returns(400, "Bad Request", nil))
}

func (alluxioClusterEndpoint *AlluxioClusterEndpoint) show(request *restful.Request, response *restful.Response) {
	// Get AlluxioClusterList
	alluxioClusterList := new(v1alpha1.AlluxioClusterList)
	err := alluxioClusterEndpoint.client.List(request.Request.Context(), alluxioClusterList, &client.ListOptions{})
	// Unable to get:
	if err != nil {
		writeError(response, 404, Error{
			Title:   "Error",
			Details: fmt.Sprintf("Could not retrieve list: %s", err),
		})
		return
	} else {
		// Convert
		newAL := AlluxioClusterConverter.AlluxioClusterList(alluxioClusterList)
		// Send
		if err := response.WriteAsJson(newAL); err != nil {
			writeError(response, 404, Error{
				Title:   "Error",
				Details: "Could not list resources",
			})
		}
	}
}

func (alluxioClusterEndpoint *AlluxioClusterEndpoint) create(request *restful.Request, response *restful.Response) {
	// Read input
	alluxioClusterConfig := &AlluxioClusterConfig{}
	err := request.ReadEntity(alluxioClusterConfig)

	if err != nil {
		writeError(response, 400, Error{
			Title:   "Bad Request",
			Details: "Could not read entity. Check the Input File Format.",
		})
		return
	}

	// TODO May need validation func
	alluxioClusterObj := &v1alpha1.AlluxioCluster{
		ObjectMeta: metav1.ObjectMeta{Name: alluxioClusterConfig.Name, Namespace: "default"},
		Spec:       *alluxioClusterConfig.Spec,
	}

	//Deploy the object
	err = alluxioClusterEndpoint.client.Create(request.Request.Context(), alluxioClusterObj, &client.CreateOptions{})
	if err != nil {
		logger.Infof("Unable to create, the error is: %s ", err)
		writeError(response, 400, Error{
			Title:   "Error",
			Details: fmt.Sprintf("Could not create object: %s", err),
		})
	} else {
		// Convert
		newAlluxioClusterObj := AlluxioClusterConverter.AlluxioClusterObject(*alluxioClusterObj)
		// Send
		if err := response.WriteAsJson(newAlluxioClusterObj); err != nil {
			writeError(response, 404, Error{
				Title:   "Error",
				Details: "Could not list resources",
			})
		}
	}
	logger.Infof("Create Alluxio Cluster: %s", alluxioClusterObj.ObjectMeta.Name)
}

func (alluxioClusterEndpoint *AlluxioClusterEndpoint) delete(request *restful.Request, response *restful.Response) {
	// Read input
	alluxioClusterConfig := &AlluxioClusterConfig{}
	err := request.ReadEntity(alluxioClusterConfig)

	if err != nil {
		writeError(response, 400, Error{
			Title:   "Bad Request",
			Details: "Could not read entity",
		})
		return
	}

	// TODO May need validation func
	alluxioClusterObj := &v1alpha1.AlluxioCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      alluxioClusterConfig.Name,
			Namespace: "default",
		},
	}

	// Delete the object
	if err = alluxioClusterEndpoint.client.Delete(request.Request.Context(), alluxioClusterObj, &client.DeleteOptions{}); err != nil {
		writeError(response, 400, Error{
			Title:   "Error",
			Details: fmt.Sprintf("Could not Delete dataset: %s", err),
		})
	}
	logger.Infof("DELETE Alluxio Cluster: %s", alluxioClusterObj.ObjectMeta.Name)
}

//func (alluxioClusterEndpoint *AlluxioClusterEndpoint) update(request *restful.Request, response *restful.Response) {
//	writeError(response, 400, Error{
//		Title:   "Error",
//		Details: "To Do",
//	})
//}
