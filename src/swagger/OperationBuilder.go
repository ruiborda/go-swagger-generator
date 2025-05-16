package swagger

import (
	openapi2 "github.com/ruiborda/go-swagger-generator/src/openapi"
	entity2 "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
	"strconv"
)

type OperationBuilder struct {
	operation   *entity2.OperationEntity
	pathBuilder *PathItemBuilder
	docBuilder  *SwaggerDocBuilder
}

func (b *OperationBuilder) Summary(summary string) openapi2.Operation {
	b.operation.Summary = summary
	return b
}
func (b *OperationBuilder) Description(description string) openapi2.Operation {
	b.operation.Description = description
	return b
}
func (b *OperationBuilder) OperationID(id string) openapi2.Operation {
	b.operation.OperationID = id
	return b
}
func (b *OperationBuilder) Tag(tag string) openapi2.Operation {
	b.operation.Tags = append(b.operation.Tags, tag)
	return b
}
func (b *OperationBuilder) Tags(tags ...string) openapi2.Operation {
	b.operation.Tags = append(b.operation.Tags, tags...)
	return b
}
func (b *OperationBuilder) Consume(mimeType string) openapi2.Operation {
	b.operation.Consumes = append(b.operation.Consumes, mimeType)
	return b
}
func (b *OperationBuilder) Consumes(mimeTypes ...string) openapi2.Operation {
	b.operation.Consumes = append(b.operation.Consumes, mimeTypes...)
	return b
}
func (b *OperationBuilder) Produce(mimeType mime.MimeType) openapi2.Operation {
	b.operation.Produces = append(b.operation.Produces, mimeType)
	return b
}
func (b *OperationBuilder) Produces(mimeTypes ...mime.MimeType) openapi2.Operation {
	//b.operation.Produces = append(b.operation.Produces, mimeTypes...)
	b.operation.Produces = append(
		b.operation.Produces,
		mimeTypes...,
	)
	return b
}

func (b *OperationBuilder) Parameter(name, in string, config func(builder openapi2.Parameter)) openapi2.Operation {
	param := entity2.ParameterEntity{Name: name, In: in}
	if in == "path" {
		param.Required = true
	}
	paramBuilder := &ParameterBuilder{param: &param, docBuilder: b.docBuilder}
	config(paramBuilder)
	b.operation.Parameters = append(b.operation.Parameters, param)
	return b
}

func (b *OperationBuilder) QueryParameter(name string, config func(openapi2.Parameter)) openapi2.Operation {
	return b.Parameter(name, "query", config)
}
func (b *OperationBuilder) PathParameter(name string, config func(openapi2.Parameter)) openapi2.Operation {
	return b.Parameter(name, "path", func(pb openapi2.Parameter) {
		pb.Required(true) // Default for path parameters
		config(pb)
	})
}
func (b *OperationBuilder) HeaderParameter(name string, config func(openapi2.Parameter)) openapi2.Operation {
	return b.Parameter(name, "header", config)
}
func (b *OperationBuilder) FormParameter(name string, config func(openapi2.Parameter)) openapi2.Operation {
	return b.Parameter(name, "formData", config)
}
func (b *OperationBuilder) BodyParameter(config func(openapi2.Parameter)) openapi2.Operation {
	return b.Parameter("body", "body", config)
}

func (b *OperationBuilder) Response(statusCode int, config func(builder openapi2.Response)) openapi2.Operation {
	resp := entity2.ResponseEntity{}
	if b.operation.Responses == nil {
		b.operation.Responses = make(map[string]entity2.ResponseEntity)
	}
	responseBuilder := &ResponseBuilder{response: &resp, docBuilder: b.docBuilder}
	config(responseBuilder)
	b.operation.Responses[strconv.Itoa(statusCode)] = resp
	return b
}
func (b *OperationBuilder) Security(schemeName string, scopes ...string) openapi2.Operation {
	if b.operation.Security == nil {
		b.operation.Security = make([]map[string][]string, 0)
	}
	secRequirement := map[string][]string{schemeName: scopes}
	if len(scopes) == 0 { // Ensure empty array if no scopes, not nil
		secRequirement[schemeName] = []string{}
	}
	b.operation.Security = append(b.operation.Security, secRequirement)
	return b
}
func (b *OperationBuilder) Deprecated(deprecated bool) openapi2.Operation {
	b.operation.Deprecated = deprecated
	return b
}
func (b *OperationBuilder) ExternalDocumentation(url string, description string) openapi2.Operation {
	b.operation.ExternalDocs = &entity2.ExternalDocumentationEntity{URL: url, Description: description}
	return b
}
func (b *OperationBuilder) Path() openapi2.PathItem {
	return b.pathBuilder
}
