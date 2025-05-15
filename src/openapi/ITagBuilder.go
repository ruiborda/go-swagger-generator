package openapi

type Tag interface {
	Description(description string) Tag
	ExternalDocumentation(url string, description string) Tag
}
