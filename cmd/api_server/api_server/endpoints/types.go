package endpoints

import (
	"errors"
	"fmt"
	"github.com/ajatprabha/operator-example/api/v1alpha1"
	"github.com/emicklei/go-restful"
	"strings"
)

type Endpoint interface {
	SetupWithWS(ws *restful.WebService)
}

type 	Source struct {
	Type       string `json:"type"`
	BaseURL    string `json:"baseUrl,omitempty"`
	AccessKey  string `json:"accessKey,omitempty"`
	SecretKey  string `json:"secretKey,omitempty"`
	Region     string `json:"region,omitempty"`
	PathPrefix string `json:"pathPrefix,omitempty"`
}

type Darkroom struct {
	Name       string                   `json:"name"`
	Version    string                   `json:"version"`
	Source     Source                   `json:"source"`
	SubDomains []string                 `json:"subDomains,omitempty"`
	Status     *v1alpha1.DarkroomStatus `json:"status,omitempty"`
}

func (d *Darkroom) Validate() error {
	var validated bool
	switch v1alpha1.Type(d.Source.Type) {
	case v1alpha1.WebFolder:
		validated = true
		if d.Source.BaseURL == "" {
			return errors.New("baseUrl can't be empty")
		}
	case v1alpha1.S3:
		validated = true
		var missingFields []string
		if d.Source.AccessKey == "" {
			missingFields = append(missingFields, "accessKey")
		}
		if d.Source.SecretKey == "" {
			missingFields = append(missingFields, "secretKey")
		}
		if d.Source.Region == "" {
			missingFields = append(missingFields, "region")
		}
		if len(missingFields) > 0 {
			return errors.New(fmt.Sprintf("required fields missing: %s", strings.Join(missingFields, ", ")))
		}
	}
	if validated {
		return nil
	}
	return errors.New("source type validation was not performed, type can only be [WebFolder,S3]")
}

type List struct {
	Items []Darkroom `json:"items"`
}

type Error struct {
	Title   string `json:"title"`
	Details string `json:"details"`
}
