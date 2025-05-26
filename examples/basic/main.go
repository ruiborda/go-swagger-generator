package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/src/middleware"
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/src/swagger"
	"net/http"
)

type UserDto struct {
	ID   int    `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

func main() {
	router := gin.Default()

	// This will add UserDto to components.schemas
	_, _ = swagger.Swagger().SchemaFromDTO(&UserDto{})

	ConfigureSwagger(router)

	router.GET("/v1/users/:id", GetUserById)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Swagger UI available at http://localhost:8080/")
	fmt.Println("Swagger JSON available at http://localhost:8080/openapi.json")
	_ = router.Run(":8080")
}

// Note: Path in doc.Path() should NOT include the base path from servers configuration.
// The server URL will be prepended by UI tools.
var _ = swagger.Swagger().Path("/users/{id}"). // Path relative to server URL
						Get(func(op openapi.Operation) {
		op.Summary("Find user by ID").
			Tag("UserController"). // Tags are for grouping, can be any string
			OperationID("getUserById").
			PathParameter("id", func(p openapi.Parameter) {
				p.Description("ID of the user to retrieve").
					Required(true).
					Schema(func(s openapi.Schema) {
						s.Type("integer").Format("int64")
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&UserDto{}) // Reference the DTO
					})
			})
	}).
	Doc()

func GetUserById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, UserDto{
		ID:   1, // Example ID, could parse from 'id' param
		Name: "John Doe for ID " + id,
	})
}

func ConfigureSwagger(router *gin.Engine) {
	router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json",
		UIPath:   "/",
		Title:    "Simple API OAS3",
	}))

	doc := swagger.Swagger()

	doc.Info(func(info openapi.Info) {
		info.Title("Simple API with OAS3").
			Version("1.0").
			Description("This is a simple API example using Swagger and OpenAPI 3.0.")
	}).
		Server("http://localhost:8080/v1", func(server openapi.Server) { // Server URL now includes base path
			server.Description("Local development server V1")
		})
	// No BasePath() or Schemes() in OAS3 at the root level, handled by Servers array.
}
