---
sidebar_position: 16
title: Production Configuration (v2)
---

# Production Configuration (v2)

When deploying your Gin application with Go-Swagger-Generator v2 to production, you might want to adjust the OpenAPI documentation Ccnfiguration for security and performance reasons. Typically, you might disable or restrict access to the Swagger UI and the OpenAPI JSON file in a production environment.

## Configuring Swagger Middleware for Production

The `middleware.SwaggerConfig` allows you to control the behavior of the Swagger UI and JSON endpoint.

```go
package main

import (
    "os"
    "github.com/gin-gonic/gin"
    "github.com/ruiborda/go-swagger-generator/src/middleware"
    // ... other necessary imports
)

func main() {
    router := gin.Default()

    // Detect if running in a production environment
    // This is just one way; you might use build tags or other config methods.
    isProd := os.Getenv("APP_ENV") == "production"

    // Configure Swagger middleware
    router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
        Enabled:  !isProd,          // Disable Swagger UI and JSON endpoint if isProd is true.
        JSONPath: "/openapi.json",  // Path for the OpenAPI JSON file (even if disabled, good to define).
        UIPath:   "/docs",          // Path for the Swagger UI (e.g., serve it at /docs).
                                     // If UIPath is "/", it serves at the root.
        Title:    "My Production API Docs", // Title for the Swagger UI page.
    }))

    // ... your existing OpenAPI documentation setup using swagger.Swagger() ...
    // ... your Gin route definitions ...

    _ = router.Run(":8080")
}
```

**Explanation:**

*   `Enabled: !isProd`: This line is key. If `isProd` is `true` (meaning the application is running in a production environment), `Enabled` becomes `false`. When `Enabled` is `false`, the `SwaggerGin` middleware will not serve the Swagger UI or the OpenAPI JSON file. This prevents public exposure of your API's detailed structure in production.
*   `JSONPath: "/openapi.json"`: Defines the path where the OpenAPI JSON specification would be served if enabled. This path is also used by the Swagger UI to fetch the specification.
*   `UIPath: "/docs"`: Defines the path where the Swagger UI would be served if enabled. You can customize this to any path, like `/api-documentation`, `/swagger`, or keep it at `/` for the root.
*   `Title: "My Production API Docs"`: Sets the title of the Swagger UI page. This is useful even if UI is conditionally disabled, as it's part of the configuration.

**Alternative Strategies for Production:**

1.  **Serve Statically:** Instead of disabling, you could generate the `openapi.json` file during your build process and serve it as a static file, potentially protected by authentication/authorization if needed. The Swagger UI itself can also be hosted statically.
2.  **Internal Access Only:** Configure your reverse proxy (e.g., Nginx, Caddy) or firewall to only allow access to the Swagger UI/JSON paths from internal IP addresses or specific authenticated users.
3.  **Separate Documentation Server:** Host your Swagger UI and OpenAPI specification on a completely separate, internally accessible documentation server.

Choosing the right strategy depends on your organization's security policies and operational requirements. Disabling it via the `Enabled` flag is the simplest way to prevent public exposure directly from the application itself.