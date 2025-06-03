package swagger

import (
	"fmt"
	openapi "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
	entity "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec"
	"reflect"
	"strings"
	"sync"
)

var swaggerDoc openapi.SwaggerDocBuilder
var once sync.Once

type SwaggerDocBuilder struct {
	doc            *entity.SwaggerDocEntity
	definitionsMux sync.Mutex // Used for DTO to Schema conversion
}

// Swagger creates a new SwaggerDocBuilder instance for OAS3.
func Swagger() openapi.SwaggerDocBuilder {
	once.Do(func() {
		docEntity := &entity.SwaggerDocEntity{
			Openapi: "3.1.0",
			Info:    entity.Info{},
			Paths:   make(map[string]*entity.PathItem),
			Servers: make([]*entity.Server, 0),
			Tags:    make([]*entity.Tag, 0),
			Components: &entity.Components{
				Schemas:         make(map[string]*entity.SchemaRef),
				Responses:       make(map[string]*entity.ResponseRef),
				Parameters:      make(map[string]*entity.ParameterRef),
				Examples:        make(map[string]*entity.ExampleRef),
				RequestBodies:   make(map[string]*entity.RequestBodyRef),
				Headers:         make(map[string]*entity.HeaderRef),
				SecuritySchemes: make(map[string]*entity.SecuritySchemeRef),
				Links:           make(map[string]*entity.LinkRef),
				Callbacks:       make(map[string]*entity.CallbackRef),
			},
			Security: make([]entity.SecurityRequirement, 0),
		}
		swaggerDoc = &SwaggerDocBuilder{doc: docEntity}
	})
	return swaggerDoc
}

func (b *SwaggerDocBuilder) OpenAPIVersion(version string) openapi.SwaggerDocBuilder {
	b.doc.Openapi = version
	return b
}

func (b *SwaggerDocBuilder) Info(config func(builder openapi.Info)) openapi.SwaggerDocBuilder {
	infoBuilder := &InfoBuilder{info: &b.doc.Info}
	config(infoBuilder)
	return b
}

func (b *SwaggerDocBuilder) Server(url string, config func(builder openapi.Server)) openapi.SwaggerDocBuilder {
	server := entity.Server{URL: url}
	serverBuilder := &ServerBuilder{server: &server}
	config(serverBuilder)
	b.doc.Servers = append(b.doc.Servers, &server)
	return b
}

func (b *SwaggerDocBuilder) Servers(servers ...entity.Server) openapi.SwaggerDocBuilder {
	for i := range servers {
		b.doc.Servers = append(b.doc.Servers, &servers[i])
	}
	return b
}

func (b *SwaggerDocBuilder) Tag(name string, config func(builder openapi.Tag)) openapi.SwaggerDocBuilder {
	tag := entity.Tag{Name: name}
	tagBuilder := &TagBuilder{tag: &tag}
	config(tagBuilder)
	b.doc.Tags = append(b.doc.Tags, &tag)
	return b
}

func (b *SwaggerDocBuilder) Path(pathPattern string) openapi.PathItem {
	if b.doc.Components == nil {
		b.doc.Components = &entity.Components{}
	}
	pathItem, exists := b.doc.Paths[pathPattern]
	if !exists || pathItem == nil {
		pathItem = &entity.PathItem{}
		b.doc.Paths[pathPattern] = pathItem
	}
	return &PathItemBuilder{
		pathItem:   pathItem,
		docPath:    pathPattern,
		docBuilder: b,
	}
}

func (b *SwaggerDocBuilder) ComponentSchema(name string, schema entity.Schema) openapi.SwaggerDocBuilder {
	b.definitionsMux.Lock()
	defer b.definitionsMux.Unlock()
	if b.doc.Components.Schemas == nil {
		b.doc.Components.Schemas = make(map[string]*entity.SchemaRef)
	}
	b.doc.Components.Schemas[name] = &entity.SchemaRef{Schema: &schema}
	return b
}

func (b *SwaggerDocBuilder) ComponentSchemaRef(name string, ref string) openapi.SwaggerDocBuilder {
	b.definitionsMux.Lock()
	defer b.definitionsMux.Unlock()
	if b.doc.Components.Schemas == nil {
		b.doc.Components.Schemas = make(map[string]*entity.SchemaRef)
	}
	b.doc.Components.Schemas[name] = &entity.SchemaRef{Ref: ref}
	return b
}

