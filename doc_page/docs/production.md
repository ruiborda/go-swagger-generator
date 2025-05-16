---
sidebar_position: 16
title: Production
---

# Production

When deploying to production, you should adjust Go-Swagger-Generator configuration for security and performance reasons.

## Configuring Swagger for Production Environments

```go
// Detect environment
isProd := os.Getenv("ENVIRONMENT") == "production"

// Configure Swagger middleware
router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
    Enabled: !isProd,             // Disable in production
    JSONPath: "/api-docs/swagger.json",  // Custom JSON path
    UIPath: "/api-docs",          // Custom UI path
}))
```

With this configuration, Swagger UI will be disabled in production environments but available during development. This approach helps reduce attack surface and avoids exposing API documentation unintentionally.