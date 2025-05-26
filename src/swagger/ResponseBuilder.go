package swagger

import (
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	entity "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
)

type ResponseBuilder struct {
	response   *entity.Response
	docBuilder *SwaggerDocBuilder
}

func (b *ResponseBuilder) Description(description string) openapi.Response {
	b.response.Description = description
	return b
}

func (b *ResponseBuilder) Header(name string, config func(openapi.Header)) openapi.Response {
	if b.response.Headers == nil {
		b.response.Headers = make(map[string]*entity.HeaderRef)
	}
	header := entity.Header{}
	headerBuilder := &HeaderBuilder{header: &header, docBuilder: b.docBuilder}
	config(headerBuilder)
	b.response.Headers[name] = &entity.HeaderRef{Header: &header}
	return b
}

func (b *ResponseBuilder) Content(mimeTypeStr mime.MimeType, config func(openapi.MediaType)) openapi.Response {
	if b.response.Content == nil {
		b.response.Content = make(map[string]*entity.MediaType)
	}
	mediaType := entity.MediaType{}
	mediaTypeBuilder := &MediaTypeBuilder{mediaType: &mediaType, docBuilder: b.docBuilder}
	config(mediaTypeBuilder)
	b.response.Content[string(mimeTypeStr)] = &mediaType
	return b
}

func (b *ResponseBuilder) Link(name string, config func(openapi.Link)) openapi.Response {
	if b.response.Links == nil {
		b.response.Links = make(map[string]*entity.LinkRef)
	}
	link := entity.Link{}
	linkBuilder := &LinkBuilder{link: &link}
	config(linkBuilder)
	b.response.Links[name] = &entity.LinkRef{Link: &link}
	return b
}

func (b *ResponseBuilder) Ref(ref string) openapi.Response {
	return b
}
