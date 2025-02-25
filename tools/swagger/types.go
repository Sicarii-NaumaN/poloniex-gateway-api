package swagger

type swaggerParamType string

const (
	query  swaggerParamType = "query"
	path   swaggerParamType = "path"
	header swaggerParamType = "header"
	body   swaggerParamType = "body"
)

type swaggerFieldType string

const (
	// Object type
	Object swaggerFieldType = "object"
	// Array type
	Array swaggerFieldType = "array"
	// String type
	String swaggerFieldType = "string"
	// Int64 type
	Int64 swaggerFieldType = "long"
	// Float64 type
	Float64 swaggerFieldType = "number"
	// Integer type
	Integer swaggerFieldType = "integer"
	// Map type
	Map swaggerFieldType = "map"
	// Boolean type
	Boolean swaggerFieldType = "boolean"
	// File type
	File swaggerFieldType = "file"
)
