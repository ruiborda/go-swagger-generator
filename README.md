# Go Swagger Generator v2

[![Go Reference](https://pkg.go.dev/badge/github.com/ruiborda/go-swagger-generator.svg)](https://pkg.go.dev/github.com/ruiborda/go-swagger-generator)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go Swagger Generator v2 is a library that makes it easy to generate **OpenAPI 3.0** documentation directly from your Go code with a fluid and elegant API.

## Features

- **Fluid and elegant API** - Chained syntax that makes OpenAPI 3.0 documentation easy to read and write.
- **Simple integration with Gin** - Works with the popular Gin web framework without complications.
- **No annotations needed** - No special comments required in your code.
- **Built-in Swagger UI** - Includes Swagger UI to interactively explore your API, rendering your OpenAPI 3.0 specification.
- **OpenAPI 3.0 Compliant** - Generates documentation following the OpenAPI 3.0 specification.

## Installation

```bash
# Install Go Swagger Generator v2
go get -u github.com/ruiborda/go-swagger-generator@v2
```

If you're using Gin:

```bash
# Install Gin Framework
go get github.com/gin-gonic/gin
```

## Quick Start (v2)

Here's a simple example showing how to integrate Go Swagger Generator v2 with Gin to produce OpenAPI 3.0 documentation:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/src/middleware" // Assuming 'src' is part of the import path for v2
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/src/swagger"
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
	doc := swagger.Swagger() // Defaults to OpenAPI 3.0 builder in v2.

	// Define general API information (Info Object).
	doc.Info(func(info openapi.Info) {
		info.Title("Simple User API v2").
			Version("1.0.0"). // API version
			Description("A simple API to manage users, documented with Go Swagger Generator v2 and OpenAPI 3.0.")
	})

	// Define the server(s) for the API.
	// Paths in Path() calls will be relative to these server URLs.
	doc.Server("http://localhost:8080/v1", func(server openapi.Server) {
		server.Description("Local development server (v1)")
	})

	// Register DTOs to make them available in #/components/schemas/
	_, _ = doc.ComponentSchemaFromDTO(&UserDto{})

	// Document the /users/{id} GET endpoint.
	// The path "/users/{id}" is relative to the server URL (e.g., http://localhost:8080/v1/users/{id}).
	var _ = doc.Path("/users/{id}").
		Get(func(op openapi.Operation) {
			op.Summary("Find user by ID").
				Tag("User Management").      // Groups operations in Swagger UI.
				OperationID("getUserByIdV2"). // Unique ID for the operation.
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

## Usage Guide (v2 with OpenAPI 3.0)

### 1. Configure OpenAPI in your Gin application

```go
func ConfigureOpenAPI(router *gin.Engine) {
	// Enable Swagger middleware
	router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json", // Path for OpenAPI 3.0 JSON
		UIPath:   "/",             // Path for Swagger UI
		Title:    "My API v2",
	}))

	// Get the OpenAPI document instance (defaults to OAS3)
	doc := swagger.Swagger()

	// Configure basic API information (Info Object)
	doc.Info(func(info openapi.Info) {
		info.Title("My Awesome API v2").
			Version("2.0.0").
			Description("This is an API example using Go Swagger Generator v2 with OpenAPI 3.0.")
	})

	// Configure the server(s) (OAS3 uses a 'servers' array)
	doc.Server("http://localhost:8080/api/v2", func(server openapi.Server) {
		server.Description("Development server")
	})
	// Add more servers if needed (e.g., staging, production)
	// doc.Server("https://api.example.com/api/v2", func(server openapi.Server) { ... })
}
```

### 2. Define your models (DTOs)

Add `yaml` tags if you plan to support YAML representations, consistent with OpenAPI examples.

```go
type UserDto struct {
	ID   int    `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

// Register DTO with OpenAPI components:
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&UserDto{})
```

### 3. Document your endpoints (OpenAPI 3.0 style)

```go
var _ = swagger.Swagger().Path("/items/{itemId}"). // Path relative to server URL(s)
	Get(func(op openapi.Operation) {
		op.Summary("Get an item by its ID").
			Tag("Items").
			OperationID("getItemByIdV2").
			PathParameter("itemId", func(p openapi.Parameter) {
				p.Description("ID of the item to fetch").
					Required(true). // Path parameters are always required
					Schema(func(s openapi.Schema) {
						s.Type("integer").Format("int64")
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("Successful operation - item found").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						// Assuming ItemDto is defined and registered
						// mt.SchemaFromDTO(&ItemDto{})
						mt.Schema(openapi.S().Type("object").Property("message", openapi.S().Type("string"))) // Example inline schema
					})
			}).
			Response(http.StatusNotFound, func(r openapi.Response) {
				r.Description("Item not found")
			})
	}).
	Doc()
```

### 4. Implement your handler functions

```go
func GetItemById(c *gin.Context) {
	itemId := c.Param("itemId")
	// ... your logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Item " + itemId + " found",
	})
}
```

## Documentation (v2)

Check the `/doc_page` directory for detailed documentation on all the features of Go Swagger Generator v2:

- [Introduction (v2)](/doc_page/docs/intro.md)
- [Quick Start (v2)](/doc_page/docs/quick-start.md)
- [Defining Models (Schemas) (v2)](/doc_page/docs/defining-models.md)
- [Path Parameters (v2)](/doc_page/docs/path-parameters.md)
- [Query Parameters (v2)](/doc_page/docs/query-parameters.md)
- [Request Bodies (v2)](/doc_page/docs/request-bodies.md)
- [Responses (v2)](/doc_page/docs/responses.md)
- [Security Schemes (v2)](/doc_page/docs/security.md)
- [Array Responses (v2)](/doc_page/docs/array-responses.md)
- [Production Configuration (v2)](/doc_page/docs/production.md)
- And more...

## Examples (v2)

Check the `/examples` directory for complete examples updated for v2:

- [Basic Example (v2)](/examples/basic/main.go) - A simple API with basic OpenAPI 3.0 features.
- [Pet Store (v2)](/examples/pet_store/main.go) - A more complex example based on the Swagger Pet Store, adapted for OpenAPI 3.0.
- [Array Response Example (v2)](/examples/array_response/main.go) - Demonstrates array request/response bodies.
- [JWT Bearer Auth Example (v2)](/examples/jwt_bearer/main.go) - Shows JWT Bearer authentication scheme.


## Testing your documented API

After implementing your API with Go Swagger Generator v2:

1. Start your application.
2. Open your browser at the configured UI path (default: `http://localhost:YOUR_PORT/`, e.g., http://localhost:8080/).
3. You should see the Swagger UI interface displaying your OpenAPI 3.0 documented API.
4. The OpenAPI 3.0 JSON specification will be available at the configured `JSONPath` (e.g., `/openapi.json`).
5. Explore the endpoints, see the required parameters, and test the calls directly from the interface.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/amazing-feature`).
3. Commit your changes (`git commit -m 'Add some amazing feature'`).
4. Push to the branch (`git push origin feature/amazing-feature`).
5. Open a Pull Request.

