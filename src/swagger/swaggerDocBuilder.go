package swagger

import (
	"fmt"
	openapi2 "github.com/ruiborda/go-swagger-generator/src/openapi"
	entity2 "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
	"reflect"
	"strings"
	"sync"
)

var swaggerDoc openapi2.SwaggerDoc

type SwaggerDocBuilder struct {
	doc            *entity2.SwaggerDocEntity
	definitionsMux sync.Mutex
}

func NewSwaggerDocBuilder() openapi2.SwaggerDoc {
	if swaggerDoc != nil {
		return swaggerDoc
	}
	swaggerDoc = &SwaggerDocBuilder{
		doc: &entity2.SwaggerDocEntity{
			Swagger:             "2.0",
			Info:                entity2.InfoEntity{},
			Paths:               make(map[string]entity2.PathItemEntity),
			Definitions:         make(map[string]entity2.SchemaEntity),
			Tags:                make([]entity2.TagEntity, 0),
			Schemes:             make([]string, 0),
			SecurityDefinitions: make(map[string]entity2.SecuritySchemeEntity),
		},
	}
	return swaggerDoc
}

func (b *SwaggerDocBuilder) SwaggerVersion(version string) openapi2.SwaggerDoc {
	b.doc.Swagger = version
	return b
}

func (b *SwaggerDocBuilder) Info(config func(builder openapi2.Info)) openapi2.SwaggerDoc {
	infoBuilder := &InfoBuilder{info: &b.doc.Info}
	config(infoBuilder)
	return b
}

func (b *SwaggerDocBuilder) Host(host string) openapi2.SwaggerDoc {
	b.doc.Host = host
	return b
}

func (b *SwaggerDocBuilder) BasePath(basePath string) openapi2.SwaggerDoc {
	b.doc.BasePath = basePath
	return b
}

func (b *SwaggerDocBuilder) Tag(name string, config func(builder openapi2.Tag)) openapi2.SwaggerDoc {
	tag := entity2.TagEntity{Name: name}
	tagBuilder := &TagBuilder{tag: &tag}
	config(tagBuilder)
	b.doc.Tags = append(b.doc.Tags, tag)
	return b
}

func (b *SwaggerDocBuilder) Scheme(scheme string) openapi2.SwaggerDoc {
	b.doc.Schemes = append(b.doc.Schemes, scheme)
	return b
}

func (b *SwaggerDocBuilder) Schemes(schemes ...string) openapi2.SwaggerDoc {
	b.doc.Schemes = append(b.doc.Schemes, schemes...)
	return b
}

func (b *SwaggerDocBuilder) Path(pathPattern string) openapi2.PathItem {
	pathItem, exists := b.doc.Paths[pathPattern]
	if !exists {
		pathItem = entity2.PathItemEntity{}
	}
	return &PathItemBuilder{
		pathItem:   &pathItem,
		docPath:    pathPattern,
		docBuilder: b,
	}
}

func (b *SwaggerDocBuilder) SecurityDefinition(name string, config func(builder openapi2.SecurityScheme)) openapi2.SwaggerDoc {
	secScheme := entity2.SecuritySchemeEntity{}
	secSchemeBuilder := &SecuritySchemeBuilder{scheme: &secScheme}
	config(secSchemeBuilder)
	if b.doc.SecurityDefinitions == nil {
		b.doc.SecurityDefinitions = make(map[string]entity2.SecuritySchemeEntity)
	}
	b.doc.SecurityDefinitions[name] = secScheme
	return b
}

func (b *SwaggerDocBuilder) Definition(name string, schema entity2.SchemaEntity) openapi2.SwaggerDoc {
	b.definitionsMux.Lock()
	defer b.definitionsMux.Unlock()
	if b.doc.Definitions == nil {
		b.doc.Definitions = make(map[string]entity2.SchemaEntity)
	}
	b.doc.Definitions[name] = schema
	return b
}

func (b *SwaggerDocBuilder) DefinitionFromDTO(dtoInstance interface{}) (string, error) {
	b.definitionsMux.Lock()
	defer b.definitionsMux.Unlock()
	if b.doc.Definitions == nil {
		b.doc.Definitions = make(map[string]entity2.SchemaEntity)
	}

	dtoType := reflect.TypeOf(dtoInstance)
	if dtoType.Kind() == reflect.Ptr {
		dtoType = dtoType.Elem()
	}
	if dtoType.Kind() != reflect.Struct {
		return "", fmt.Errorf("DTO must be a struct or pointer to struct, got %s", dtoType.Kind())
	}

	dtoName := dtoType.Name()
	if _, exists := b.doc.Definitions[dtoName]; !exists {
		_, err := b.GenerateSchemaFromGoType(dtoType, make(map[string]bool))
		if err != nil {
			return "", fmt.Errorf("failed to generate schema for DTO %s: %w", dtoName, err)
		}
	}
	return dtoName, nil
}

func (b *SwaggerDocBuilder) ExternalDocumentation(url string, description string) openapi2.SwaggerDoc {
	b.doc.ExternalDocs = &entity2.ExternalDocumentationEntity{URL: url, Description: description}
	return b
}

