package kallax

import "fmt"

// Schema represents a table schema in the database. Contains some information
// like the table name, its columns, its identifier and so on.
type Schema interface {
	// Alias returns the name of the alias used in queries for this schema.
	Alias() string
	// Table returns the table name.
	Table() string
	// ID returns the name of the identifier of the table.
	ID() SchemaField
	// Columns returns the list of columns in the schema.
	Columns() []SchemaField
	// ForeignKey returns the name of the foreign key of the given model field.
	ForeignKey(string) (SchemaField, bool)
	// WithAlias returns a new schema with the given string added to the
	// default alias.
	// Calling WithAlias on a schema returned by WithAlias not return a
	// schema based on the child, but another based on the parent.
	WithAlias(string) Schema
}

// BaseSchema is the basic implementation of Schema.
type BaseSchema struct {
	alias       string
	table       string
	foreignKeys ForeignKeys
	id          SchemaField
	columns     []SchemaField
}

// NewBaseSchema creates a new schema with the given table, alias, identifier
// and columns.
func NewBaseSchema(table, alias string, id SchemaField, fks ForeignKeys, columns ...SchemaField) *BaseSchema {
	return &BaseSchema{
		alias:       alias,
		table:       table,
		foreignKeys: fks,
		id:          id,
		columns:     columns,
	}
}

func (s *BaseSchema) Alias() string          { return s.alias }
func (s *BaseSchema) Table() string          { return s.table }
func (s *BaseSchema) ID() SchemaField        { return s.id }
func (s *BaseSchema) Columns() []SchemaField { return s.columns }
func (s *BaseSchema) ForeignKey(field string) (SchemaField, bool) {
	k, ok := s.foreignKeys[field]
	return k, ok
}
func (s *BaseSchema) WithAlias(field string) Schema {
	return &aliasSchema{s, field}
}

type aliasSchema struct {
	*BaseSchema
	alias string
}

func (s *aliasSchema) Alias() string {
	return fmt.Sprintf("%s_%s", s.BaseSchema.Alias(), s.alias)
}

// ForeignKeys is a mapping between relationships and their foreign key field.
type ForeignKeys map[string]SchemaField

// SchemaField is a named field in the table schema.
type SchemaField interface {
	isSchemaField()
	// String returns the string representation of the field. That is, its name.
	String() string
	// QualifiedString returns the name of the field qualified by the alias of
	// the given schema.
	QualifiedName(Schema) string
}

// BaseSchemaField is a basic schema field with name.
type BaseSchemaField struct {
	name string
}

// NewSchemaField creates a new schema field with the given name.
func NewSchemaField(name string) SchemaField {
	return &BaseSchemaField{name}
}

func (*BaseSchemaField) isSchemaField() {}

func (f BaseSchemaField) String() string {
	return f.name
}

func (f *BaseSchemaField) QualifiedName(schema Schema) string {
	alias := schema.Alias()
	if alias != "" {
		return fmt.Sprintf("%s.%s", alias, f.name)
	}
	return f.name
}

// Relationship is a relationship with its schema and the field of te relation
// in the record.
type Relationship struct {
	// Field is the field in the record where the relationship is.
	Field string
	// Schema is the schema of the relationship.
	Schema Schema
}

// ColumnNames returns the names of the given schema fields.
func ColumnNames(columns []SchemaField) []string {
	var names = make([]string, len(columns))
	for i, v := range columns {
		names[i] = v.String()
	}
	return names
}