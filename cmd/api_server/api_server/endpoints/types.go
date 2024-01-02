package endpoints

import (
	"github.com/Alluxio/k8s-operator/api/v1alpha1"
	"github.com/emicklei/go-restful"
)

type Endpoint interface {
	SetupWithWS(ws *restful.WebService)
}

type Source struct {
	Type       string `json:"type"`
	BaseURL    string `json:"baseUrl,omitempty"`
	AccessKey  string `json:"accessKey,omitempty"`
	SecretKey  string `json:"secretKey,omitempty"`
	Region     string `json:"region,omitempty"`
	PathPrefix string `json:"pathPrefix,omitempty"`
}

type Error struct {
	Title   string `json:"title"`
	Details string `json:"details"`
}

// AlluxioCluster Part todo check below!!!
type AlluxioClusterConfig struct {
	Name string                       `json:"name"`
	Spec *v1alpha1.AlluxioClusterSpec `json:"spec,omitempty"`
}

type AlluxioCluster struct {
	AlluxioClusterConfig AlluxioClusterConfig           `json:"alluxio-cluster-config"`
	Status               *v1alpha1.AlluxioClusterStatus `json:"status,omitempty"`
}

type AlluxioClusterList struct {
	AlluxioClusters []AlluxioCluster `json:"alluxio-clusters"`
}

type DatasetConfig struct {
	Name        string            `json:"name"`
	Path        string            `json:"path"`
	Credentials map[string]string `json:"credentials,omitempty"`
}

// Dataset todo change the `json:"datasetConfig"` captial letter
type Dataset struct {
	DatasetConfig DatasetConfig           `json:"datasetConfig"`
	Status        *v1alpha1.DatasetStatus `json:"status,omitempty"`
}

type DatasetList struct {
	Datasets []Dataset `json:"datasets"`
}
