package swagger

import (
	openapi2 "github.com/ruiborda/go-swagger-generator/src/openapi"
	entity2 "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
	"net/http"
	"strings"
)

type PathItemBuilder struct {
	pathItem   *entity2.PathItemEntity
	docPath    string
	docBuilder *SwaggerDocBuilder
}

func (b *PathItemBuilder) operation(method string, config func(builder openapi2.Operation)) openapi2.PathItem {
	op := &entity2.OperationEntity{
		Responses:  make(map[string]entity2.ResponseEntity),
		Parameters: make([]entity2.ParameterEntity, 0),
	}
	opBuilder := &OperationBuilder{operation: op, pathBuilder: b, docBuilder: b.docBuilder}
	config(opBuilder)

	switch strings.ToUpper(method) {
	case http.MethodGet:
		b.pathItem.Get = op
	case http.MethodPost:
		b.pathItem.Post = op
	case http.MethodPut:
		b.pathItem.Put = op
	case http.MethodDelete:
		b.pathItem.Delete = op
	case http.MethodOptions:
		b.pathItem.Options = op
	case http.MethodHead:
		b.pathItem.Head = op
	case http.MethodPatch:
		b.pathItem.Patch = op
	}
	b.docBuilder.doc.Paths[b.docPath] = *b.pathItem
	return b
}

func (b *PathItemBuilder) Get(config func(openapi2.Operation)) openapi2.PathItem {
	return b.operation(http.MethodGet, config)
}
func (b *PathItemBuilder) Post(config func(openapi2.Operation)) openapi2.PathItem {
	return b.operation(http.MethodPost, config)
}
func (b *PathItemBuilder) Put(config func(openapi2.Operation)) openapi2.PathItem {
	return b.operation(http.MethodPut, config)
}
func (b *PathItemBuilder) Delete(config func(openapi2.Operation)) openapi2.PathItem {
	return b.operation(http.MethodDelete, config)
}
func (b *PathItemBuilder) Options(config func(openapi2.Operation)) openapi2.PathItem {
	return b.operation(http.MethodOptions, config)
}
func (b *PathItemBuilder) Head(config func(openapi2.Operation)) openapi2.PathItem {
	return b.operation(http.MethodHead, config)
}
func (b *PathItemBuilder) Patch(config func(openapi2.Operation)) openapi2.PathItem {
	return b.operation(http.MethodPatch, config)
}

func (b *PathItemBuilder) Parameter(name, in string, config func(builder openapi2.Parameter)) openapi2.PathItem {
	param := entity2.ParameterEntity{Name: name, In: in}
	if in == "path" {
		param.Required = true
	}
	paramBuilder := &ParameterBuilder{param: &param, docBuilder: b.docBuilder}
	config(paramBuilder)
	b.pathItem.Parameters = append(b.pathItem.Parameters, param)
	b.docBuilder.doc.Paths[b.docPath] = *b.pathItem
	return b
}

func (b *PathItemBuilder) Doc() openapi2.SwaggerDoc {
	b.docBuilder.doc.Paths[b.docPath] = *b.pathItem
	return b.docBuilder
}
