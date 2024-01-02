package endpoints

import (
	"github.com/Alluxio/k8s-operator/api/v1alpha1"
	"github.com/emicklei/go-restful"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubelog "sigs.k8s.io/controller-runtime/pkg/log"
)

// writeError is a helper func that write error message to rest response
func writeError(response *restful.Response, httpStatus int, err Error) {
	if err := response.WriteHeaderAndJson(httpStatus, err, "application/json"); err != nil {
		kubelog.Log.Error(err, "Could not write the error response")
	}
}

// DatasetConverter is used to simplify Dataset CRD to struct in types.go
var AlluxioClusterConverter = &alluxioClusterConverter{}

type alluxioClusterConverter struct{}

func (c *alluxioClusterConverter) AlluxioClusterObject(d *v1alpha1.AlluxioCluster) *AlluxioCluster {
	return &AlluxioCluster{
		AlluxioClusterConfig: AlluxioClusterConfig{
			Name: d.Name,
			Spec: &d.Spec,
		},
		Status: &d.Status,
	}
}

func (c *alluxioClusterConverter) AlluxioClusterList(list *v1alpha1.AlluxioClusterList) *AlluxioClusterList {
	items := make([]AlluxioCluster, len(list.Items))
	for i, r := range list.Items {
		items[i] = *c.AlluxioClusterObject(&r)
	}
	return &AlluxioClusterList{
		AlluxioClusters: items,
	}
}

// DatasetConverter is used to simplify Dataset CRD to struct in types.go
var DatasetConverter = &datasetConverter{}

type datasetConverter struct{}

func (c *datasetConverter) DatasetObject(d *v1alpha1.Dataset) *Dataset {
	return &Dataset{
		DatasetConfig: DatasetConfig{
			Name:        d.Name,
			Path:        d.Spec.Dataset.Path,
			Credentials: d.Spec.Dataset.Credentials,
		},
		Status: &d.Status,
	}
}

func (c *datasetConverter) DatasetList(list *v1alpha1.DatasetList) *DatasetList {
	datasets := make([]Dataset, len(list.Items))
	for i, r := range list.Items {
		datasets[i] = *c.DatasetObject(&r)
	}
	return &DatasetList{
		Datasets: datasets,
	}
}

func (c *datasetConverter) K8SDatasetObject(datasetConfig DatasetConfig) *v1alpha1.Dataset {
	return &v1alpha1.Dataset{
		ObjectMeta: metav1.ObjectMeta{Name: datasetConfig.Name, Namespace: "default"},
		Spec: v1alpha1.DatasetSpec{
			Dataset: v1alpha1.DatasetConf{
				Path:        datasetConfig.Path,
				Credentials: datasetConfig.Credentials,
			},
		},
	}
}
