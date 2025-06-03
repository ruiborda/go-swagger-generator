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

type LoginRequest struct {
	Username string `json:"username" binding:"required" yaml:"username"`
	Password string `json:"password" binding:"required" yaml:"password"`
}

type TokenResponse struct {
	Token   string `json:"token" yaml:"token"`
	Type    string `json:"type" yaml:"type"`             // e.g. Bearer
	Expires int    `json:"expires_in" yaml:"expires_in"` // e.g. 3600 (seconds)
}

func main() {
	router := gin.Default()

	// Register DTOs for components.schemas
	_, _ = swagger.Swagger().SchemaFromDTO(&LoginRequest{})
	_, _ = swagger.Swagger().SchemaFromDTO(&TokenResponse{})

	router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json",
		UIPath:   "/",
		Title:    "JWT Bearer Auth API OAS3",
	}))
	configureSwagger()

	router.POST("/v1/auth/login", handleLogin)
	// Example of a secured route (handler not implemented for brevity)
	// router.GET("/v1/secure/data", AuthMiddleware(), handleSecureData)

	fmt.Println("Swagger UI available at http://localhost:8080/")
	fmt.Println("Swagger JSON available at http://localhost:8080/openapi.json")
	_ = router.Run(":8080")
}

func handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de solicitud inválido"})
		return
	}

	if req.Username != "user" || req.Password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		Type:    "Bearer",
		Expires: 3600,
	})
}

func configureSwagger() {
	doc := swagger.Swagger()

	doc.Info(func(info openapi.Info) {
		info.Title("API de Autenticación JWT con OAS3").
			Version("1.0").
			Description("Ejemplo API con autenticación JWT Bearer (OpenAPI 3.0)")
	}).
		Server("http://localhost:8080/v1", func(s openapi.Server) {
			s.Description("Servidor de Desarrollo V1")
		})

	// Definición del esquema de seguridad JWT Bearer
	doc.ComponentSecurityScheme("BearerAuth", func(ss openapi.SecurityScheme) {
		ss.Type("http").
			Scheme("bearer").
			BearerFormat("JWT"). // Optional, for documentation
			Description("Autenticación JWT Bearer. Ingrese 'Bearer {token}' en el campo de autorización.")
	})

	// Documentación del endpoint POST /auth/login (relative to server URL)
	_ = doc.Path("/auth/login").Post(func(op openapi.Operation) {
		op.Summary("Login con usuario y contraseña").
			Description("Autentica al usuario y devuelve un token JWT Bearer.").
			Tag("Autenticación").
			OperationID("loginUser").
			RequestBody(func(rb openapi.RequestBody) {
				rb.Description("Credenciales de acceso").
					Required(true).
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&LoginRequest{})
					})
			}).
			Response(http.StatusOK, func(r openapi.Response) {
				r.Description("Autenticación exitosa, token JWT emitido.").
					Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
						mt.SchemaFromDTO(&TokenResponse{})
					})
			}).
			Response(http.StatusBadRequest, func(r openapi.Response) {
				r.Description("Formato de solicitud inválido")
			}).
			Response(http.StatusUnauthorized, func(r openapi.Response) {
				r.Description("Credenciales incorrectas")
			})
	})

	doc.GlobalSecurity(map[string][]string{
		"BearerAuth": {},
	})

}