func (b *SwaggerDocBuilder) Build() entity2.SwaggerDocEntity {
	return *b.doc
}

func (b *SwaggerDocBuilder) GenerateSchemaFromGoType(t reflect.Type, visited map[string]bool) (*entity2.SchemaEntity, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Handle recursive types by creating a reference if already visited
	if t.Kind() == reflect.Struct {
		typeName := t.PkgPath() + "." + t.Name() // Unique identifier for the type
		if visited[typeName] {
			// If this type name is already in Definitions, it's a known DTO
			if _, exists := b.doc.Definitions[t.Name()]; exists {
				return &entity2.SchemaEntity{Ref: "#/definitions/" + t.Name()}, nil
			}
			// Otherwise, it's a recursive call within the same DTO generation, return a temporary ref
			// This case might need more robust handling if complex anonymous struct recursions are expected
			return &entity2.SchemaEntity{Ref: "#/definitions/" + t.Name()}, nil // Hope t.Name() is unique enough
		}
		visited[typeName] = true
		defer delete(visited, typeName) // Clean up after processing this type
	}

	schema := &entity2.SchemaEntity{}

	switch t.Kind() {
	case reflect.String:
		schema.Type = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		schema.Type = "integer"
		if t.Kind() == reflect.Int64 || t.Kind() == reflect.Uint64 {
			schema.Format = "int64"
		} else {
			schema.Format = "int32"
		}
	case reflect.Float32:
		schema.Type = "number"
		schema.Format = "float"
	case reflect.Float64:
		schema.Type = "number"
		schema.Format = "double"
	case reflect.Bool:
		schema.Type = "boolean"
	case reflect.Slice, reflect.Array:
		schema.Type = "array"
		elemType := t.Elem()
		// Check if the element type is the same as the parent type (recursive array)
		if elemType.Kind() == reflect.Ptr { // Dereference pointer if element is a pointer
			elemType = elemType.Elem()
		}
		// Simple self-reference check for arrays/slices of the DTO itself
		if elemType.Name() == t.Name() && elemType.PkgPath() == t.PkgPath() {
			schema.Items = &entity2.SchemaEntity{Ref: "#/definitions/" + elemType.Name()}
		} else {
			itemSchema, err := b.GenerateSchemaFromGoType(elemType, visited)
			if err != nil {
				return nil, fmt.Errorf("failed to generate item schema for array/slice of %s: %w", elemType.Name(), err)
			}
			schema.Items = itemSchema
		}
	case reflect.Struct:
		// Handle special cases like time.Time
		if t.PkgPath() == "time" && t.Name() == "Time" {
			schema.Type = "string"
			schema.Format = "date-time"
			return schema, nil
		}

		// This struct will be a definition
		dtoName := t.Name()
		schema.Ref = "#/definitions/" + dtoName

		// If this definition doesn't exist yet, create it
		if _, exists := b.doc.Definitions[dtoName]; !exists {
			// Temporarily add to definitions to handle self-references within fields
			// b.doc.Definitions[dtoName] = SchemaEntity{Ref: "#/definitions/" + dtoName} // Placeholder for recursion

			fullStructSchema := entity2.SchemaEntity{
				Type:       "object",
				Properties: make(map[string]*entity2.SchemaEntity),
				Required:   []string{},
			}
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				if !field.IsExported() { // Skip unexported fields
					continue
				}

				jsonTag := field.Tag.Get("json")
				jsonFieldName := field.Name // Default to field name
				omitempty := false
				if jsonTag != "" {
					parts := strings.Split(jsonTag, ",")
					if parts[0] == "-" { // Field is ignored
						continue
					}
					if parts[0] != "" {
						jsonFieldName = parts[0]
					}
					for _, part := range parts[1:] {
						if part == "omitempty" {
							omitempty = true
							break
						}
					}
				}

				propSchema, err := b.GenerateSchemaFromGoType(field.Type, visited)
				if err != nil {
					return nil, fmt.Errorf("failed to generate schema for field %s in struct %s: %w", field.Name, dtoName, err)
				}
				fullStructSchema.Properties[jsonFieldName] = propSchema
				if !omitempty { // Add to required if not omitempty
					fullStructSchema.Required = append(fullStructSchema.Required, jsonFieldName)
				}
			}
			if len(fullStructSchema.Required) == 0 {
				fullStructSchema.Required = nil // omit if empty
			}
			b.doc.Definitions[dtoName] = fullStructSchema // Replace placeholder with full schema
		}
	case reflect.Map:
		schema.Type = "object"
		// Swagger 2.0 only supports string keys for maps.
		// For `additionalProperties`, we use the schema of the map's value type.
		valType := t.Elem()
		addPropsSchema, err := b.GenerateSchemaFromGoType(valType, visited)
		if err != nil {
			return nil, fmt.Errorf("failed to generate schema for map value type %s: %w", valType.Name(), err)
		}
		schema.AdditionalProperties = addPropsSchema

	default:
		return nil, fmt.Errorf("unsupported type for DTO schema generation: %s", t.Kind())
	}

	return schema, nil
}
