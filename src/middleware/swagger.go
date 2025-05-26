package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/src/swagger" // Updated import path if Swagger() is in swagger package directly
	"html/template"
	"net/http"
)

// SwaggerConfig holds configuration for the SwaggerGin middleware
type SwaggerConfig struct {
	// Enabled determines if the SwaggerGin UI is enabled
	Enabled bool
	// JSONPath is the path where the SwaggerGin JSON will be served
	JSONPath string
	// UIPath is the path where the SwaggerGin UI will be served
	UIPath string
	// Title for the Swagger UI page
	Title string
}

// DefaultSwaggerConfig returns the default SwaggerGin configuration
func DefaultSwaggerConfig() SwaggerConfig {
	return SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json",
		UIPath:   "/",
	}
}

// SwaggerGin returns a gin middleware for serving Swagger UI and JSON for OAS3
func SwaggerGin(config ...SwaggerConfig) gin.HandlerFunc {
	cfg := DefaultSwaggerConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Use swagger.Swagger() to get the builder instance
	docBuilder := swagger.Swagger()

	swaggerTemplate := template.Must(template.New("swagger").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="description" content="SwaggerUI" />
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
<style>body{margin:0;padding:0;}</style>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-standalone-preset.js" crossorigin></script>
<script>
    window.onload = () => {
        window.ui = SwaggerUIBundle({
            url: '{{.JSONPath}}',
            dom_id: '#swagger-ui',
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIStandalonePreset
            ],
            layout: "StandaloneLayout",
        });
    };
</script>
</body>
</html>`))

	return func(c *gin.Context) {
		reqPath := c.Request.URL.Path

		if reqPath == cfg.JSONPath {
			c.Header("Content-Type", "application/json")
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
			c.JSON(http.StatusOK, docBuilder.Build())
			c.Abort()
			return
		}

		if reqPath == cfg.UIPath || (cfg.UIPath == "/" && reqPath == "/") {
			data := struct {
				Title    string
				JSONPath string
			}{
				Title:    cfg.Title,
				JSONPath: cfg.JSONPath,
			}

			c.Header("Content-Type", "text/html; charset=utf-8")
			if err := swaggerTemplate.Execute(c.Writer, data); err != nil {
				_ = c.Error(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.Abort()
			return
		}

		c.Next()
	}
}