// SchemaFromDTO generates a schema from a DTO and adds it to components.schemas.
// It returns the name of the generated schema.
func (b *SwaggerDocBuilder) SchemaFromDTO(dtoInstance interface{}) (string, error) {
	b.definitionsMux.Lock()
	defer b.definitionsMux.Unlock()

	if b.doc.Components == nil {
		b.doc.Components = &entity.Components{}
	}
	if b.doc.Components.Schemas == nil {
		b.doc.Components.Schemas = make(map[string]*entity.SchemaRef)
	}
	// processedInThisCall is a map to track types processed within this top-level SchemaFromDTO call
	// to avoid re-processing the same component if encountered multiple times or through recursion.
	return b.schemaFromDTORecursive(dtoInstance, make(map[reflect.Type]string))
}

// schemaFromDTORecursive is the internal implementation for SchemaFromDTO.
// It assumes definitionsMux is already locked.
// processedInThisCall tracks types (reflect.Type) to their component names (string)
// that have already been initiated or completed within the *current* top-level SchemaFromDTO call chain.
func (b *SwaggerDocBuilder) schemaFromDTORecursive(dtoInstance interface{}, processedInThisCall map[reflect.Type]string) (string, error) {
	originalDtoType := reflect.TypeOf(dtoInstance)
	currentDtoType := originalDtoType

	// Handle nil input gracefully, perhaps by returning an error or a specific schema name for null type.
	if dtoInstance == nil {
		// Depending on desired behavior, could define a generic 'Null' schema or similar.
		// For now, error out or treat as an issue, as DTOs are expected to be typed.
		return "", fmt.Errorf("SchemaFromDTO called with nil DTO instance")
	}

	if currentDtoType.Kind() == reflect.Ptr {
		// If pointer is nil, create a new instance of the pointed-to type to proceed with type analysis.
		// This is important for recursive calls like `elementInstance` for slices.
		if reflect.ValueOf(dtoInstance).IsNil() {
			dtoInstance = reflect.New(currentDtoType.Elem()).Interface()
		}
		currentDtoType = currentDtoType.Elem()
	}

	// Check if this exact original type has already been processed or started in this call chain.
	if componentName, ok := processedInThisCall[originalDtoType]; ok {
		return componentName, nil
	}

	// Handle slices/arrays
	if currentDtoType.Kind() == reflect.Slice || currentDtoType.Kind() == reflect.Array {
		elementType := currentDtoType.Elem() // Type of the elements in the slice/array
		baseElementType := elementType
		if baseElementType.Kind() == reflect.Ptr {
			baseElementType = baseElementType.Elem()
		}

		if baseElementType.Kind() != reflect.Struct {
			return "", fmt.Errorf("slice/array element DTO must be a struct or pointer to struct, got %s for element of %s", baseElementType.Kind(), currentDtoType.String())
		}

		// Create a zero-value instance of the base element type for recursive call.
		elementInstance := reflect.New(baseElementType).Interface()

		elementComponentName, err := b.schemaFromDTORecursive(elementInstance, processedInThisCall)
		if err != nil {
			return "", fmt.Errorf("failed to generate schema for slice/array element type %s: %w", baseElementType.Name(), err)
		}

		listDtoName := "ListOf" + elementComponentName

		// Check if this list component (e.g., "ListOfUserResponse") is already globally defined.
		if _, exists := b.doc.Components.Schemas[listDtoName]; exists {
			processedInThisCall[originalDtoType] = listDtoName // Cache for this call chain.
			return listDtoName, nil
		}

		arraySchema := &entity.Schema{
			Type: "array",
			Items: &entity.SchemaRef{
				Ref: "#/components/schemas/" + elementComponentName,
			},
		}
		b.doc.Components.Schemas[listDtoName] = &entity.SchemaRef{Schema: arraySchema}
		processedInThisCall[originalDtoType] = listDtoName // Cache for this call chain.
		return listDtoName, nil
	}

	// Handle structs
	if currentDtoType.Kind() == reflect.Struct {
		structDtoName := currentDtoType.Name()
		if structDtoName == "" {
			return "", fmt.Errorf("anonymous structs cannot be registered as top-level DTOs via SchemaFromDTO")
		}

		// Check if this struct component is already globally defined (e.g., from a previous SchemaFromDTO call or deeper recursion).
		if _, exists := b.doc.Components.Schemas[structDtoName]; exists {
			// If it exists, its processing has either completed or is in progress higher up the stack.
			// processedInThisCall check at the beginning handles if we started it in this chain.
			// If not in processedInThisCall, it's either fully done globally or a placeholder from recursion.
			// In any case, its name is now known and usable.
			processedInThisCall[originalDtoType] = structDtoName
			return structDtoName, nil
		}

		// New struct component for global registration: add placeholder and mark as processed for this call chain.
		b.doc.Components.Schemas[structDtoName] = &entity.SchemaRef{Ref: "#/components/schemas/" + structDtoName} // Placeholder for recursion
		processedInThisCall[originalDtoType] = structDtoName                                                      // Mark as started in this chain

		generatedSchema, err := b.generateSchemaFromGoType(currentDtoType, make(map[string]bool), structDtoName)
		if err != nil {
			delete(b.doc.Components.Schemas, structDtoName) // Clean up placeholder
			delete(processedInThisCall, originalDtoType)    // Clean up cache entry for this attempt
			return "", fmt.Errorf("failed to generate schema for DTO struct %s: %w", structDtoName, err)
		}

		// generateSchemaFromGoType for the top-level struct (structDtoName) should NOT return a Ref itself.
		if generatedSchema.Ref != "" {
			delete(b.doc.Components.Schemas, structDtoName) // Clean up placeholder
			delete(processedInThisCall, originalDtoType)
			return "", fmt.Errorf("internal error: generateSchemaFromGoType for component %s returned a $ref ('%s'), expected full schema definition", structDtoName, generatedSchema.Ref)
		}

		// Update placeholder with the fully generated schema.
		b.doc.Components.Schemas[structDtoName] = &entity.SchemaRef{Schema: generatedSchema}
		return structDtoName, nil
	}

	// If not a slice, array, or struct, it's an unsupported DTO type for top-level registration.
	return "", fmt.Errorf("DTO must be a struct, pointer to struct, slice/array of structs, or pointer to slice/array of structs, got %s", originalDtoType.String())
}

