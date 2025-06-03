package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi"
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/v2/src/swagger"
	"net/http"
)

type UserDto struct {
	ID   int    `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

func main() {
	router := gin.Default()

	// Registrar el esquema UserDto en Swagger
	doc := swagger.Swagger()
	_, _ = doc.SchemaFromDTO(&UserDto{})

	ConfigureSwagger(doc)

	router.GET("/v1/users/:id", GetUserById)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Swagger UI available at http://localhost:8080/")
	fmt.Println("Swagger JSON available at http://localhost:8080/openapi.json")
	_ = router.Run(":8080")
}

// Documentación de la ruta usando la API v2
func ConfigureSwagger(doc openapi.SwaggerDocBuilder) {
	doc.Info(func(info openapi.Info) {
		info.Title("Simple API with OAS3").
			Version("1.0").
			Description("This is a simple API example using Swagger and OpenAPI 3.0.")
	}).
		Server("http://localhost:8080/v1", func(server openapi.Server) {
			server.Description("Local development server V1")
		})

	doc.Path("/users/{id}", func(path openapi.PathItem) {
		path.Get(func(op openapi.Operation) {
			op.Summary("Find user by ID").
				Tag("UserController").
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
							mt.SchemaFromDTO(&UserDto{})
						})
				})
		})
	})
}

func GetUserById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, UserDto{
		ID:   1, // Example ID, could parse from 'id' param
		Name: "John Doe for ID " + id,
	})
}
