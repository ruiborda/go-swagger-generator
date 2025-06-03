# Go Swagger Generator v2

[![Go Reference](https://pkg.go.dev/badge/github.com/ruiborda/go-swagger-generator.svg)](https://pkg.go.dev/github.com/ruiborda/go-swagger-generator)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go Swagger Generator v2 es una librería que facilita la generación de documentación **OpenAPI 3.0** directamente desde tu código Go con una API fluida y elegante.

## Características

- **API fluida y elegante** - Sintaxis encadenada que hace que la documentación OpenAPI 3.0 sea fácil de leer y escribir.
- **Integración simple con Gin** - Funciona con el popular framework web Gin sin complicaciones.
- **Sin anotaciones necesarias** - No requiere comentarios especiales en tu código.
- **Swagger UI incorporado** - Incluye Swagger UI para explorar tu API de forma interactiva, mostrando tu especificación OpenAPI 3.0.
- **Compatible con OpenAPI 3.0** - Genera documentación siguiendo la especificación OpenAPI 3.0.

## Instalación

```bash
# Instala Go Swagger Generator v2
 go get -u github.com/ruiborda/go-swagger-generator/v2
```

Si usas Gin:

```bash
# Instala Gin Framework
 go get github.com/gin-gonic/gin
```

## Inicio Rápido (v2)

Aquí tienes un ejemplo simple mostrando cómo integrar Go Swagger Generator v2 con Gin para producir documentación OpenAPI 3.0:

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

// UserDto representa el DTO de usuario.
type UserDto struct {
	ID   int    `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

func main() {
	router := gin.Default()

	// Configura la documentación OpenAPI 3.0
	ConfigureOpenAPI(router)

	// Define la ruta de la API. El prefijo /v1 debe coincidir con la URL del servidor en la config OpenAPI.
	router.GET("/v1/users/:id", GetUserByIdHandler)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	fmt.Println("Swagger UI disponible en http://localhost:8080/")
	fmt.Println("OpenAPI 3.0 JSON disponible en http://localhost:8080/openapi.json")
	_ = router.Run(":8080")
}

// Handler para obtener usuario por ID.
func GetUserByIdHandler(c *gin.Context) {
	idStr := c.Param("id")
	c.JSON(http.StatusOK, UserDto{
		ID:   1,
		Name: "John Doe (User " + idStr + ")",
	})
}

// Configura la documentación OpenAPI 3.0.
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

## Guía de Uso (v2 con OpenAPI 3.0)

### 1. Configura OpenAPI en tu aplicación Gin

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

### 2. Define tus modelos (DTOs)

Agrega etiquetas `yaml` si planeas soportar representaciones YAML, consistente con los ejemplos OpenAPI.

```go
type UserDto struct {
	ID   int    `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

// Registra el DTO con OpenAPI components:
// _, _ = swagger.Swagger().ComponentSchemaFromDTO(&UserDto{})
```

### 3. Documenta tus endpoints (estilo OpenAPI 3.0)

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

### 4. Implementa tus funciones handler

```go
func GetItemById(c *gin.Context) {
	itemId := c.Param("itemId")
	c.JSON(http.StatusOK, gin.H{
		"message": "Item " + itemId + " found",
	})
}
```

## Documentación (v2)

Consulta el directorio `/doc_page` para documentación detallada sobre todas las características de Go Swagger Generator v2:

- [Introducción (v2)](/doc_page/docs/intro.md)
- [Inicio Rápido (v2)](/doc_page/docs/quick-start.md)
- [Definición de Modelos (v2)](/doc_page/docs/defining-models.md)
- [Parámetros de Ruta (v2)](/doc_page/docs/path-parameters.md)
- [Parámetros de Consulta (v2)](/doc_page/docs/query-parameters.md)
- [Cuerpos de Petición (v2)](/doc_page/docs/request-bodies.md)
- [Respuestas (v2)](/doc_page/docs/responses.md)
- [Esquemas de Seguridad (v2)](/doc_page/docs/security.md)
- [Respuestas de Array (v2)](/doc_page/docs/array-responses.md)
- [Configuración de Producción (v2)](/doc_page/docs/production.md)
- Y más...

## Ejemplos (v2)

Consulta el directorio `/examples` para ejemplos completos actualizados para v2:

- [Ejemplo Básico (v2)](/examples/basic/main.go)
- [Pet Store (v2)](/examples/pet_store/main.go)
- [Ejemplo de Respuesta Array (v2)](/examples/array_response/main.go)
- [Ejemplo JWT Bearer Auth (v2)](/examples/jwt_bearer/main.go)

## Probar tu API documentada

1. Inicia tu aplicación.
2. Abre tu navegador en la ruta configurada para la UI (por defecto: `http://localhost:TU_PUERTO/`, por ejemplo, http://localhost:8080/).
3. Deberías ver la interfaz Swagger UI mostrando tu API documentada con OpenAPI 3.0.
4. La especificación OpenAPI 3.0 JSON estará disponible en el `JSONPath` configurado (por ejemplo, `/openapi.json`).
5. Explora los endpoints, revisa los parámetros requeridos y prueba las llamadas directamente desde la interfaz.

## Licencia

Este proyecto está licenciado bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.

## Contribuir

¡Las contribuciones son bienvenidas! No dudes en abrir un issue o enviar un pull request.

1. Haz fork del repositorio.
2. Crea tu rama de feature (`git checkout -b feature/mi-feature`).
3. Haz commit de tus cambios (`git commit -m 'Agrega una nueva feature'`).
4. Haz push a la rama (`git push origin feature/mi-feature`).
5. Abre un Pull Request.