// generateSchemaFromGoType converts a Go type to an OpenAPI Schema object.
// visited map is used to handle recursive types for the *current* struct's fields.
// currentlyDefiningCompName is the name of the component schema that SchemaFromDTO is currently trying to define.
// This is used to prevent generateSchemaFromGoType from short-circuiting with a $ref for the component it's meant to be defining.
func (b *SwaggerDocBuilder) generateSchemaFromGoType(t reflect.Type, visited map[string]bool, currentlyDefiningCompName string) (*entity.Schema, error) {
	originalTypeFullName := t.PkgPath() + "." + t.Name() // For visited map key

	// Dereference pointer types to get the underlying type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Handle named types (structs, interfaces potentially) that might be components or part of a recursion cycle.
	if t.Name() != "" && t.Kind() == reflect.Struct { // Only consider structs for component/recursion logic here
		// 1. Cycle detection within the current struct's field expansion:
		// If this exact type (package + name) is already in the visited set for the current struct's fields.
		if visited[originalTypeFullName] {
			return &entity.Schema{Ref: "#/components/schemas/" + t.Name()}, nil
		}

		// 2. Pre-existing component or different recursive struct:
		// If this type is NOT the one being primarily defined by SchemaFromDTO (i.e., t.Name() != currentlyDefiningCompName)
		// AND it IS already registered as a component (e.g., by a prior SchemaFromDTO call or another branch of current one).
		if t.Name() != currentlyDefiningCompName {
			if _, isComponent := b.doc.Components.Schemas[t.Name()]; isComponent {
				// This refers to a different, already known component. Safe to return $ref.
				return &entity.Schema{Ref: "#/components/schemas/" + t.Name()}, nil
			}
		}

		// If we are here, it means for a named struct:
		// - It's NOT a cycle in the current field expansion (not in `visited`).
		// - AND (it IS the `currentlyDefiningCompName` OR it's another struct not yet a component).
		// In these cases, we need to expand its fields.
		// Add to visited map for this struct's expansion cycle detection.
		visited[originalTypeFullName] = true
		defer delete(visited, originalTypeFullName) // Remove from visited after its fields are processed
	}

	// Special handling for time.Time
	if t.PkgPath() == "time" && t.Name() == "Time" && t.Kind() == reflect.Struct {
		return &entity.Schema{Type: "string", Format: "date-time"}, nil
	}

	schema := &entity.Schema{}

	switch t.Kind() {
	case reflect.String:
		schema.Type = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		schema.Type = "integer"
		schema.Format = "int32"
	case reflect.Int64:
		schema.Type = "integer"
		schema.Format = "int64"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		schema.Type = "integer"
		schema.Format = "int32" // OpenAPI does not have uint32, use int32 with minimum: 0
		minValue := 0.0
		schema.Minimum = &minValue
	case reflect.Uint64:
		schema.Type = "integer"
		schema.Format = "int64" // OpenAPI does not have uint64, use int64 with minimum: 0
		minValue := 0.0
		schema.Minimum = &minValue
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

		// Recursively generate schema for the element type, passing current `visited` and `currentlyDefiningCompName`.
		// If elemType is a struct that matches `currentlyDefiningCompName`, generateSchemaFromGoType will correctly return a $ref to it
		// because it will hit the `visited[originalTypeFullName]` check (if direct recursion) or the `t.Name() == currentlyDefiningCompName` logic combined with `isComponent`
		// (if that component's definition process has started).
		itemSchema, err := b.generateSchemaFromGoType(elemType, visited, currentlyDefiningCompName)
		if err != nil {
			return nil, fmt.Errorf("failed to generate item schema for array/slice element type %s: %w", elemType.String(), err)
		}

		if itemSchema.Ref != "" {
			schema.Items = &entity.SchemaRef{Ref: itemSchema.Ref}
		} else {
			schema.Items = &entity.SchemaRef{Schema: itemSchema}
		}

	case reflect.Struct:
		// Note: time.Time handled above. Named structs also partially handled above for recursion/component checks.
		schema.Type = "object"
		schema.Properties = make(map[string]*entity.SchemaRef)
		var requiredFields []string

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}

			jsonTag := field.Tag.Get("json")
			fieldName := field.Name
			omitempty := false

			if jsonTag != "" {
				parts := strings.Split(jsonTag, ",")
				if parts[0] == "-" {
					continue
				}
				if parts[0] != "" {
					fieldName = parts[0]
				}
				for _, part := range parts[1:] {
					if part == "omitempty" {
						omitempty = true
						break
					}
				}
			}

			propSchema, err := b.generateSchemaFromGoType(field.Type, visited, currentlyDefiningCompName)
			if err != nil {
				return nil, fmt.Errorf("failed to generate schema for field '%s' in struct '%s': %w", field.Name, t.Name(), err)
			}

			if propSchema.Ref != "" {
				schema.Properties[fieldName] = &entity.SchemaRef{Ref: propSchema.Ref}
			} else {
				schema.Properties[fieldName] = &entity.SchemaRef{Schema: propSchema}
			}

			if !omitempty {
				requiredFields = append(requiredFields, fieldName)
			}
		}

		if len(requiredFields) > 0 {
			schema.Required = requiredFields
		}

	case reflect.Map:
		schema.Type = "object"
		if t.Key().Kind() != reflect.String {
			// Consider logging a warning, as OpenAPI map keys must be strings.
		}

		valType := t.Elem()
		addPropsSchema, err := b.generateSchemaFromGoType(valType, visited, currentlyDefiningCompName)
		if err != nil {
			return nil, fmt.Errorf("failed to generate schema for map value type %s: %w", valType.Name(), err)
		}

		if addPropsSchema.Ref != "" {
			schema.AdditionalProperties = &entity.SchemaRef{Ref: addPropsSchema.Ref}
		} else {
			schema.AdditionalProperties = &entity.SchemaRef{Schema: addPropsSchema}
		}

	default:
		return nil, fmt.Errorf("unsupported type for DTO schema generation: %s (Kind: %s)", t.String(), t.Kind())
	}

	return schema, nil
}

