---
sidebar_position: 1
title: Quick Start (v2)
---

# Quick Start Guide (v2)

This guide will help you start using Go Swagger Generator v2 in just a few minutes to create OpenAPI 3.0 documentation.

## 1. Installation

Install Go Swagger Generator v2 and the Gin framework:

```bash
# Install Go Swagger Generator v2
go get -u github.com/ruiborda/go-swagger-generator@v2

# Install Gin Framework
go get github.com/gin-gonic/gin
```

## 2. Copy and paste this example

Create a `main.go` file with the following content. This example sets up a simple API with one endpoint and its OpenAPI 3.0 documentation.

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/v2/src/middleware" // Assuming 'src' is part of the import path
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
    "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
    "github.com/ruiborda/go-swagger-generator/v2/src/swagger"
)

// UserDto represents the data transfer object for a user.
type UserDto struct {
    ID   int    `json:"id" yaml:"id"`
    Name string `json:"name" yaml:"name"`
}

func main() {
    router := gin.Default()

    // Configure OpenAPI 3.0 documentation
    ConfigureOpenAPI(router)

    // Define the API route. The /v1 prefix matches the server URL in OpenAPI config.
    router.GET("/v1/users/:id", GetUserByIdHandler)

    fmt.Println("Server running on http://localhost:8080")
    fmt.Println("Swagger UI available at http://localhost:8080/")
    fmt.Println("OpenAPI 3.0 JSON available at http://localhost:8080/openapi.json")
    _ = router.Run(":8080")
}

// GetUserByIdHandler is the Gin handler for retrieving a user.
func GetUserByIdHandler(c *gin.Context) {
    idStr := c.Param("id")
    // In a real application, convert idStr to int and fetch user data.
    // For simplicity, we return a dummy UserDto.
    c.JSON(http.StatusOK, UserDto{
        ID:   1, // Example ID
        Name: "John Doe (User " + idStr + ")",
    })
}

// ConfigureOpenAPI sets up the OpenAPI 3.0 documentation.
func ConfigureOpenAPI(router *gin.Engine) {
    // Register the SwaggerGin middleware.
    router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
        Enabled:  true,            // Enable Swagger UI and JSON endpoint.
        JSONPath: "/openapi.json", // Path to serve the OpenAPI JSON.
        UIPath:   "/",             // Path to serve the Swagger UI.
        Title:    "My API with OpenAPI 3.0 (v2)", // Title for the Swagger UI page.
    }))

    // Get the global OpenAPI document builder instance.
    doc := swagger.Swagger() // Defaults to OpenAPI 3.0 builder.

    // Define general API information (Info Object).
    doc.Info(func(info openapi.Info) {
        info.Title("Simple User API").
            Version("1.0.0"). // API version
            Description("A simple API to manage users, documented with Go Swagger Generator v2.")
    })

    // Define the server(s) for the API.
    // Paths in Path() calls will be relative to these server URLs.
    doc.Server("http://localhost:8080/v1", func(server openapi.Server) {
        server.Description("Local development server (v1)")
    })
    // You can add more servers (e.g., production, staging).

    // It's good practice to register DTOs used in responses/request bodies.
    // This makes them available in the #/components/schemas/ section.
    _, _ = doc.SchemaFromDTO(&UserDto{})

    // Document the /users/{id} GET endpoint.
    // The path "/users/{id}" is relative to the server URL (e.g., http://localhost:8080/v1/users/{id}).
    var _ = doc.Path("/users/{id}").
        Get(func(op openapi.Operation) {
            op.Summary("Find user by ID").
                Tag("User Management"). // Groups operations in Swagger UI.
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
                        Content(mime.ApplicationJSON, func(mt openapi.MediaType) { // Define response content type and schema.
                            mt.SchemaFromDTO(&UserDto{}) // Use UserDto for the response schema.
                        })
                }).
                Response(http.StatusNotFound, func(r openapi.Response) {
                    r.Description("User not found")
                })
        }).
        Doc() // Finalizes and registers this path documentation.
}

```

## 3. Run your application

Initialize your Go module if you haven't already:
```bash
go mod init yourprojectname
```

Then, ensure all dependencies are downloaded:
```bash
go mod tidy
```

Finally, run your application:
```bash
go run main.go
```

## 4. View OpenAPI documentation

Open your browser at [http://localhost:8080](http://localhost:8080) to access the Swagger UI interface.
The OpenAPI 3.0 JSON specification will be available at [http://localhost:8080/openapi.json](http://localhost:8080/openapi.json).

## Next steps

For more information, see:
- [Introduction (v2)](/docs/intro)
- [Defining Models (Schemas)](/docs/defining-models)
- [Security Schemes](/docs/security)
- [Production Configuration](/docs/production)