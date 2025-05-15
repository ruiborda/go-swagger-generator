package swagger

import (
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec"
)

type SecuritySchemeBuilder struct {
	scheme *openapi_spec.SecuritySchemeEntity
}

func (b *SecuritySchemeBuilder) Type(secType string) openapi.SecurityScheme {
	b.scheme.Type = secType
	return b
}
func (b *SecuritySchemeBuilder) Description(description string) openapi.SecurityScheme {
	b.scheme.Description = description
	return b
}
func (b *SecuritySchemeBuilder) Name(name string) openapi.SecurityScheme {
	b.scheme.Name = name
	return b
}
func (b *SecuritySchemeBuilder) In(in string) openapi.SecurityScheme {
	b.scheme.In = in
	return b
}
func (b *SecuritySchemeBuilder) Flow(flow string) openapi.SecurityScheme {
	b.scheme.Flow = flow
	return b
}
func (b *SecuritySchemeBuilder) AuthorizationURL(url string) openapi.SecurityScheme {
	b.scheme.AuthorizationURL = url
	return b
}
func (b *SecuritySchemeBuilder) TokenURL(url string) openapi.SecurityScheme {
	b.scheme.TokenURL = url
	return b
}
func (b *SecuritySchemeBuilder) Scope(scopeName string, description string) openapi.SecurityScheme {
	if b.scheme.Scopes == nil {
		b.scheme.Scopes = make(map[string]string)
	}
	b.scheme.Scopes[scopeName] = description
	return b
}
