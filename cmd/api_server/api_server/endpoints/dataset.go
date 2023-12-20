package endpoints

import (
	"fmt"
	"github.com/Alluxio/k8s-operator/api/v1alpha1"
	"github.com/emicklei/go-restful"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DatasetEndpoint struct {
	client client.Client
}

func NewDataSetEndpoint(client client.Client) *DatasetEndpoint {
	return &DatasetEndpoint{client: client}
}

func (datasetEndpoint *DatasetEndpoint) SetupWithWS(ws *restful.WebService) {
	ws.Route(ws.GET("dataset").To(datasetEndpoint.show).
		Doc("List of Datasets").
		Returns(200, "OK", DatasetList{}))

	ws.Route(ws.POST("dataset").To(datasetEndpoint.create).
		Doc("Create a new Dataset").
		Returns(200, "OK", Dataset{}).
		Returns(400, "Bad Request", nil))

	ws.Route(ws.DELETE("dataset").To(datasetEndpoint.delete).
		Doc("Delete Current Dataset").
		Returns(200, "OK", nil).
		Returns(400, "Bad Request", nil))

	// todo add more route
}

func (datasetEndpoint *DatasetEndpoint) show(request *restful.Request, response *restful.Response) {
	datasetList := new(v1alpha1.DatasetList)
	err := datasetEndpoint.client.List(request.Request.Context(), datasetList, &client.ListOptions{})

	if err != nil {
		writeError(response, 404, Error{
			Title:   "Error",
			Details: fmt.Sprintf("Could not retrieve list: %s", err),
		})
	} else {
		// Convert
		newDL := DatasetConverter.DatasetList(datasetList)
		// Send
		if err := response.WriteAsJson(newDL); err != nil {
			writeError(response, 404, Error{
				Title:   "Error",
				Details: "Could not list resources",
			})
		}
	}

}

func (datasetEndpoint *DatasetEndpoint) create(request *restful.Request, response *restful.Response) {
	// Read input
	dataset := new(Dataset)
	err := request.ReadEntity(dataset)

	if err != nil {
		writeError(response, 400, Error{
			Title:   "Bad Request",
			Details: "Could not read entity",
		})
		return
	}

	// TODO May need validation func

	datasetObj := &v1alpha1.Dataset{
		ObjectMeta: metav1.ObjectMeta{Name: dataset.Name},
		Spec: v1alpha1.DatasetSpec{
			Dataset: v1alpha1.DatasetConf{
				Path:        dataset.Path,
				Credentials: dataset.Credentials,
			},
		},
	}

	// Deploy the object
	err = datasetEndpoint.client.Create(request.Request.Context(), datasetObj, &client.CreateOptions{})
	if err != nil {
		writeError(response, 400, Error{
			Title:   "Error",
			Details: fmt.Sprintf("Could not create object: %s", err),
		})
	} else {
		// Convert
		newDatasetObj := DatasetConverter.DatasetObject(datasetObj)
		// Send
		if err := response.WriteAsJson(newDatasetObj); err != nil {
			writeError(response, 404, Error{
				Title:   "Error",
				Details: "Could not list resources",
			})
		}
	}
}

func (datasetEndpoint *DatasetEndpoint) delete(request *restful.Request, response *restful.Response) {
	//writeError(response, 404, Error{
	//	Title:   "Error",
	//	Details: "To Do",
	//})
	dataset := new(Dataset)
	err := request.ReadEntity(dataset)
	if err != nil {
		writeError(response, 400, Error{
			Title:   "Bad Request",
			Details: "Could not read entity",
		})
		return
	}

	// TODO May need validation func

	datasetObj := &v1alpha1.Dataset{
		ObjectMeta: metav1.ObjectMeta{Name: dataset.Name},
		Spec: v1alpha1.DatasetSpec{
			Dataset: v1alpha1.DatasetConf{
				Path:        dataset.Path,
				Credentials: dataset.Credentials,
			},
		},
	}

	// Delete the object
	if err = datasetEndpoint.client.Delete(request.Request.Context(), datasetObj, &client.DeleteOptions{}); err != nil {
		writeError(response, 400, Error{
			Title:   "Error",
			Details: fmt.Sprintf("Could not Delete object: %s", err),
		})
	}
}
