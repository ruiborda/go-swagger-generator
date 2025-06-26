package swagger

import (
	"fmt"
	openapi "github.com/ruiborda/go-swagger-generator/v2/src/openapi"
	entity "github.com/ruiborda/go-swagger-generator/v2/src/openapi_spec"
	"reflect"
	"regexp"
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

// generateComponentSchemaName creates a unique schema name for a type by sanitizing its full string representation.
// This is crucial for handling generic types correctly.
func (b *SwaggerDocBuilder) generateComponentSchemaName(typ reflect.Type) string {
	nameStr := typ.String() // e.g., "main.Response[main.UserData]", "[]*main.UserData", "main.NonGenericStruct"

	// Replace package paths (e.g., "main.Type" -> "main_Type", "pkg.Type" -> "pkg_Type")
	s := strings.ReplaceAll(nameStr, ".", "_")

	// Handle pointers and slices in a consistent way for type arguments or nested generics
	s = strings.ReplaceAll(s, "*", "Ptr")     // *pkg_Type -> PtrPkg_Type
	s = strings.ReplaceAll(s, "[]", "ListOf") // []pkg_Type -> ListOfPkg_Type (consistent with top-level slice naming)

	// Handle generic type bracketing and multiple type arguments
	s = strings.ReplaceAll(s, "[", "_") // pkg_Response[pkg_Data] -> pkg_Response_pkg_Data]
	s = strings.ReplaceAll(s, "]", "")  // pkg_Response_pkg_Data] -> pkg_Response_pkg_Data
	s = strings.ReplaceAll(s, ",", "_") // For types like T[A, B]

	// Clean up any resulting multiple underscores and leading/trailing underscores
	reg := regexp.MustCompile(`_+`)
	s = reg.ReplaceAllString(s, "_")
	s = strings.Trim(s, "_")

	// If the original type was a simple non-generic struct like "main.MyData",
	// it would become "main_MyData". If we prefer "MyData", we might strip common package prefixes.
	// For now, this comprehensive naming ensures uniqueness.
	return s
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
	return b.schemaFromDTORecursive(dtoInstance, make(map[string]string))
}

// schemaFromDTORecursive is the internal implementation for SchemaFromDTO.
// It assumes definitionsMux is already locked.
// processedInThisCall tracks types (reflect.Type.String()) to their component names (string)
// that have already been initiated or completed within the *current* top-level SchemaFromDTO call chain.
func (b *SwaggerDocBuilder) schemaFromDTORecursive(dtoInstance interface{}, processedInThisCall map[string]string) (string, error) {
	if dtoInstance == nil {
		return "", fmt.Errorf("SchemaFromDTO called with nil DTO instance")
	}

	originalDtoType := reflect.TypeOf(dtoInstance)
	currentDtoType := originalDtoType
	originalDtoTypeString := originalDtoType.String() // Use .String() for unique key with generics

	if currentDtoType.Kind() == reflect.Ptr {
		if reflect.ValueOf(dtoInstance).IsNil() {
			dtoInstance = reflect.New(currentDtoType.Elem()).Interface()
		}
		currentDtoType = currentDtoType.Elem()
	}

	if componentName, ok := processedInThisCall[originalDtoTypeString]; ok {
		return componentName, nil
	}

	if currentDtoType.Kind() == reflect.Slice || currentDtoType.Kind() == reflect.Array {
		elementType := currentDtoType.Elem()
		baseElementType := elementType
		if baseElementType.Kind() == reflect.Ptr {
			baseElementType = baseElementType.Elem()
		}

		if baseElementType.Kind() != reflect.Struct {
			// Allow slices of basic types, they don't become components but are valid schemas.
			// For SchemaFromDTO, the top level should ideally be a struct or slice of structs.
			// If it's a slice of primitives, it won't be a named component.
			// This function is for creating named components primarily.
			// For a slice of primitives, a schema can be generated but not registered as a named component via this path.
			return "", fmt.Errorf("slice/array element DTO for component registration must be a struct or pointer to struct, got %s for element of %s", baseElementType.Kind(), currentDtoType.String())
		}

		elementInstance := reflect.New(baseElementType).Interface()
		elementComponentName, err := b.schemaFromDTORecursive(elementInstance, processedInThisCall)
		if err != nil {
			return "", fmt.Errorf("failed to generate schema for slice/array element type %s: %w", baseElementType.Name(), err)
		}

		// Use generateComponentSchemaName for consistency, although ListOf<Name> is also common.
		// listDtoName := "ListOf" + elementComponentName
		listDtoName := b.generateComponentSchemaName(currentDtoType) // e.g. ListOf_pkg_MyData

		if _, exists := b.doc.Components.Schemas[listDtoName]; exists {
			processedInThisCall[originalDtoTypeString] = listDtoName
			return listDtoName, nil
		}

		arraySchema := &entity.Schema{
			Type: "array",
			Items: &entity.SchemaRef{
				Ref: "#/components/schemas/" + elementComponentName,
			},
		}
		b.doc.Components.Schemas[listDtoName] = &entity.SchemaRef{Schema: arraySchema}
		processedInThisCall[originalDtoTypeString] = listDtoName
		return listDtoName, nil
	}

	if currentDtoType.Kind() == reflect.Struct {
		structDtoName := b.generateComponentSchemaName(currentDtoType)
		if currentDtoType.Name() == "" && !strings.Contains(currentDtoType.String(), "[") { // Check if it's truly anonymous and not a generic struct
			// Anonymous, non-generic structs cannot be top-level components.
			// Generic structs like somepkg.Gen[T] will have Name() but are fine.
			return "", fmt.Errorf("anonymous structs cannot be registered as top-level DTOs via SchemaFromDTO: %s", currentDtoType.String())
		}

		if _, exists := b.doc.Components.Schemas[structDtoName]; exists {
			processedInThisCall[originalDtoTypeString] = structDtoName
			return structDtoName, nil
		}

		b.doc.Components.Schemas[structDtoName] = &entity.SchemaRef{Ref: "#/components/schemas/" + structDtoName} // Placeholder
		processedInThisCall[originalDtoTypeString] = structDtoName

		generatedSchema, err := b.generateSchemaFromGoType(currentDtoType, make(map[string]bool), structDtoName)
		if err != nil {
			delete(b.doc.Components.Schemas, structDtoName)
			delete(processedInThisCall, originalDtoTypeString)
			return "", fmt.Errorf("failed to generate schema for DTO struct %s: %w", structDtoName, err)
		}

		if generatedSchema.Ref != "" {
			delete(b.doc.Components.Schemas, structDtoName)
			delete(processedInThisCall, originalDtoTypeString)
			return "", fmt.Errorf("internal error: generateSchemaFromGoType for component %s returned a $ref ('%s'), expected full schema definition", structDtoName, generatedSchema.Ref)
		}

		b.doc.Components.Schemas[structDtoName] = &entity.SchemaRef{Schema: generatedSchema}
		return structDtoName, nil
	}

	return "", fmt.Errorf("DTO must be a struct, pointer to struct, slice/array of structs, or pointer to slice/array of structs, got %s", originalDtoType.String())
}

// generateSchemaFromGoType converts a Go type to an OpenAPI Schema object.
// visited map (keyed by type.String()) is used to handle recursive types for the *current* struct's fields.
// currentlyDefiningCompName is the name of the component schema that SchemaFromDTO is currently trying to define.
func (b *SwaggerDocBuilder) generateSchemaFromGoType(t reflect.Type, visited map[string]bool, currentlyDefiningCompName string) (*entity.Schema, error) {
	originalTypeString := t.String() // Use .String() for unique key with generics

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Name() != "" && t.Kind() == reflect.Struct { // Named structs (generic or non-generic)
		potentialComponentName := b.generateComponentSchemaName(t)

		if visited[originalTypeString] {
			// Cycle detected for the current struct's field expansion. This type is already being processed.
			// If this is the component being defined, this recursive call should refer to it.
			// If it's another component, it should also refer to its global name.
			return &entity.Schema{Ref: "#/components/schemas/" + potentialComponentName}, nil
		}

		if potentialComponentName != currentlyDefiningCompName {
			if _, isComponent := b.doc.Components.Schemas[potentialComponentName]; isComponent {
				// This refers to a different, already known (or placeholder for) component.
				return &entity.Schema{Ref: "#/components/schemas/" + potentialComponentName}, nil
			}
		}
		// If potentialComponentName == currentlyDefiningCompName, we must expand it (not ref).
		// If it's not currentlyDefiningCompName AND not yet a component, expand it inline (unless it becomes one via SchemaFromDTO during field processing).

		visited[originalTypeString] = true
		defer delete(visited, originalTypeString)
	}

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
		schema.Format = "int32"
		minValue := 0.0
		schema.Minimum = &minValue
	case reflect.Uint64:
		schema.Type = "integer"
		schema.Format = "int64"
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
		itemSchema, err := b.generateSchemaFromGoType(elemType, visited, currentlyDefiningCompName)
		if err != nil {
			return nil, fmt.Errorf("failed to generate item schema for array/slice element type %s: %w", elemType.String(), err)
		}

		if itemSchema.Ref != "" {
			schema.Items = &entity.SchemaRef{Ref: itemSchema.Ref}
		} else {
			// If itemSchema is complex but not a ref, ensure it's fully defined.
			// If elemType is a struct that wasn't registered as a component (e.g. anonymous or not top-level DTO),
			// its schema will be inline here.
			schema.Items = &entity.SchemaRef{Schema: itemSchema}
		}

	case reflect.Struct:
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
			// OpenAPI map keys must be strings. Consider logging a warning or error.
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
	return b
}
