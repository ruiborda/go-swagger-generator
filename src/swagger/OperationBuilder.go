package swagger

import (
	openapi "github.com/ruiborda/go-swagger-generator/src/openapi"
	entity "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
	"strconv"
)

type OperationBuilder struct {
	operation   *entity.Operation
	pathBuilder *PathItemBuilder
	docBuilder  *SwaggerDocBuilder
}

func (b *OperationBuilder) Summary(summary string) openapi.Operation {
	b.operation.Summary = summary
	return b
}

func (b *OperationBuilder) Description(description string) openapi.Operation {
	b.operation.Description = description
	return b
}

func (b *OperationBuilder) OperationID(id string) openapi.Operation {
	b.operation.OperationID = id
	return b
}

func (b *OperationBuilder) Tag(tag string) openapi.Operation {
	b.operation.Tags = append(b.operation.Tags, tag)
	return b
}

func (b *OperationBuilder) Tags(tags ...string) openapi.Operation {
	b.operation.Tags = append(b.operation.Tags, tags...)
	return b
}

func (b *OperationBuilder) Parameter(name, in string, config func(openapi.Parameter)) openapi.Operation {
	param := entity.Parameter{Name: name, In: in}
	if in == "path" {
		param.Required = true
	}
	paramBuilder := &ParameterBuilder{param: &param, docBuilder: b.docBuilder}
	config(paramBuilder)

	// Parameters are stored as ParameterRef in Operation
	if b.operation.Parameters == nil {
		b.operation.Parameters = make([]*entity.ParameterRef, 0)
	}
	b.operation.Parameters = append(b.operation.Parameters, &entity.ParameterRef{Parameter: &param})
	return b
}

func (b *OperationBuilder) QueryParameter(name string, config func(openapi.Parameter)) openapi.Operation {
	return b.Parameter(name, "query", config)
}

func (b *OperationBuilder) PathParameter(name string, config func(openapi.Parameter)) openapi.Operation {
	return b.Parameter(name, "path", func(pb openapi.Parameter) {
		pb.Required(true) // Default for path parameters, also set in Parameter() above but good to be explicit.
		config(pb)
	})
}

func (b *OperationBuilder) HeaderParameter(name string, config func(openapi.Parameter)) openapi.Operation {
	return b.Parameter(name, "header", config)
}

func (b *OperationBuilder) CookieParameter(name string, config func(openapi.Parameter)) openapi.Operation {
	return b.Parameter(name, "cookie", config)
}

func (b *OperationBuilder) RequestBody(config func(openapi.RequestBody)) openapi.Operation {
	if b.operation.RequestBody == nil {
		b.operation.RequestBody = &entity.RequestBodyRef{RequestBody: &entity.RequestBody{}}
	}
	if b.operation.RequestBody.RequestBody == nil { // If it was a Ref, ensure RequestBody part is initialized
		b.operation.RequestBody.RequestBody = &entity.RequestBody{}
	}

	// Ensure Content map is initialized
	if b.operation.RequestBody.RequestBody.Content == nil {
		b.operation.RequestBody.RequestBody.Content = make(map[string]*entity.MediaType)
	}

	reqBodyBuilder := &RequestBodyBuilder{requestBody: b.operation.RequestBody.RequestBody, docBuilder: b.docBuilder}
	config(reqBodyBuilder)
	return b
}

func (b *OperationBuilder) Response(statusCode int, config func(openapi.Response)) openapi.Operation {
	return b.responseInternal(strconv.Itoa(statusCode), config)
}

func (b *OperationBuilder) DefaultResponse(config func(openapi.Response)) openapi.Operation {
	return b.responseInternal("default", config)
}

func (b *OperationBuilder) responseInternal(statusKey string, config func(openapi.Response)) openapi.Operation {
	if b.operation.Responses == nil {
		b.operation.Responses = make(map[string]*entity.ResponseRef)
	}
	resp := entity.Response{}
	responseBuilder := &ResponseBuilder{response: &resp, docBuilder: b.docBuilder}
	config(responseBuilder)
	b.operation.Responses[statusKey] = &entity.ResponseRef{Response: &resp}
	return b
}

func (b *OperationBuilder) Security(schemeName string, scopes ...string) openapi.Operation {
	if b.operation.Security == nil {
		b.operation.Security = make([]entity.SecurityRequirement, 0)
	}
	secRequirement := entity.SecurityRequirement{schemeName: scopes}
	if len(scopes) == 0 {
		secRequirement[schemeName] = []string{} // Ensure empty array if no scopes
	}
	b.operation.Security = append(b.operation.Security, secRequirement)
	return b
}

func (b *OperationBuilder) Deprecated(deprecated bool) openapi.Operation {
	b.operation.Deprecated = deprecated
	return b
}

func (b *OperationBuilder) ExternalDocumentation(url string, description string) openapi.Operation {
	b.operation.ExternalDocs = &entity.ExternalDocumentation{URL: url, Description: description}
	return b
}

func (b *OperationBuilder) Server(url string, config func(openapi.Server)) openapi.Operation {
	srv := entity.Server{URL: url}
	serverBuilder := &ServerBuilder{server: &srv}
	config(serverBuilder)
	b.operation.Servers = append(b.operation.Servers, &srv)
	return b
}

func (b *OperationBuilder) Path() openapi.PathItem {
	return b.pathBuilder
}
