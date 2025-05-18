---
slug: integrating-swagger-in-gin
title: How to Integrate Swagger in your Gin API with Go Swagger Generator
authors: [rui]
tags: [go, swagger, gin, api, documentation]
---

# How to Integrate Swagger in your Gin API with Go Swagger Generator

In modern API development, clear and accessible documentation is as important as the code itself. A well-documented API facilitates its adoption, reduces friction during integration, and saves time for both internal and external developers.

In this tutorial, we will learn how to integrate Swagger (OpenAPI) into a Gin application using **Go Swagger Generator**, a library that makes it easy to generate OpenAPI documentation directly from your Go code.

<!-- truncate -->

## Why use Go Swagger Generator?

- **Fluid and elegant API** - Chained syntax that makes documentation easy to read and write
- **Simple integration with Gin** - Works with the popular Gin web framework without complications
- **No annotations needed** - No special comments required in your code
- **Built-in Swagger UI** - Includes Swagger UI to interactively explore your API

## Step 1: Installing dependencies

The first thing we need to do is install both Gin and the Swagger generator:

```bash
# Install Gin Framework
go get github.com/gin-gonic/gin

# Install Go Swagger Generator
go get -u github.com/ruiborda/go-swagger-generator@v1
```

## Step 2: Defining models (DTOs)

Let's start by defining a simple structure that will be part of our API. In this case, we'll define a `UserDto`:

```go
type UserDto struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
```

This structure represents the user data that the API will return.

## Step 3: Configuring Swagger

Now we need to configure Swagger in our Gin application. We'll create a dedicated function for this:

```go
func SwaggerConfig(router *gin.Engine) {
	// Enable Swagger middleware
	router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json", // Path for OpenAPI JSON
		UIPath:   "/",             // Path for Swagger UI
	}))

	// Get the Swagger document instance
	doc := swagger.Swagger()

	// Configure basic API information
	doc.Info(func(info openapi.Info) {
		info.Title("Simple Api").
			Version("1.0").
			Description("This is a simple API example using SwaggerGin.")
	}).
		// Configure the server
		Server("/", func(server openapi.Server) {
			server.Description("Local development server")
		}).
		// Configure base path and schemas
		BasePath("/v1").
		Schemes("http", "https")
}
```

This function does several important things:

1. Registers the Swagger middleware in the Gin router
2. Configures the routes where Swagger will be exposed (JSON and UI)
3. Defines basic API metadata such as title, version, and description
4. Sets up server information and accepted schemes

## Step 4: Defining an endpoint

Now, we'll define a simple endpoint to retrieve a user by their ID:

```go
func GetUserById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"id":   id,
		"name": "John Doe",
	})
}
```

This function gets the user ID from the route parameters and returns user data in JSON format.

## Step 5: Documenting the endpoint with Swagger

This is where Go Swagger Generator shines. To document our endpoint, we use a fluid syntax:

```go
var _ = swagger.Swagger().Path("/users/{id}").
	Get(func(op openapi.Operation) {
		op.Summary("Find user by ID").
			Tag("UserController").
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
```

This documentation:

1. Defines an endpoint at the route `/users/{id}` with GET method
2. Provides a clear summary and tag for grouping endpoints
3. Specifies that the endpoint produces JSON
4. Documents the route parameter `id` as a required integer
5. Defines the successful response that returns a `UserDto` object

The `.Doc()` method at the end registers this documentation in the Swagger instance.

## Step 6: Implementing the main function

Finally, we put everything together in our `main` function:

```go
func main() {
	router := gin.Default()

	SwaggerConfig(router)

	router.GET("/v1/users/:id", GetUserById)

	fmt.Println("Server running on http://localhost:8080")
	_ = router.Run(":8080")
}
```

Here:
1. We create a Gin router
2. We configure Swagger in the router
3. We register our GET route for retrieving users
4. We start the server on port 8080

## Complete code

Here's the complete code for our application:

```go
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
			server.Description("Local development server")
		}).
		BasePath("/v1").
		Schemes("http", "https")
}
```

## Testing our documented API

To test our implementation:

1. Save the above code in a `main.go` file
2. Run `go mod tidy` to ensure you have all dependencies
3. Start the application with `go run main.go`
4. Open your browser at [http://localhost:8080](http://localhost:8080)

You should see the Swagger UI interface displaying your documented API. You can explore the endpoints, see the required parameters, and test the calls directly from the interface.

## Conclusion

Integrating Swagger into a Gin API using Go Swagger Generator is a simple process that offers great benefits. In just a few minutes, you get interactive and professional documentation that evolves along with your code.

Go Swagger Generator provides an elegant and fluid syntax for documenting your APIs, making the process more enjoyable and less error-prone than comment-based solutions.

Do you have any questions about integrating Swagger into your Go API? Let us know in the comments!