func (b *SwaggerDocBuilder) ExternalDocumentation(url string, description string) openapi.SwaggerDocBuilder {
	b.doc.ExternalDocs = &entity.ExternalDocumentation{URL: url, Description: description}
	return b
}

func (b *SwaggerDocBuilder) GlobalSecurity(requirement entity.SecurityRequirement) openapi.SwaggerDocBuilder {
	b.doc.Security = append(b.doc.Security, requirement)
	return b
}

func (b *SwaggerDocBuilder) Build() entity.SwaggerDocEntity {
	return *b.doc
}

// Component related methods (stubs or full implementations)
func (b *SwaggerDocBuilder) ComponentResponse(name string, config func(openapi.Response)) openapi.SwaggerDocBuilder {
	if b.doc.Components.Responses == nil {
		b.doc.Components.Responses = make(map[string]*entity.ResponseRef)
	}
	resp := entity.Response{}
	builder := &ResponseBuilder{response: &resp, docBuilder: b}
	config(builder)
	b.doc.Components.Responses[name] = &entity.ResponseRef{Response: &resp}
	return b
}

func (b *SwaggerDocBuilder) ComponentParameter(name string, config func(openapi.Parameter)) openapi.SwaggerDocBuilder {
	if b.doc.Components.Parameters == nil {
		b.doc.Components.Parameters = make(map[string]*entity.ParameterRef)
	}
	param := entity.Parameter{}
	builder := &ParameterBuilder{param: &param, docBuilder: b}
	config(builder)
	b.doc.Components.Parameters[name] = &entity.ParameterRef{Parameter: &param}
	return b
}

