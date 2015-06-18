// Author: kris@keypr.com
// Notes: Primary use case for this is to transform JSON schemas into something that can
// be represented as a go struct. Therefore, for all object properties one and
// only one type can be defined.
// [string] == OK
// [string, number] != OK
//
package gojsonschema

import (
	"fmt"
)

type SchemaDescription struct {
	Id          string
	Name        string
	Description string
	Properties  []Property
	schema      *subSchema
}

type Property struct {
	Name        string
	Description string
	Type        string
	Tags        map[string]string
	Required    bool
	reference   *subSchema
}

type SchemaDescriber interface {
	EmbeddedSchemaNames() []string
}

type PropertyDescriber interface {
	IsReference() bool
	RefSchemaName() string
}

type Inspector interface {
	GetObjectDescription() SchemaDescription
	GetResolvedProperties() []Property
	ValidateStructIntegrity() error
}

// Schema Describer Interface
// See if this schema references any embedded schemas; may not be needed.
func (s SchemaDescription) HasEmbeddedSchema() bool {
	return len(s.EmbeddedSchemaNames()) > 0
}

// Get the names of any embedded schemas. This wraps a private variadic function
func (s SchemaDescription) EmbeddedSchemaNames() []string {
	return extractSchemaRefs(s.schema.allOf, s.schema.anyOf, s.schema.oneOf)
}

// Private (variadic) function to actually pull the schema names
func extractSchemaRefs(s ...[]*subSchema) []string {
	collector := make([]string, 0)

	for _, schemaCol := range s {
		if len(schemaCol) > 1 { // if there's nothing in this collection skip it.
			for _, item := range schemaCol {
				if item.refSchema != nil {
					collector = append(collector, *item.refSchema.title)
				}
			}
		}
	}
	return collector
}

// Inspector interface
// Resolve this document's property definition
func (s *Schema) GetResolvedProperties() []Property {
	var ref *subSchema

	rtn := make([]Property, 0)

	// Check for Root properties
	for _, prop := range s.rootSchema.propertiesChildren {
		if prop.refSchema != nil {
			ref = prop.refSchema
		} else {
			ref = nil
		}

		rtn = append(rtn, Property{
			Name:        prop.property,
			Description: pointerToString(prop.description),
			Type:        prop.types.String(),
			Required:    s.IsRequiredProperty(prop.property),
			reference:   ref,
		})
	}

	// Collect properties that are defined in combiners
	rtn = append(rtn, extractSchemaProperties(s, s.rootSchema.allOf, s.rootSchema.anyOf, s.rootSchema.oneOf)...)
	return rtn
}

func extractSchemaProperties(parent *Schema, schemas ...[]*subSchema) []Property {
	var ref *subSchema
	rtn := make([]Property, 0)

	for _, schemaCol := range schemas {
		if len(schemaCol) > 1 {
			for _, item := range schemaCol {
				for _, prop := range item.propertiesChildren {
					if prop.refSchema != nil { // Account for references in the property
						ref = item.refSchema
					} else {
						ref = nil
					}

					rtn = append(rtn, Property{
						Name:        prop.property,
						Description: pointerToString(prop.description),
						Type:        prop.types.String(),
						Required:    parent.IsRequiredProperty(prop.property),
						reference:   ref,
					})
				}
			}
		}
	}
	return rtn
}

// Get Object Description
func (s *Schema) GetObjectDescription() SchemaDescription {
	props := s.GetResolvedProperties()
	// fmt.Println(*s.rootSchema.id)
	// fmt.Printf("%+v\n", s.rootSchema)

	return SchemaDescription{
		Id:          *s.rootSchema.id,
		Name:        *s.rootSchema.title,
		Description: *s.rootSchema.description,
		Properties:  props,
		schema:      s.rootSchema,
	}
}

// Pass in a property name and see if it is required in the schema
func (s *Schema) IsRequiredProperty(propName string) bool {
	rtn := false

	for _, required := range s.rootSchema.required {
		rtn = rtn || (required == propName)
	}

	return rtn
}

// PropertyDescriber interface

func (p *Property) IsReference() bool {
	return p.reference != nil
}

func (p *Property) RefSchemaName() string {
	if p.reference == nil {
		return ""
	}

	// This should already be the refSchema getting passed in here.
	return *p.reference.title
}

// Needed? needs to be finished if so...
func isLeafSchema(s *subSchema) bool {
	combinators := (len(s.oneOf) == 0 && len(s.anyOf) == 0 && len(s.allOf) == 0)
	return combinators
}

// utility function to parse missing / nil string pointers
func pointerToString(ptr *string) string {
	if ptr != nil {
		return *ptr
	} else {
		return ""
	}
}

func (s *Schema) Probe() {
	for _, x := range s.rootSchema.allOf {
		if x.refSchema != nil {
			for _, y := range (*x.refSchema).propertiesChildren {
				fmt.Println(">>>CHILD<<<")
				fmt.Printf("%+v\n", y)
			}
		}
		fmt.Printf("%+v\n", *x)
	}

}
