package swagger

import (
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi"
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec"
)

type LicenseBuilder struct {
	license *openapi_spec.License
}

func (b *LicenseBuilder) Name(name string) openapi.License {
	b.license.Name = name
	return b
}
func (b *LicenseBuilder) URL(url string) openapi.License {
	b.license.URL = url
	return b
}
