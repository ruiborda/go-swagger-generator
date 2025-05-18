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
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	router := gin.Default()

	SwaggerConfig(router)

	router.GET("/v1/users/:id", GetUserById)

	fmt.Println("Server running on http://localhost:8080")
	_ = router.Run(":8080")
}

var _ = swagger.Swagger().Path("/users/{id}").
	Get(func(op openapi.Operation) {
		op.Summary("Find user by ID").
			Tag("UserController").
			Consumes(mime.ApplicationJSON).
			Produce(mime.ApplicationJSON).
			PathParameter("id", func(p openapi.Parameter) {
				p.
					Required(true).
					Type("integer").
					CollectionFormat("int64")
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("successful operation").
					SchemaFromDTO(UserDto{})
			})
	}).
	Doc()

func GetUserById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"id":   id,
		"name": "John Doe",
	})
}

func SwaggerConfig(router *gin.Engine) {
	router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json",
		UIPath:   "/",
	}))

	doc := swagger.Swagger()

	doc.Info(func(info openapi.Info) {
		info.Title("Simple Api").
			Version("1.0").
			Description("This is a simple API example using SwaggerGin.")
	}).
		Server("/", func(server openapi.Server) {
			server.Description("Servidor de desarrollo local")
		}).
		BasePath("/v1").
		Schemes("http", "https")
}
