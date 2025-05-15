package swagger

import (
	openapi2 "github.com/ruiborda/go-swagger-generator/src/openapi"
	entity2 "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
)

type InfoBuilder struct {
	info *entity2.InfoEntity
}

func (b *InfoBuilder) Title(title string) openapi2.Info {
	b.info.Title = title
	return b
}
func (b *InfoBuilder) Version(version string) openapi2.Info {
	b.info.Version = version
	return b
}
func (b *InfoBuilder) Description(description string) openapi2.Info {
	b.info.Description = description
	return b
}
func (b *InfoBuilder) TermsOfService(terms string) openapi2.Info {
	b.info.TermsOfService = terms
	return b
}
func (b *InfoBuilder) Contact(config func(builder openapi2.Contact)) openapi2.Info {
	if b.info.Contact == nil {
		b.info.Contact = &entity2.ContactEntity{}
	}
	contactBuilder := &ContactBuilder{contact: b.info.Contact}
	config(contactBuilder)
	return b
}
func (b *InfoBuilder) License(config func(builder openapi2.License)) openapi2.Info {
	if b.info.License == nil {
		b.info.License = &entity2.LicenseEntity{}
	}
	licenseBuilder := &LicenseBuilder{license: b.info.License}
	config(licenseBuilder)
	return b
}