func (b *SwaggerDocBuilder) ComponentExample(name string, config func(openapi.Example)) openapi.SwaggerDocBuilder {
	if b.doc.Components.Examples == nil {
		b.doc.Components.Examples = make(map[string]*entity.ExampleRef)
	}
	ex := entity.Example{}
	builder := &ExampleBuilder{example: &ex}
	config(builder)
	b.doc.Components.Examples[name] = &entity.ExampleRef{Example: &ex}
	return b
}

func (b *SwaggerDocBuilder) ComponentRequestBody(name string, config func(openapi.RequestBody)) openapi.SwaggerDocBuilder {
	if b.doc.Components.RequestBodies == nil {
		b.doc.Components.RequestBodies = make(map[string]*entity.RequestBodyRef)
	}
	reqBody := entity.RequestBody{}
	builder := &RequestBodyBuilder{requestBody: &reqBody, docBuilder: b}
	config(builder)
	b.doc.Components.RequestBodies[name] = &entity.RequestBodyRef{RequestBody: &reqBody}
	return b
}

func (b *SwaggerDocBuilder) ComponentHeader(name string, config func(openapi.Header)) openapi.SwaggerDocBuilder {
	if b.doc.Components.Headers == nil {
		b.doc.Components.Headers = make(map[string]*entity.HeaderRef)
	}
	hdr := entity.Header{}
	builder := &HeaderBuilder{header: &hdr, docBuilder: b}
	config(builder)
	b.doc.Components.Headers[name] = &entity.HeaderRef{Header: &hdr}
	return b
}

func (b *SwaggerDocBuilder) ComponentSecurityScheme(name string, config func(openapi.SecurityScheme)) openapi.SwaggerDocBuilder {
	if b.doc.Components.SecuritySchemes == nil {
		b.doc.Components.SecuritySchemes = make(map[string]*entity.SecuritySchemeRef)
	}
	secScheme := entity.SecurityScheme{}
	builder := &SecuritySchemeBuilder{scheme: &secScheme}
	config(builder)
	b.doc.Components.SecuritySchemes[name] = &entity.SecuritySchemeRef{SecurityScheme: &secScheme}
	return b
}

func (b *SwaggerDocBuilder) ComponentLink(name string, config func(openapi.Link)) openapi.SwaggerDocBuilder {
	if b.doc.Components.Links == nil {
		b.doc.Components.Links = make(map[string]*entity.LinkRef)
	}
	link := entity.Link{}
	builder := &LinkBuilder{link: &link}
	config(builder)
	b.doc.Components.Links[name] = &entity.LinkRef{Link: &link}
	return b
}

func (b *SwaggerDocBuilder) ComponentCallback(name string, config func(openapi.Callback)) openapi.SwaggerDocBuilder {
	if b.doc.Components.Callbacks == nil {
		b.doc.Components.Callbacks = make(map[string]*entity.CallbackRef)
	}
	// Proper CallbackBuilder would be needed for full implementation.
	// This is a placeholder based on current structure.
	// Callbacks are map[string]PathItem or map[string]$ref. Example structure:
	// cbVal := entity.PathItem{}
	// cbRef := &entity.CallbackRef{PathItem: &cbVal} // if CallbackRef wraps PathItem
	// builder := &CallbackBuilder{...} // Requires CallbackBuilder and Callback interface logic
	// config(builder)
	// b.doc.Components.Callbacks[name] = cbRef
	return b
}
