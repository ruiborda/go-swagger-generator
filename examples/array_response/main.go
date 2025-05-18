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

// GetUserByIdResponse defines a user DTO for the response
type GetUserByIdResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Age      int    `json:"age"`
	Active   bool   `json:"active"`
}

// Define two endpoints: one returning a single user and another returning an array of users
func main() {
	router := gin.Default()

	SwaggerConfig(router)

	// Individual user endpoint
	router.GET("/v1/users/:id", GetUserById)

	// Array of users endpoint - demonstrating the array response feature
	router.GET("/v1/users", GetAllUsers)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Try accessing the array endpoint at: http://localhost:8080/v1/users")
	fmt.Println("Swagger UI available at: http://localhost:8080/")
	_ = router.Run(":8080")
}

// Swagger documentation for the single user endpoint
var _ = swagger.Swagger().Path("/users/{id}").
	Get(func(op openapi.Operation) {
		op.Summary("Find user by ID").
			Description("Returns a single user by ID").
			Tag("Users").
			Consumes(mime.ApplicationJSON).
			Produce(mime.ApplicationJSON).
			PathParameter("id", func(p openapi.Parameter) {
				p.
					Description("ID of the user to retrieve").
					Required(true).
					Type("integer").
					Format("int64")
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("Successful operation").
					SchemaFromDTO(&GetUserByIdResponse{})
			}).
			Response(http.StatusNotFound, func(r openapi.Response) {
				r.Description("User not found")
			})
	}).
	Doc()

// Swagger documentation for the array of users endpoint
var _ = swagger.Swagger().Path("/users").
	Get(func(op openapi.Operation) {
		op.Summary("Get all users").
			Description("Returns an array of all registered users").
			Tag("Users").
			Consumes(mime.ApplicationJSON).
			Produce(mime.ApplicationJSON).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("Successful operation").
					// This is the key part demonstrating our fixed functionality
					// Using SchemaFromDTO with an array type
					SchemaFromDTO(&[]*GetUserByIdResponse{})
			})
	}).
	Doc()

// Handler for getting a user by ID
func GetUserById(c *gin.Context) {
	id := c.Param("id")

	// Return a sample user
	user := GetUserByIdResponse{
		ID:       1,
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
	// Return a sample array of users
	users := []*GetUserByIdResponse{
		{
			ID:       1,
			Username: "user1",
			Email:    "user1@example.com",
			FullName: "User One",
			Age:      30,
			Active:   true,
		},
		{
			ID:       2,
			Username: "user2",
			Email:    "user2@example.com",
			FullName: "User Two",
			Age:      25,
			Active:   true,
		},
		{
			ID:       3,
			Username: "user3",
			Email:    "user3@example.com",
			FullName: "User Three",
			Age:      40,
			Active:   false,
		},
	}

	c.JSON(http.StatusOK, users)
}

// Swagger configuration
func SwaggerConfig(router *gin.Engine) {
	router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json",
		UIPath:   "/",
	}))

	doc := swagger.Swagger()

	doc.Info(func(info openapi.Info) {
		info.Title("Array Response Example API").
			Version("1.0").
			Description("This example demonstrates how to use array responses with SchemaFromDTO in go-swagger-generator")
	}).
		Server("/", func(server openapi.Server) {
			server.Description("Development server")
		}).
		BasePath("/v1").
		Schemes("http", "https")
}
