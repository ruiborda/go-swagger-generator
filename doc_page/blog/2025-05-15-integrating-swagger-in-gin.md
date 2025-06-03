---
slug: integrating-openapi-v2-in-gin
title: How to Integrate OpenAPI 3.0 in your Gin API with Go Swagger Generator v2
authors: [rui]
tags: [go, openapi, oas3, gin, api, documentation, v2]
---

# How to Integrate OpenAPI 3.0 in your Gin API with Go Swagger Generator v2

In modern API development, clear and accessible documentation is as important as the code itself. A well-documented API facilitates its adoption, reduces friction during integration, and saves time for both internal and external developers.

In this tutorial, we will learn how to integrate OpenAPI 3.0 into a Gin application using **Go Swagger Generator v2**, a library that makes it easy to generate OpenAPI documentation directly from your Go code.

<!-- truncate -->

## Why use Go Swagger Generator v2?

- **Fluid and elegant API** - Chained syntax that makes documentation easy to read and write for OpenAPI 3.0.
- **Simple integration with Gin** - Works with the popular Gin web framework without complications.
- **No annotations needed** - No special comments required in your code.
- **Built-in Swagger UI** - Includes Swagger UI to interactively explore your API, rendering your OpenAPI 3.0 spec.

## Step 1: Installing dependencies

The first thing we need to do is install both Gin and the Go Swagger Generator v2:

```bash
# Install Gin Framework
go get github.com/gin-gonic/gin

# Install Go Swagger Generator v2
go get -u github.com/ruiborda/go-swagger-generator@v2
```

## Step 2: Defining models (DTOs)

Let's start by defining a simple structure that will be part of our API. In this case, we'll define a `UserDto`:

```go
type UserDto struct {
	ID   int    `json:"id" yaml:"id"` // Added yaml tag for consistency with examples
	Name string `json:"name" yaml:"name"`
}
```

This structure represents the user data that the API will return.

## Step 3: Configuring OpenAPI 3.0 Documentation

Now we need to configure OpenAPI 3.0 documentation in our Gin application. We'll create a dedicated function for this:

```go
func ConfigureOpenAPI(router *gin.Engine) {
	// Enable Swagger middleware for OpenAPI 3.0
	router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json", // Path for OpenAPI JSON
		UIPath:   "/",             // Path for Swagger UI
		Title:    "Simple API with OpenAPI 3.0", // Title for the Swagger UI page
	}))

	// Get the OpenAPI document instance (default is OAS3)
	doc := swagger.Swagger()

	// Configure basic API information (Info Object)
	doc.Info(func(info openapi.Info) {
		info.Title("Simple API").
			Version("1.0.0"). // Semantic versioning for your API
			Description("This is a simple API example using Go Swagger Generator v2 and OpenAPI 3.0.")
	})

	// Configure Servers (OAS3 replaces BasePath and Schemes)
	doc.Server("http://localhost:8080/v1", func(server openapi.Server) {
		server.Description("Local development server - API version 1")
	})
	// You can add more servers (e.g., staging, production)
	// doc.Server("https://api.example.com/v1", func(server openapi.Server) {
	// 	server.Description("Production server - API version 1")
	// })

    // Register DTOs to be available in #/components/schemas/
    // This is good practice, though SchemaFromDTO in operations also registers them.
    _, _ = doc.SchemaFromDTO(&UserDto{})
}
```

This function does several important things:

1. Registers the Swagger middleware in the Gin router.
2. Configures the routes where OpenAPI JSON and Swagger UI will be exposed.
3. Defines basic API metadata (Info Object) such as title, version, and description.
4. Sets up server information using the OpenAPI 3.0 `servers` array.
5. Registers DTOs for use in the documentation.

## Step 4: Defining an endpoint handler

Now, we'll define a simple Gin handler to retrieve a user by their ID:

```go
func GetUserById(c *gin.Context) {
	idStr := c.Param("id")
    // In a real app, parse idStr to int and fetch user data
	c.JSON(http.StatusOK, UserDto{
		ID:   1, // Example, use parsed id
		Name: "John Doe (id: " + idStr + ")",
	})
}
```

This function gets the user ID from the route parameters and returns user data in JSON format.

## Step 5: Documenting the endpoint with OpenAPI 3.0 syntax

This is where Go Swagger Generator v2 shines. To document our endpoint, we use a fluid syntax adhering to OpenAPI 3.0 principles:

