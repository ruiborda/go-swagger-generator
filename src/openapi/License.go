package openapi

type License interface {
	Name(name string) License
	URL(url string) License
}
