package openapi

type Server interface {
	Description(description string) Server
	Variable(name string, defaultValue string, config func(ServerVariable)) Server
}

type ServerVariable interface {
	Enum(values ...string) ServerVariable
	Description(description string) ServerVariable
}
