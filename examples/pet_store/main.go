package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/src/middleware"
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/swagger"
)

func main() {
	router := gin.Default()

	// Set up Swagger
	router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json",
		UIPath:   "/",
	}))

	doc := swagger.Swagger()

	doc.Info(func(info openapi.Info) {
		info.Title("SwaggerGin Petstore").
			Version("1.0.7").
			Description("This is a sample server Petstore server. You can find out more about SwaggerGin at [http://swagger.io](http://swagger.io) or on [irc.freenode.net, #swagger](http://swagger.io/irc/). For this sample, you can use the api key `special-key` to test the authorization filters.").
			TermsOfService("http://swagger.io/terms/")
	}).
		Server("/", func(server openapi.Server) {
			server.Description("Servidor de desarrollo local")
		}).
		BasePath("/v2").
		Schemes("http", "https")

	doc.SecurityDefinition("api_key", func(sd openapi.SecurityScheme) {
		sd.Type("apiKey").Name("api_key").In("header")
	})
	doc.SecurityDefinition("petstore_auth", func(sd openapi.SecurityScheme) {
		sd.Type("oauth2").
			AuthorizationURL("https://petstore.swagger.io/oauth/authorize").
			Flow("implicit").
			Scope("read:pets", "read your pets").
			Scope("write:pets", "modify pets in your account")
	})

	doc.ExternalDocumentation("http://swagger.io", "Find out more about SwaggerGin")

	// Define DTOs and their schemas
	//_, _ = doc.DefinitionFromDTO(&controller.Pet{})
	//_, _ = doc.DefinitionFromDTO(&controller.Category{})
	//_, _ = doc.DefinitionFromDTO(&controller.Tag{})
	//_, _ = doc.DefinitionFromDTO(&controller.ApiResponse{})
	//_, _ = doc.DefinitionFromDTO(&controller.Order{})
	//_, _ = doc.DefinitionFromDTO(&controller.User{})

	setupRoutes(router)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("SwaggerGin UI available at http://localhost:8080/")
	fmt.Println("SwaggerGin JSON available at http://localhost:8080/openapi.json")
	_ = router.Run(":8080")
}
