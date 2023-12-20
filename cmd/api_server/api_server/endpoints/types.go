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
	Spec   *v1alpha1.AlluxioClusterSpec   `json:"spec,omitempty"`
	Status *v1alpha1.AlluxioClusterStatus `json:"status,omitempty"`
}

type AlluxioClusterList struct {
	Items []AlluxioCluster `json:"items"`
}

// Dataset Part
type Dataset struct {
	Name        string                  `json:"name"`
	Path        string                  `json:"path"`
	Credentials map[string]string       `json:"credentials,omitempty"`
	Status      *v1alpha1.DatasetStatus `json:"status,omitempty"`
}

type DatasetList struct {
	Items []Dataset `json:"items"`
}
