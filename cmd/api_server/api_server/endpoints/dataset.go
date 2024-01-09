package endpoints

import (
	"fmt"
	"github.com/Alluxio/k8s-operator/api/v1alpha1"
	"github.com/Alluxio/k8s-operator/pkg/logger"
	"github.com/emicklei/go-restful"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DatasetEndpoint struct {
	client client.Client
}

func NewDataSetEndpoint(client client.Client) *DatasetEndpoint {
	return &DatasetEndpoint{client: client}
}

func (datasetEndpoint *DatasetEndpoint) SetupWithWS(ws *restful.WebService) {
	ws.Route(ws.GET("dataset").To(datasetEndpoint.show).
		Doc("List of All Datasets").
		Returns(200, "OK", DatasetList{}))

	ws.Route(ws.POST("dataset").To(datasetEndpoint.create).
		Doc("Create a new Dataset").
		Returns(200, "OK", Dataset{}).
		Returns(400, "Bad Request", nil))

	ws.Route(ws.DELETE("dataset").To(datasetEndpoint.delete).
		Doc("Delete a Current Dataset").
		Returns(200, "OK", nil).
		Returns(400, "Bad Request", nil))

	//ws.Route(ws.PUT("dataset").To(datasetEndpoint.update).
	//	Doc("Update a Current Dataset").
	//	Returns(200, "OK", Dataset{}).
	//	Returns(400, "Bad Request", nil))
	//
}

// show takes GET request.
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

// create takes POST request with json format datasetConfig types.
func (datasetEndpoint *DatasetEndpoint) create(request *restful.Request, response *restful.Response) {
	// Read input
	datasetConfig := &DatasetConfig{}
	err := request.ReadEntity(datasetConfig)

	if err != nil {
		writeError(response, 400, Error{
			Title:   "Bad Request",
			Details: "Could not read entity. Check the Input File Format.",
		})
		return
	}

	// TODO May need validation func
	datasetObj := &v1alpha1.Dataset{
		ObjectMeta: metav1.ObjectMeta{Name: datasetConfig.Name, Namespace: "default"},
		Spec: v1alpha1.DatasetSpec{
			Dataset: v1alpha1.DatasetConf{
				Path:        datasetConfig.Spec.Path,
				Credentials: datasetConfig.Spec.Credentials,
			},
		},
	}

	//Deploy the object
	err = datasetEndpoint.client.Create(request.Request.Context(), datasetObj, &client.CreateOptions{})
	if err != nil {
		logger.Infof("Unable to create, the error is: %s ", err)
		writeError(response, 400, Error{
			Title:   "Error",
			Details: fmt.Sprintf("Could not create object: %s", err),
		})
	} else {
		// Convert
		newDatasetObj := DatasetConverter.DatasetObject(*datasetObj)
		// Send
		if err := response.WriteAsJson(newDatasetObj); err != nil {
			writeError(response, 404, Error{
				Title:   "Error",
				Details: "Could not list resources",
			})
		}
	}
	logger.Infof("Create Dataset: %s", datasetObj.ObjectMeta.Name)
}

// delete takes DELETE request
func (datasetEndpoint *DatasetEndpoint) delete(request *restful.Request, response *restful.Response) {
	// Read input
	datasetConfig := &DatasetConfig{}
	err := request.ReadEntity(datasetConfig)

	if err != nil {
		writeError(response, 400, Error{
			Title:   "Bad Request",
			Details: "Could not read entity",
		})
		return
	}

	datasetObj := &v1alpha1.Dataset{
		ObjectMeta: metav1.ObjectMeta{Name: datasetConfig.Name, Namespace: "default"},
	}

	// Delete the object
	if err = datasetEndpoint.client.Delete(request.Request.Context(), datasetObj, &client.DeleteOptions{}); err != nil {
		logger.Infof("DELETE Dataset: %s Has Error", datasetObj.ObjectMeta.Name)
		writeError(response, 400, Error{
			Title:   "Error",
			Details: fmt.Sprintf("Could not Delete dataset: %s", err),
		})
	}
	logger.Infof("DELETE Dataset: %s", datasetObj.ObjectMeta.Name)
}

//func (datasetEndpoint *DatasetEndpoint) update(request *restful.Request, response *restful.Response) {
//	writeError(response, 400, Error{
//		Title:   "Error",
//		Details: "To Do, can read entity",
//	})
//}
