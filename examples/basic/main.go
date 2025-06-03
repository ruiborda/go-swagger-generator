package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/v2/src/middleware"
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi"
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/v2/src/swagger"
)

type UserDto struct {
	ID   int    `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

func main() {
	router := gin.Default()

	ConfigureOpenAPI(router)

	router.GET("/v1/users/:id", GetUserByIdHandler)

	fmt.Println("Swagger UI available at http://localhost:8080/")
	fmt.Println("OpenAPI 3.0 JSON available at http://localhost:8080/openapi.json")
	_ = router.Run(":8080")
}

func GetUserByIdHandler(c *gin.Context) {
	idStr := c.Param("id")
	c.JSON(http.StatusOK, UserDto{
		ID:   1,
		Name: "John Doe (User " + idStr + ")",
	})
}

// ConfigureOpenAPI sets up the OpenAPI 3.0 documentation.
func ConfigureOpenAPI(router *gin.Engine) {
	// Register the SwaggerGin middleware.
	router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
		Enabled:  true,                           // Enable Swagger UI and JSON endpoint.
		JSONPath: "/openapi.json",                // Path to serve the OpenAPI JSON.
		UIPath:   "/",                            // Path to serve the Swagger UI.
		Title:    "My API with OpenAPI 3.0 (v2)", // Title for the Swagger UI page.
	}))

	// Get the global OpenAPI document builder instance.
	doc := swagger.Swagger() // Defaults to OpenAPI 3.0 builder.

	doc.Info(func(info openapi.Info) {
		info.Title("Simple User API").
			Version("1.0.0"). // API version
			Description("A simple API to manage users, documented with Go Swagger Generator v2.")
	})

	doc.Server("http://localhost:8080/v1", func(server openapi.Server) {
		server.Description("Local development server (v1)")
	})

	_, _ = doc.SchemaFromDTO(&UserDto{})

	var _ = doc.Path("/users/{id}").
		Get(func(op openapi.Operation) {
			op.Summary("Find user by ID").
				Tag("User Management").     // Groups operations in Swagger UI.
				OperationID("getUserById"). // Unique ID for the operation.
				PathParameter("id", func(p openapi.Parameter) {
					p.Description("ID of the user to retrieve").
						Required(true).
						Schema(func(s openapi.Schema) { // Define schema for the path parameter.
							s.Type("integer").Format("int64")
						})
				}).
				Response(http.StatusOK, func(r openapi.Response) {
					r.Description("Successful operation - user details returned").
						Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
							mt.SchemaFromDTO(&UserDto{})
						})
				}).
				Response(http.StatusNotFound, func(r openapi.Response) {
					r.Description("User not found")
				})
		}).
		Doc()
}
