package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/src/swagger"
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
}

// DefaultSwaggerConfig returns the default SwaggerGin configuration
func DefaultSwaggerConfig() SwaggerConfig {
	return SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json",
		UIPath:   "/",
	}
}

// SwaggerGin returns a gin middleware for serving SwaggerGin UI and JSON
func SwaggerGin(config ...SwaggerConfig) gin.HandlerFunc {
	// Use default config if none provided
	cfg := DefaultSwaggerConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// If disabled, return an empty middleware
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// HTML template for SwaggerGin UI
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
		// Skip this middleware if the route doesn't match
		reqPath := c.Request.URL.Path

		// Check if the request is for the SwaggerGin JSON
		if reqPath == cfg.JSONPath {
			c.Header("Content-Type", "application/json")
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
			c.JSON(http.StatusOK, swagger.NewSwaggerDocBuilder().Build())
			c.Abort()
			return
		}

		// Check if the request is for the SwaggerGin UI
		if reqPath == cfg.UIPath {
			data := struct {
				Title    string
				JSONPath string
			}{
				Title:    "SwaggerUI",
				JSONPath: cfg.JSONPath,
			}

			c.Header("Content-Type", "text/html")
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
