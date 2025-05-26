package swagger

import (
	openapi "github.com/ruiborda/go-swagger-generator/src/openapi"
	entity "github.com/ruiborda/go-swagger-generator/src/openapi_spec"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
)

type RequestBodyBuilder struct {
	requestBody *entity.RequestBody
	docBuilder  *SwaggerDocBuilder
}

func (b *RequestBodyBuilder) Description(description string) openapi.RequestBody {
	b.requestBody.Description = description
	return b
}

func (b *RequestBodyBuilder) Content(mimeTypeStr mime.MimeType, config func(openapi.MediaType)) openapi.RequestBody {
	if b.requestBody.Content == nil {
		b.requestBody.Content = make(map[string]*entity.MediaType)
	}
	mediaType := &entity.MediaType{}
	mediaTypeBuilder := &MediaTypeBuilder{mediaType: mediaType, docBuilder: b.docBuilder}
	config(mediaTypeBuilder)
	b.requestBody.Content[string(mimeTypeStr)] = mediaType
	return b
}

func (b *RequestBodyBuilder) Required(required bool) openapi.RequestBody {
	b.requestBody.Required = required
	return b
}

func (b *RequestBodyBuilder) Ref(ref string) openapi.RequestBody {
	// This would imply the RequestBody itself is a reference.
	// Handled by OperationBuilder when setting RequestBodyRef.
	return b
}

type MediaTypeBuilder struct {
	mediaType  *entity.MediaType
	docBuilder *SwaggerDocBuilder
}

func (b *MediaTypeBuilder) Schema(config func(openapi.Schema)) openapi.MediaType {
	if b.mediaType.Schema == nil {
		b.mediaType.Schema = &entity.SchemaRef{Schema: &entity.Schema{}}
	}
	if b.mediaType.Schema.Schema == nil { // If it was a Ref, ensure Schema part is initialized
		b.mediaType.Schema.Schema = &entity.Schema{}
	}
	schemaBuilder := &SchemaBuilder{schema: b.mediaType.Schema.Schema, docBuilder: b.docBuilder}
	config(schemaBuilder)
	return b
}

func (b *MediaTypeBuilder) SchemaFromDTO(dto interface{}) openapi.MediaType {
	dtoName, err := b.docBuilder.SchemaFromDTO(dto)
	if err != nil {
		return b
	}
	b.mediaType.Schema = &entity.SchemaRef{Ref: "#/components/schemas/" + dtoName}
	return b
}

func (b *MediaTypeBuilder) SchemaRef(ref string) openapi.MediaType {
	b.mediaType.Schema = &entity.SchemaRef{Ref: ref}
	return b
}

func (b *MediaTypeBuilder) Example(value interface{}) openapi.MediaType {
	b.mediaType.Example = value
	return b
}

func (b *MediaTypeBuilder) Examples(name string, config func(openapi.Example)) openapi.MediaType {
	if b.mediaType.Examples == nil {
		b.mediaType.Examples = make(map[string]*entity.ExampleRef)
	}
	ex := entity.Example{}
	exBuilder := &ExampleBuilder{example: &ex}
	config(exBuilder)
	b.mediaType.Examples[name] = &entity.ExampleRef{Example: &ex}
	return b
}

func (b *MediaTypeBuilder) Encoding(propertyName string, config func(openapi.Encoding)) openapi.MediaType {
	if b.mediaType.Encoding == nil {
		b.mediaType.Encoding = make(map[string]*entity.Encoding)
	}
	enc := entity.Encoding{}
	encBuilder := &EncodingBuilder{encoding: &enc}
	config(encBuilder)
	b.mediaType.Encoding[propertyName] = &enc
	return b
}

type ExampleBuilder struct {
	example *entity.Example
}

func (b *ExampleBuilder) Summary(summary string) openapi.Example {
	b.example.Summary = summary
	return b
}
func (b *ExampleBuilder) Description(description string) openapi.Example {
	b.example.Description = description
	return b
}
func (b *ExampleBuilder) Value(value interface{}) openapi.Example {
	b.example.Value = value
	return b
}
func (b *ExampleBuilder) ExternalValue(url string) openapi.Example {
	b.example.ExternalValue = url
	return b
}
func (b *ExampleBuilder) Ref(ref string) openapi.Example {
	// This example is now a ref.
	// Logic to handle this if entity.Example also had a $ref field or if ExampleRef used.
	return b
}

type EncodingBuilder struct {
	encoding *entity.Encoding
}

func (b *EncodingBuilder) ContentType(contentType string) openapi.Encoding {
	b.encoding.ContentType = contentType
	return b
}
func (b *EncodingBuilder) Header(name string, config func(openapi.Header)) openapi.Encoding {
	if b.encoding.Headers == nil {
		b.encoding.Headers = make(map[string]*entity.HeaderRef)
	}
	hdr := entity.Header{}
	hdrBuilder := &HeaderBuilder{header: &hdr} // Assuming HeaderBuilder exists and is compatible
	config(hdrBuilder)
	b.encoding.Headers[name] = &entity.HeaderRef{Header: &hdr}
	return b
}
func (b *EncodingBuilder) Style(style string) openapi.Encoding {
	b.encoding.Style = style
	return b
}
func (b *EncodingBuilder) Explode(explode bool) openapi.Encoding {
	b.encoding.Explode = explode
	return b
}
func (b *EncodingBuilder) AllowReserved(allowReserved bool) openapi.Encoding {
	b.encoding.AllowReserved = allowReserved
	return b
}

// LinkBuilder implementation (if not in its own file)
type LinkBuilder struct {
	link *entity.Link
}

func (b *LinkBuilder) OperationRef(ref string) openapi.Link {
	b.link.OperationRef = ref
	return b
}
func (b *LinkBuilder) OperationID(id string) openapi.Link {
	b.link.OperationId = id
	return b
}
func (b *LinkBuilder) Parameter(name string, expressionOrValue interface{}) openapi.Link {
	if b.link.Parameters == nil {
		b.link.Parameters = make(map[string]interface{})
	}
	b.link.Parameters[name] = expressionOrValue
	return b
}
func (b *LinkBuilder) RequestBody(expressionOrValue interface{}) openapi.Link {
	b.link.RequestBody = expressionOrValue
	return b
}
func (b *LinkBuilder) Description(description string) openapi.Link {
	b.link.Description = description
	return b
}
func (b *LinkBuilder) Server(url string, config func(openapi.Server)) openapi.Link {
	srv := entity.Server{URL: url}
	b.link.Server = &srv
	return b
}
func (b *LinkBuilder) Ref(ref string) openapi.Link {
	return b
}
