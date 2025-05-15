package openapi

type Contact interface {
	Name(name string) Contact
	URL(url string) Contact
	Email(email string) Contact
}
