package swagger

import (
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec"
)

type ContactBuilder struct {
	contact *openapi_spec.ContactEntity
}

func (b *ContactBuilder) Name(name string) openapi.Contact {
	b.contact.Name = name
	return b
}
func (b *ContactBuilder) URL(url string) openapi.Contact {
	b.contact.URL = url
	return b
}
func (b *ContactBuilder) Email(email string) openapi.Contact {
	b.contact.Email = email
	return b
}
