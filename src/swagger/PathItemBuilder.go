package swagger

import (
	openapi "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
	entity "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec"
	"net/http"
	"strings"
)

type PathItemBuilder struct {
	pathItem   *entity.PathItem
	docPath    string // The path string, e.g., "/pets"
	docBuilder *SwaggerDocBuilder
}

func (b *PathItemBuilder) operation(method string, config func(openapi.Operation)) openapi.PathItem {
	op := &entity.Operation{
		Responses:  make(map[string]*entity.ResponseRef),
		Parameters: make([]*entity.ParameterRef, 0),
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
	case http.MethodTrace:
		b.pathItem.Trace = op
	}
	return b
}

func (b *PathItemBuilder) Get(config func(openapi.Operation)) openapi.PathItem {
	return b.operation(http.MethodGet, config)
}

func (b *PathItemBuilder) Post(config func(openapi.Operation)) openapi.PathItem {
	return b.operation(http.MethodPost, config)
}

func (b *PathItemBuilder) Put(config func(openapi.Operation)) openapi.PathItem {
	return b.operation(http.MethodPut, config)
}

func (b *PathItemBuilder) Delete(config func(openapi.Operation)) openapi.PathItem {
	return b.operation(http.MethodDelete, config)
}

func (b *PathItemBuilder) Options(config func(openapi.Operation)) openapi.PathItem {
	return b.operation(http.MethodOptions, config)
}

func (b *PathItemBuilder) Head(config func(openapi.Operation)) openapi.PathItem {
	return b.operation(http.MethodHead, config)
}

func (b *PathItemBuilder) Patch(config func(openapi.Operation)) openapi.PathItem {
	return b.operation(http.MethodPatch, config)
}

// Parameter adds a parameter common to all operations in this path.
func (b *PathItemBuilder) Parameter(name, in string, config func(openapi.Parameter)) openapi.PathItem {
	param := entity.Parameter{Name: name, In: in}
	if in == "path" {
		param.Required = true
	}
	paramBuilder := &ParameterBuilder{param: &param, docBuilder: b.docBuilder}
	config(paramBuilder)

	if b.pathItem.Parameters == nil {
		b.pathItem.Parameters = make([]*entity.ParameterRef, 0)
	}
	b.pathItem.Parameters = append(b.pathItem.Parameters, &entity.ParameterRef{Parameter: &param})
	return b
}

func (b *PathItemBuilder) Doc() openapi.SwaggerDocBuilder {
	return b.docBuilder
}
