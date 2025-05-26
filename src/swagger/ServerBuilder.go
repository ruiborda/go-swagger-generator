package swagger

import (
	openapi "github.com/ruiborda/go-swagger-generator/src/openapi"
	entity "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
)

type ServerBuilder struct {
	server *entity.Server
}

func (b *ServerBuilder) Description(description string) openapi.Server {
	b.server.Description = description
	return b
}

func (b *ServerBuilder) Variable(name string, defaultValue string, config func(openapi.ServerVariable)) openapi.Server {
	if b.server.Variables == nil {
		b.server.Variables = make(map[string]*entity.ServerVariable)
	}
	variable := entity.ServerVariable{Default: defaultValue}
	variableBuilder := &ServerVariableBuilder{variable: &variable}
	config(variableBuilder)
	b.server.Variables[name] = &variable
	return b
}

type ServerVariableBuilder struct {
	variable *entity.ServerVariable
}

func (b *ServerVariableBuilder) Enum(values ...string) openapi.ServerVariable {
	b.variable.Enum = append(b.variable.Enum, values...)
	return b
}

func (b *ServerVariableBuilder) Description(description string) openapi.ServerVariable {
	b.variable.Description = description
	return b
}
