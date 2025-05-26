package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/src/middleware"
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/src/swagger"
)

type UserResponse struct {
	ID       int    `json:"id" yaml:"id"`
	Username string `json:"username" yaml:"username"`
	Email    string `json:"email" yaml:"email"`
	FullName string `json:"fullName" yaml:"fullName"`
	Age      int    `json:"age" yaml:"age"`
	Active   bool   `json:"active" yaml:"active"`
}

func main() {
	fmt.Println("Server running on http://localhost:8080")
	router := gin.Default()

	_, _ = swagger.Swagger().SchemaFromDTO(&UserResponse{})

	ConfigureSwagger(router)

	// Individual user endpoint
	router.GET("/v1/users/:id", GetUserById)

	// Array of users endpoint
	router.GET("/v1/users", GetAllUsers)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Try accessing the array endpoint at: http://localhost:8080/v1/users")
	fmt.Println("Swagger UI available at: http://localhost:8080/")
	router.Run(":8080")
}

// Swagger documentation for the single user endpoint
var _ = swagger.Swagger().Path("/users/{id}"). // Path relative to server base path
						Get(func(op openapi.Operation) {
		op.Summary("Find user by ID").
			Description("Returns a single user by ID").
			Tag("Users").
			OperationID("getUserById").
			PathParameter("id", func(p openapi.Parameter) {
				p.Description("ID of the user to retrieve").
					Required(true).
					Schema(func(s openapi.Schema) {
						s.Type("integer").Format("int64")
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("Successful operation").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&UserResponse{})
					})
			}).
			Response(http.StatusNotFound, func(r openapi.Response) {
				r.Description("User not found")
			})
	}).
	Doc()

var _ = swagger.Swagger().Path("/users").
	Get(func(op openapi.Operation) {
		op.Summary("Get all users").
			Description("Returns an array of all registered users").
			Tag("Users").
			OperationID("getAllUsers").
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("Successful operation").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&[]*UserResponse{})
					})
			})
	}).
	Doc()

// Handler for getting a user by ID
func GetUserById(c *gin.Context) {
	id := c.Param("id")
	user := UserResponse{
		ID:       1, // Example, parse 'id'
		Username: "user" + id,
		Email:    "user" + id + "@example.com",
		FullName: "User " + id,
		Age:      30,
		Active:   true,
	}
	c.JSON(http.StatusOK, user)
}

// Handler for getting all users - returns an array of users
func GetAllUsers(c *gin.Context) {
	users := []*UserResponse{
		{ID: 1, Username: "user1", Email: "user1@example.com", FullName: "User One", Age: 30, Active: true},
		{ID: 2, Username: "user2", Email: "user2@example.com", FullName: "User Two", Age: 25, Active: true},
		{ID: 3, Username: "user3", Email: "user3@example.com", FullName: "User Three", Age: 40, Active: false},
	}
	c.JSON(http.StatusOK, users)
}

// Swagger configuration
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
		Server("http://localhost:8080/v1", func(server openapi.Server) {
			server.Description("Local development server V1")
		})
}