```go
// Document the /users/{id} GET endpoint
// Path is relative to the server URLs defined in doc.Server(...)
var _ = swagger.Swagger().Path("/users/{id}"). 
	Get(func(op openapi.Operation) {
		op.Summary("Find user by ID").
			Tag("User Operations"). // For grouping in Swagger UI
			OperationID("getUserById"). // Unique ID for the operation
			PathParameter("id", func(p openapi.Parameter) {
				p.Description("ID of the user to retrieve").
					Required(true).
					Schema(func(s openapi.Schema) { // Define schema for the path parameter
						s.Type("integer").Format("int64")
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("Successful operation - user found").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) { // Define content type and schema
						mt.SchemaFromDTO(&UserDto{}) // Use the UserDto for the response schema
					})
			}).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("User not found")
            })
	}).
	Doc() // Registers this documentation path
```

This documentation:

1. Defines an endpoint at the route `/users/{id}` (relative to server URLs) with the GET method.
2. Provides a clear summary, tag for grouping, and a unique operation ID.
3. Specifies the content type (JSON) and schema for the successful response using `Content` and `SchemaFromDTO`.
4. Documents the path parameter `id` as a required 64-bit integer using `PathParameter` and its `Schema`.
5. Defines possible responses, including `200 OK` and `404 Not Found`.

The `.Doc()` method at the end registers this documentation in the OpenAPI instance.

## Step 6: Implementing the main function

Finally, we put everything together in our `main` function:

```go
func main() {
	router := gin.Default()

	// Configure OpenAPI documentation
	ConfigureOpenAPI(router)

	// Register our GET route for retrieving users (note the /v1 prefix matching server config)
	router.GET("/v1/users/:id", GetUserById)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Swagger UI available at http://localhost:8080/")
	fmt.Println("OpenAPI 3.0 JSON available at http://localhost:8080/openapi.json")
	_ = router.Run(":8080")
}
```

Here:
1. We create a Gin router.
2. We configure OpenAPI documentation in the router.
3. We register our GET route `/v1/users/:id` for retrieving users. The `/v1` prefix matches one of our defined server URLs.
4. We start the server on port 8080.

## Complete code

Here's the complete code for our application:

```go
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

type UserDto struct {
	ID   int    `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

func main() {
	router := gin.Default()

	ConfigureOpenAPI(router)

	router.GET("/v1/users/:id", GetUserById)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Swagger UI available at http://localhost:8080/")
	fmt.Println("OpenAPI 3.0 JSON available at http://localhost:8080/openapi.json")
	_ = router.Run(":8080")
}

// Document the /users/{id} GET endpoint
var _ = swagger.Swagger().Path("/users/{id}").
	Get(func(op openapi.Operation) {
		op.Summary("Find user by ID").
			Tag("User Operations").
			OperationID("getUserById").
			PathParameter("id", func(p openapi.Parameter) {
				p.Description("ID of the user to retrieve").
					Required(true).
					Schema(func(s openapi.Schema) {
						s.Type("integer").Format("int64")
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("Successful operation - user found").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&UserDto{})
					})
			}).
            Response(http.StatusNotFound, func(r openapi.Response) {
                r.Description("User not found")
            })
	}).
	Doc()

func GetUserById(c *gin.Context) {
	idStr := c.Param("id")
	c.JSON(http.StatusOK, UserDto{
		ID:   1, // Example ID
		Name: "John Doe (id: " + idStr + ")",
	})
}

func ConfigureOpenAPI(router *gin.Engine) {
	router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json",
		UIPath:   "/",
		Title:    "Simple API with OpenAPI 3.0",
	}))

	doc := swagger.Swagger()

	doc.Info(func(info openapi.Info) {
		info.Title("Simple API").
			Version("1.0.0").
			Description("This is a simple API example using Go Swagger Generator v2 and OpenAPI 3.0.")
	})
	doc.Server("http://localhost:8080/v1", func(server openapi.Server) {
		server.Description("Local development server - API version 1")
	})

    _, _ = doc.SchemaFromDTO(&UserDto{})
}

```

## Testing our documented API

To test our implementation:

1. Save the above code in a `main.go` file.
2. Run `go mod init yourprojectname` (if you haven't already).
3. Run `go mod tidy` to ensure you have all dependencies (it will fetch `go-swagger-generator@v2`).
4. Start the application with `go run main.go`.
5. Open your browser at [http://localhost:8080](http://localhost:8080).

You should see the Swagger UI interface displaying your documented API. You can explore the endpoints, see the required parameters, and test the calls directly from the interface.

## Conclusion

Integrating OpenAPI 3.0 into a Gin API using Go Swagger Generator v2 is a straightforward process that offers great benefits. In just a few minutes, you get interactive and professional documentation that evolves along with your code.

Go Swagger Generator v2 provides an elegant and fluid syntax for documenting your APIs, making the process more enjoyable and less error-prone than comment-based solutions, fully leveraging the power of OpenAPI 3.0.

Do you have any questions about integrating OpenAPI 3.0 into your Go API with v2? Let us know in the comments!