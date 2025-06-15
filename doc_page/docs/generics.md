---
sidebar_position: 17
title: Generic DTO Support (v2)
---

# Generic DTO Support (v2)

Go Swagger Generator v2 supports generating OpenAPI schemas from Go generic types (structs). This allows you to define reusable generic Data Transfer Objects (DTOs) and have their specific instantiations correctly represented in your OpenAPI documentation.

## Defining Generic DTOs

You can define generic DTOs in Go as you normally would using type parameters. For example, a common use case is a generic API response wrapper:

```go
package main // or your DTO package

// GenericResponse is a generic DTO for API responses.
// The type parameter T will be replaced by a concrete type upon instantiation.

// @schema(description="A generic API response wrapper.")
// @property(name="success", type="boolean", description="Indicates if the request was successful.")
// @property(name="message", type="string", format="", description="Optional message, typically for errors or info.")
// @property(name="data", description="The actual data payload, type varies based on T.")
type GenericResponse[T any] struct {
    Success bool   `json:"success"`
    Message string `json:"message,omitempty"`
    Data    T      `json:"data,omitempty"`
}

// UserData is a concrete DTO for user information.

// @schema(description="Detailed user information.")
// @property(name="id", type="integer", format="int32", description="User's unique identifier.")
// @property(name="username", type="string", description="User's login name.")
// @property(name="email", type="string", format="email", description="User's email address.")
type UserData struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

// ProductData is another concrete DTO for product information.

// @schema(description="Detailed product information.")
// @property(name="sku", type="string", description="Stock Keeping Unit.")
// @property(name="productName", type="string", description="Name of the product.")
// @property(name="price", type="number", format="double", description="Price of the product.")
type ProductData struct {
    SKU         string  `json:"sku"`
    ProductName string  `json:"productName"`
    Price       float64 `json:"price"`
}
```

## Using `SchemaFromDTO` with Generic Types

When you use `doc.SchemaFromDTO()` (or its variants like `mt.SchemaFromDTO()`, `p.SchemaFromDTO()`, etc.) with an instantiated generic type, the library will generate a unique component schema name and definition for that specific instantiation.

```go
// In your ConfigureOpenAPI function:
doc := swagger.Swagger() 

// Registering an instance of GenericResponse[UserData]
_, _ = doc.SchemaFromDTO(&GenericResponse[UserData]{})

// Registering an instance of GenericResponse[ProductData]
_, _ = doc.SchemaFromDTO(&GenericResponse[ProductData]{})

// Registering an instance of GenericResponse with a slice type argument
_, _ = doc.SchemaFromDTO(&GenericResponse[[]*UserData]{})

// It's also good practice to register the concrete types if they are used elsewhere
_, _ = doc.SchemaFromDTO(&UserData{})
_, _ = doc.SchemaFromDTO(&ProductData{})
```

### Generated Component Names

The generator will create unique names for these instantiated generic schemas in your `#/components/schemas/` section. The naming convention aims to combine the generic type's name with its type arguments. For example:

-   `&GenericResponse[UserData]{}` might generate a schema named `main_GenericResponse_main_UserData` (if defined in `main` package) or similar, reflecting the package and type names.
-   `&GenericResponse[ProductData]{}` could become `main_GenericResponse_main_ProductData`.
-   `&GenericResponse[[]*UserData]{}` might be named `main_GenericResponse_ListOf_main_UserData` or `main_GenericResponse_ListOf_Ptr_main_UserData`.

The exact naming might vary slightly based on sanitization rules but will be unique for each instantiation.

### Schema Structure

The generated schema for `main_GenericResponse_main_UserData` would look something like this (in JSON format):

```json
{
  "type": "object",
  "properties": {
    "success": {
      "type": "boolean"
    },
    "message": {
      "type": "string"
    },
    "data": {
      "$ref": "#/components/schemas/main_UserData" // Reference to the UserData schema
    }
  },
  "required": ["success"] // Based on Go struct tags (omitempty)
}
```

And for `main_GenericResponse_ListOf_main_UserData`:

```json
{
  "type": "object",
  "properties": {
    "success": {
      "type": "boolean"
    },
    "message": {
      "type": "string"
    },
    "data": {
      "type": "array",
      "items": {
        "$ref": "#/components/schemas/main_UserData" // Reference to UserData
      }
    }
  },
  "required": ["success"]
}
```

## Example Usage in an Endpoint Definition

```go
// Documenting an endpoint that returns GenericResponse[UserData]
doc.Path("/userinfo").
    Get(func(op openapi.Operation) {
        op.Summary("Get User Information").
            Tag("Generic Examples").
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation - user data returned").
                    Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                        // Use the instantiated generic type here
                        mt.SchemaFromDTO(&GenericResponse[UserData]{})
                    })
            })
    }).Doc()

// Documenting an endpoint that returns GenericResponse[[]*UserData]
doc.Path("/userslist").
    Get(func(op openapi.Operation) {
        op.Summary("Get List of Users Information").
            Tag("Generic Examples").
            Response(http.StatusOK, func(r openapi.Response) {
                r.Description("Successful operation - list of user data returned").
                    Content(mime.ApplicationJSON, func(mt openapi.MediaType) {
                        mt.SchemaFromDTO(&GenericResponse[[]*UserData]{})
                    })
            })
    }).Doc()
```

## Key Considerations

1.  **Instantiate Generics**: Always pass an *instance* of the fully instantiated generic type to `SchemaFromDTO` (e.g., `&GenericResponse[MyConcreteType]{}`).
2.  **Type Argument Registration**: If type arguments (like `MyConcreteType` in the example above) are complex structs themselves, ensure they are also registered via `SchemaFromDTO` so they can be correctly referenced.
3.  **Naming Uniqueness**: The library strives to create unique component names. If you encounter issues with name collisions or excessively long names in complex scenarios, review your type structures or consider manual schema definitions for highly complex cases.
4.  **Reflection Limitations**: Go's reflection on generics has some limitations compared to non-generic types. The current support focuses on instantiated generic structs. More complex generic scenarios (e.g., generic functions, methods with generic receivers used in schema generation) might not be fully supported for schema inference.

By leveraging generic DTOs, you can significantly reduce boilerplate code in your API definitions while maintaining clear and accurate OpenAPI documentation.