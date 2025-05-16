package swagger

import (
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec"
)

type ServerBuilder struct {
	server *openapi_spec.ServerEntity
}

func (b *ServerBuilder) Description(description string) openapi.Server {
	b.server.Description = description
	return b
}
