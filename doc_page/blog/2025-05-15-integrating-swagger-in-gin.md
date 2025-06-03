---
slug: integrating-openapi-v2-in-gin
title: Cómo integrar OpenAPI 3.0 en tu API Gin con Go Swagger Generator v2
authors: [rui]
tags: [go, openapi, oas3, gin, api, documentation, v2]
---

# Cómo integrar OpenAPI 3.0 en tu API Gin con Go Swagger Generator v2

En el desarrollo moderno de APIs, una documentación clara y accesible es tan importante como el propio código. Una API bien documentada facilita su adopción, reduce la fricción durante la integración y ahorra tiempo tanto a desarrolladores internos como externos.

En este tutorial aprenderás a integrar OpenAPI 3.0 en una aplicación Gin usando **Go Swagger Generator v2**, una librería que facilita la generación de documentación OpenAPI directamente desde tu código Go.

<!-- truncate -->

## ¿Por qué usar Go Swagger Generator v2?

- **API fluida y elegante** - Sintaxis encadenada que facilita la documentación OpenAPI 3.0.
- **Integración simple con Gin** - Funciona con el framework Gin sin complicaciones.
- **Sin anotaciones necesarias** - No requiere comentarios especiales en tu código.
- **Swagger UI incorporado** - Incluye Swagger UI para explorar tu API de forma interactiva.

## Paso 1: Instalar dependencias

Lo primero es instalar Gin y Go Swagger Generator v2:

```bash
# Instala Gin Framework
go get github.com/gin-gonic/gin

# Instala Go Swagger Generator v2
go get -u github.com/ruiborda/go-swagger-generator/v2
```

## Paso 2: Definir modelos (DTOs)

Define una estructura simple que será parte de tu API, por ejemplo `UserDto`:

```go
type UserDto struct {
	ID   int    `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}
```

Esta estructura representa los datos de usuario que la API devolverá.

## Paso 3: Configurar la documentación OpenAPI 3.0

Crea una función dedicada para configurar la documentación:

```go
func ConfigureOpenAPI(router *gin.Engine) {
	// Habilita el middleware Swagger para OpenAPI 3.0
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
			Description("Este es un ejemplo de API usando Go Swagger Generator v2 y OpenAPI 3.0.")
	})
	doc.Server("http://localhost:8080/v1", func(server openapi.Server) {
		server.Description("Servidor local de desarrollo - API versión 1")
	})
	// Puedes agregar más servidores si lo necesitas
	// doc.Server("https://api.ejemplo.com/v1", ...)

	// Registra los DTOs en #/components/schemas/
	_, _ = doc.ComponentSchemaFromDTO(&UserDto{})
}
```

Esta función:

1. Registra el middleware Swagger en el router de Gin.
2. Configura las rutas donde se expondrá el JSON de OpenAPI y Swagger UI.
3. Define metadatos básicos de la API.
4. Configura la información de servidores.
5. Registra los DTOs para la documentación.

## Paso 4: Definir un handler para el endpoint

Define un handler simple para obtener un usuario por su ID:

```go
func GetUserById(c *gin.Context) {
	idStr := c.Param("id")
	c.JSON(http.StatusOK, UserDto{
		ID:   1,
		Name: "John Doe (id: " + idStr + ")",
	})
}
```

## Paso 5: Documentar el endpoint con sintaxis OpenAPI 3.0

Aquí es donde brilla Go Swagger Generator v2. Documenta el endpoint usando sintaxis fluida:

```go
// Documenta el endpoint GET /users/{id}
var _ = swagger.Swagger().Path("/users/{id}").
	Get(func(op openapi.Operation) {
		op.Summary("Buscar usuario por ID").
			Tag("Operaciones de Usuario").
			OperationID("getUserById").
			PathParameter("id", func(p openapi.Parameter) {
				p.Description("ID del usuario a recuperar").
					Required(true).
					Schema(func(s openapi.Schema) {
						s.Type("integer").Format("int64")
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("Operación exitosa - usuario encontrado").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&UserDto{})
					})
			}).
			Response(http.StatusNotFound, func(r openapi.Response) {
				r.Description("Usuario no encontrado")
			})
	}).
	Doc()
```

Esta documentación:

1. Define el endpoint `/users/{id}` con método GET.
2. Proporciona resumen, tag y operationId.
3. Especifica el tipo de contenido y el esquema de respuesta.
4. Documenta el parámetro de ruta `id` como entero requerido.
5. Define respuestas posibles: `200 OK` y `404 Not Found`.

El método `.Doc()` registra la documentación en la instancia OpenAPI.

## Paso 6: Implementar la función main

Finalmente, une todo en la función `main`:

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

type UserDto struct {
	ID   int    `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

func main() {
	router := gin.Default()

	ConfigureOpenAPI(router)

	router.GET("/v1/users/:id", GetUserById)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	fmt.Println("Swagger UI disponible en http://localhost:8080/")
	fmt.Println("OpenAPI 3.0 JSON disponible en http://localhost:8080/openapi.json")
	_ = router.Run(":8080")
}

// Documenta el endpoint GET /users/{id}
var _ = swagger.Swagger().Path("/users/{id}").
	Get(func(op openapi.Operation) {
		op.Summary("Buscar usuario por ID").
			Tag("Operaciones de Usuario").
			OperationID("getUserById").
			PathParameter("id", func(p openapi.Parameter) {
				p.Description("ID del usuario a recuperar").
					Required(true).
					Schema(func(s openapi.Schema) {
						s.Type("integer").Format("int64")
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("Operación exitosa - usuario encontrado").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&UserDto{})
					})
			}).
			Response(http.StatusNotFound, func(r openapi.Response) {
				r.Description("Usuario no encontrado")
			})
	}).
	Doc()

func GetUserById(c *gin.Context) {
	idStr := c.Param("id")
	c.JSON(http.StatusOK, UserDto{
		ID:   1,
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
			Description("Este es un ejemplo de API usando Go Swagger Generator v2 y OpenAPI 3.0.")
	})
	doc.Server("http://localhost:8080/v1", func(server openapi.Server) {
		server.Description("Servidor local de desarrollo - API versión 1")
	})
	_, _ = doc.ComponentSchemaFromDTO(&UserDto{})
}
```

## Probar la API documentada

1. Guarda el código anterior en un archivo `main.go`.
2. Ejecuta `go mod init tunombreproyecto` si no tienes un módulo Go.
3. Ejecuta `go mod tidy` para instalar dependencias (incluye go-swagger-generator/v2).
4. Inicia la app con `go run main.go`.
5. Abre tu navegador en [http://localhost:8080](http://localhost:8080).

Deberías ver la interfaz Swagger UI mostrando tu API documentada. Puedes explorar los endpoints y probar llamadas desde la interfaz.

## Conclusión

Integrar OpenAPI 3.0 en una API Gin usando Go Swagger Generator v2 es sencillo y aporta grandes beneficios. En minutos tendrás documentación interactiva y profesional que evoluciona junto a tu código.

¿Tienes dudas sobre cómo integrar OpenAPI 3.0 en tu API Go con v2? ¡Déjalas en los comentarios!

