package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/examples/pet_store/controller"
	"github.com/ruiborda/go-swagger-generator/src/middleware"
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/swagger"
)

func main() {
	router := gin.Default()

	// Set up Swagger
	router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json", // Path for the OAS3 JSON file
		UIPath:   "/",             // Path for the Swagger UI
		Title:    "Swagger Petstore OAS3",
	}))

	doc := swagger.Swagger() // This will use the OAS3 builder

	// Register DTOs so they are available in #/components/schemas/
	// The SchemaFromDTO method in builders will automatically do this if not already present.
	// It's good practice to ensure they are registered if you know you'll reference them.
	_, _ = doc.SchemaFromDTO(&controller.Pet{})
	_, _ = doc.SchemaFromDTO(&controller.Category{})
	_, _ = doc.SchemaFromDTO(&controller.Tag{})
	_, _ = doc.SchemaFromDTO(&controller.ApiResponse{})
	_, _ = doc.SchemaFromDTO(&controller.Order{})
	_, _ = doc.SchemaFromDTO(&controller.User{})

	doc.Info(func(info openapi.Info) {
		info.Title("Swagger Petstore - OpenAPI 3.0").
			Version("1.0.7").
			Description("This is a sample Pet Store Server based on the OpenAPI 3.0 specification. You can find out more about Swagger at [https://swagger.io](https://swagger.io). For this sample, you can use the api key `special-key` to test the authorization filters.").
			TermsOfService("http://swagger.io/terms/").
			License(func(license openapi.License) {
				license.Name("Apache 2.0").URL("http://www.apache.org/licenses/LICENSE-2.0.html")
			})
	}).
		Server("http://localhost:8080/v2", func(server openapi.Server) { // Server URL includes base path
			server.Description("Local development server")
		}).
		Server("https://petstore.swagger.io/v2", func(server openapi.Server) {
			server.Description("Production server")
		})

	doc.ComponentSecurityScheme("api_key", func(ss openapi.SecurityScheme) {
		ss.Type("apiKey").Name("api_key").In("header")
	})
	doc.ComponentSecurityScheme("petstore_auth", func(ss openapi.SecurityScheme) {
		ss.Type("oauth2").
			Description("OAuth2 implicit grant").
			Flows(func(flows openapi.OAuthFlows) {
				flows.Implicit(func(flow openapi.OAuthFlow) {
					flow.AuthorizationURL("https://petstore.swagger.io/oauth/authorize").
						Scope("write:pets", "modify pets in your account").
						Scope("read:pets", "read your pets")
				})
			})
	})

	// global security requirements
	doc.GlobalSecurity(map[string][]string{"api_key": {}})

	doc.ExternalDocumentation("http://swagger.io", "Find out more about Swagger")

	setupRoutes(router) // This will define paths and operations

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Swagger UI available at http://localhost:8080/")
	fmt.Println("Swagger JSON available at http://localhost:8080/openapi.json")
	_ = router.Run(":8080")
}
