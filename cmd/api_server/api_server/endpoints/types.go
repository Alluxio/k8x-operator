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

// AlluxioCluster Part
type AlluxioCluster struct {
	AlluxioClusterConfig AlluxioClusterConfig           `json:"alluxio-cluster-config"`
	Status               *v1alpha1.AlluxioClusterStatus `json:"status,omitempty"`
}

type AlluxioClusterConfig struct {
	Name string                       `json:"name"`
	Spec *v1alpha1.AlluxioClusterSpec `json:"spec,omitempty"`
}

type AlluxioClusterList struct {
	AlluxioClusters []AlluxioCluster `json:"alluxio-clusters"`
}

// Dataset Part
type Dataset struct {
	DatasetConfig DatasetConfig           `json:"dataset-config"`
	Status        *v1alpha1.DatasetStatus `json:"status,omitempty"`
}

type DatasetConfig struct {
	Name string                `json:"name"`
	Spec *v1alpha1.DatasetConf `json:"spec,omitempty"`
}

type DatasetList struct {
	Datasets []Dataset `json:"datasets"`
}
