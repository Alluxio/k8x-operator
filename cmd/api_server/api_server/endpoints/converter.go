package endpoints

import (
	"github.com/ajatprabha/operator-example/api/v1alpha1"
)

var From = &from{}

type from struct{}

func (c *from) Object(d *v1alpha1.Darkroom) *Darkroom {
	return &Darkroom{
		Name:    d.Name,
		Version: d.Spec.Version,
		Source: Source{
			Type:    string(d.Spec.Source.Type),
			BaseURL: d.Spec.Source.BaseURL,
		},
		Status: &d.Status,
	}
}

func (c *from) List(list *v1alpha1.DarkroomList) *List {
	items := make([]Darkroom, len(list.Items))
	for i, r := range list.Items {
		items[i] = *c.Object(&r)
	}
	return &List{
		Items: items,
	}
}
