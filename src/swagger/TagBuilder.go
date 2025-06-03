package swagger

import (
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi"
	entity2 "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec"
)

type TagBuilder struct {
	tag *entity2.Tag
}

func (b *TagBuilder) Description(description string) openapi.Tag {
	b.tag.Description = description
	return b
}
func (b *TagBuilder) ExternalDocumentation(url string, description string) openapi.Tag {
	b.tag.ExternalDocs = &entity2.ExternalDocumentation{URL: url, Description: description}
	return b
}
