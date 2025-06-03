package swagger

import (
	openapi "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
	entity "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec"
)

type SecuritySchemeBuilder struct {
	scheme *entity.SecurityScheme
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

func (b *SecuritySchemeBuilder) Scheme(scheme string) openapi.SecurityScheme {
	b.scheme.Scheme = scheme
	return b
}

func (b *SecuritySchemeBuilder) BearerFormat(format string) openapi.SecurityScheme {
	b.scheme.BearerFormat = format
	return b
}

func (b *SecuritySchemeBuilder) Flows(config func(openapi.OAuthFlows)) openapi.SecurityScheme {
	if b.scheme.Flows == nil {
		b.scheme.Flows = &entity.OAuthFlows{}
	}
	flowsBuilder := &OAuthFlowsBuilder{flows: b.scheme.Flows}
	config(flowsBuilder)
	return b
}

func (b *SecuritySchemeBuilder) OpenIDConnectURL(url string) openapi.SecurityScheme {
	b.scheme.OpenIdConnectUrl = url
	return b
}
func (b *SecuritySchemeBuilder) Ref(ref string) openapi.SecurityScheme {
	// Similar to other builders, this would mean the entire SecurityScheme is a ref.
	// This builder builds entity.SecurityScheme, not entity.SecuritySchemeRef.
	// Handled by ComponentSecurityScheme in SwaggerDocBuilder when a ref is provided.
	return b
}

type OAuthFlowsBuilder struct {
	flows *entity.OAuthFlows
}

func (b *OAuthFlowsBuilder) Implicit(config func(openapi.OAuthFlow)) openapi.OAuthFlows {
	if b.flows.Implicit == nil {
		b.flows.Implicit = &entity.OAuthFlow{}
	}
	flowBuilder := &OAuthFlowBuilder{flow: b.flows.Implicit}
	config(flowBuilder)
	return b
}

func (b *OAuthFlowsBuilder) Password(config func(openapi.OAuthFlow)) openapi.OAuthFlows {
	if b.flows.Password == nil {
		b.flows.Password = &entity.OAuthFlow{}
	}
	flowBuilder := &OAuthFlowBuilder{flow: b.flows.Password}
	config(flowBuilder)
	return b
}

func (b *OAuthFlowsBuilder) ClientCredentials(config func(openapi.OAuthFlow)) openapi.OAuthFlows {
	if b.flows.ClientCredentials == nil {
		b.flows.ClientCredentials = &entity.OAuthFlow{}
	}
	flowBuilder := &OAuthFlowBuilder{flow: b.flows.ClientCredentials}
	config(flowBuilder)
	return b
}

func (b *OAuthFlowsBuilder) AuthorizationCode(config func(openapi.OAuthFlow)) openapi.OAuthFlows {
	if b.flows.AuthorizationCode == nil {
		b.flows.AuthorizationCode = &entity.OAuthFlow{}
	}
	flowBuilder := &OAuthFlowBuilder{flow: b.flows.AuthorizationCode}
	config(flowBuilder)
	return b
}

type OAuthFlowBuilder struct {
	flow *entity.OAuthFlow
}

func (b *OAuthFlowBuilder) AuthorizationURL(url string) openapi.OAuthFlow {
	b.flow.AuthorizationUrl = url
	return b
}

func (b *OAuthFlowBuilder) TokenURL(url string) openapi.OAuthFlow {
	b.flow.TokenUrl = url
	return b
}

func (b *OAuthFlowBuilder) RefreshURL(url string) openapi.OAuthFlow {
	b.flow.RefreshUrl = url
	return b
}

func (b *OAuthFlowBuilder) Scope(scopeName string, description string) openapi.OAuthFlow {
	if b.flow.Scopes == nil {
		b.flow.Scopes = make(map[string]string)
	}
	b.flow.Scopes[scopeName] = description
	return b
}

func (b *OAuthFlowBuilder) Scopes(scopes map[string]string) openapi.OAuthFlow {
	b.flow.Scopes = scopes
	return b
}
