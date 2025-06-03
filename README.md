# Go Swagger Generator v2

[![Go Reference](https://pkg.go.dev/badge/github.com/ruiborda/go-swagger-generator.svg)](https://pkg.go.dev/github.com/ruiborda/go-swagger-generator)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go Swagger Generator v2 is a library that facilitates the generation of **OpenAPI 3.0** documentation directly from your Go code with a fluent and elegant API.

## Features

- **Fluent and elegant API** - Chained syntax that makes OpenAPI 3.0 documentation easy to read and write.
- **Simple integration with Gin** - Works with the popular Gin web framework without complications.
- **No annotations required** - Does not require special comments in your code.
- **Built-in Swagger UI** - Includes Swagger UI to explore your API interactively, displaying your OpenAPI 3.0 specification.
- **Compatible with OpenAPI 3.0** - Generates documentation following the OpenAPI 3.0 specification.

## Installation

```bash
# Install Go Swagger Generator v2
 go get -u github.com/ruiborda/go-swagger-generator/v2
```

If you use Gin:

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
	"github.com/ruiborda/go-swagger-generator/v2/src/middleware"
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi"
	"github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/v2/src/swagger"
)

// UserDto represents the user DTO.
type UserDto struct {
	ID   int    `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

func main() {
	router := gin.Default()

	// Configure OpenAPI 3.0 documentation
	ConfigureOpenAPI(router)

	// Define API route. The /v1 prefix should match the server URL in the OpenAPI config.
	router.GET("/v1/users/:id", GetUserByIdHandler)

	fmt.Println("Server running at http://localhost:8080")
	fmt.Println("Swagger UI available at http://localhost:8080/")
	fmt.Println("OpenAPI 3.0 JSON available at http://localhost:8080/openapi.json")
	_ = router.Run(":8080")
}

// Handler to get user by ID.
func GetUserByIdHandler(c *gin.Context) {
	idStr := c.Param("id")
	c.JSON(http.StatusOK, UserDto{
		ID:   1,
		Name: "John Doe (User " + idStr + ")",
	})
}

// Configure OpenAPI 3.0 documentation.
func ConfigureOpenAPI(router *gin.Engine) {
	router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json",
		UIPath:   "/",
		Title:    "My API with OpenAPI 3.0 (v2)",
	}))

	doc := swagger.Swagger()

	doc.Info(func(info openapi.Info) {
		info.Title("Simple User API v2").
			Version("1.0.0").
			Description("A simple API to manage users, documented with Go Swagger Generator v2 and OpenAPI 3.0.")
	})

	doc.Server("http://localhost:8080/v1", func(server openapi.Server) {
		server.Description("Local development server (v1)")
	})

	_, _ = doc.ComponentSchemaFromDTO(&UserDto{})

	var _ = doc.Path("/users/{id}").
		Get(func(op openapi.Operation) {
			op.Summary("Find user by ID").
				Tag("User Management").
				OperationID("getUserByIdV2").
				PathParameter("id", func(p openapi.Parameter) {
					p.Description("ID of the user to retrieve").
						Required(true).
						Schema(func(s openapi.Schema) {
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
```

## Usage Guide (v2 with OpenAPI 3.0)

### 1. Configure OpenAPI in your Gin application

```go
func ConfigureOpenAPI(router *gin.Engine) {
	router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json",
		UIPath:   "/",
		Title:    "My API v2",
	}))
	doc := swagger.Swagger()
	doc.Info(func(info openapi.Info) {
		info.Title("My Awesome API v2").
			Version("2.0.0").
			Description("This is an API example using Go Swagger Generator v2 with OpenAPI 3.0.")
	})
	doc.Server("http://localhost:8080/api/v2", func(server openapi.Server) {
		server.Description("Development server")
	})
}
```

### 2. Define your models (DTOs)

Add `yaml` tags if you plan to support YAML representations, consistent with OpenAPI examples.

```go
type UserDto struct {
	ID   int    `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

// Register the DTO with OpenAPI components:
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&UserDto{})
```

### 3. Document your endpoints (OpenAPI 3.0 style)

```go
var _ = swagger.Swagger().Path("/items/{itemId}")
	.Get(func(op openapi.Operation) {
		op.Summary("Get an item by its ID").
			Tag("Items").
			OperationID("getItemByIdV2").
			PathParameter("itemId", func(p openapi.Parameter) {
				p.Description("ID of the item to fetch").
					Required(true).
					Schema(func(s openapi.Schema) {
						s.Type("integer").Format("int64")
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("Successful operation - item found").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.Schema(openapi.S().Type("object").Property("message", openapi.S().Type("string")))
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
	c.JSON(http.StatusOK, gin.H{
		"message": "Item " + itemId + " found",
	})
}
```

## Documentation (v2)

Check the `/doc_page` directory for detailed documentation on all features of Go Swagger Generator v2:

- [Introduction (v2)](/doc_page/docs/intro.md)
- [Quick Start (v2)](/doc_page/docs/quick-start.md)
- [Defining Models (v2)](/doc_page/docs/defining-models.md)
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

- [Basic Example (v2)](/examples/basic/main.go)
- [Pet Store (v2)](/examples/pet_store/main.go)
- [Array Response Example (v2)](/examples/array_response/main.go)
- [JWT Bearer Auth Example (v2)](/examples/jwt_bearer/main.go)

## Testing your documented API

1. Start your application.
2. Open your browser at the configured UI path (default: `http://localhost:YOUR_PORT/`, for example, http://localhost:8080/).
3. You should see the Swagger UI interface displaying your API documented with OpenAPI 3.0.
4. The OpenAPI 3.0 JSON specification will be available at the configured `JSONPath` (for example, `/openapi.json`).
5. Explore the endpoints, review the required parameters, and test the calls directly from the interface.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/my-feature`).
3. Commit your changes (`git commit -m 'Add a new feature'`).
4. Push to the branch (`git push origin feature/my-feature`).
5. Open a Pull Request.
