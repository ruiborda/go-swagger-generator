# Go Swagger Generator

[![Go Reference](https://pkg.go.dev/badge/github.com/ruiborda/go-swagger-generator.svg)](https://pkg.go.dev/github.com/ruiborda/go-swagger-generator)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go Swagger Generator is a library that makes it easy to generate OpenAPI documentation directly from your Go code with a fluid and elegant API.

## Features

- **Fluid and elegant API** - Chained syntax that makes documentation easy to read and write
- **Simple integration with Gin** - Works with the popular Gin web framework without complications
- **No annotations needed** - No special comments required in your code
- **Built-in Swagger UI** - Includes Swagger UI to interactively explore your API

## Installation

```bash
# Install Go Swagger Generator
go get github.com/ruiborda/go-swagger-generator
```

If you're using Gin:

```bash
# Install Gin Framework
go get github.com/gin-gonic/gin
```

## Quick Start

Here's a simple example showing how to integrate Go Swagger Generator with Gin:

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

## Usage Guide

### 1. Configure Swagger in your Gin application

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

### 2. Define your models (DTOs)

```go
type UserDto struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
```

### 3. Document your endpoints

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

### 4. Implement your handler functions

```go
func GetUserById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"id":   id,
		"name": "John Doe",
	})
}
```

## Documentation

Check the `/doc_page` directory for detailed documentation on all the features of Go Swagger Generator:

- [Introduction](/doc_page/docs/intro.md)
- [Quick Start](/doc_page/docs/quick-start.md)
- [Defining Models](/doc_page/docs/defining-models.md)
- [Path Parameters](/doc_page/docs/path-parameters.md)
- [Query Parameters](/doc_page/docs/query-parameters.md)
- [Request Bodies](/doc_page/docs/request-bodies.md)
- [Responses](/doc_page/docs/responses.md)
- [Security](/doc_page/docs/security.md)
- And more...

## Examples

Check the `/examples` directory for complete examples:

- [Basic Example](/examples/basic/main.go) - A simple API with basic features
- [Pet Store](/examples/pet_store/main.go) - A more complex example based on the Swagger Pet Store

## Testing your documented API

After implementing your API with Go Swagger Generator:

1. Start your application
2. Open your browser at the configured UI path (default: http://localhost:8080)
3. You should see the Swagger UI interface displaying your documented API
4. Explore the endpoints, see the required parameters, and test the calls directly from the interface

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request