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
	Name        string
	Description string
	Properties  []Property
}

type Property struct {
	Name        string
	Description string
	Type        string
	Required    bool
}

type Inspector interface {
	GetObjectDescription() SchemaDescription
	GetResolvedProperties() []Property
	ValidateStructIntegrity() error
}

// Resolve this document's property definition
func (s *Schema) GetResolvedProperties() []Property {
	rtn := make([]Property, 0)

	for _, prop := range s.rootSchema.propertiesChildren {
		rtn = append(rtn, Property{prop.property,
			pointerToString(prop.description),
			prop.types.String(),
			s.IsRequiredProperty(prop.property)})
	}
	return rtn
}

// Get Object Description
func (s *Schema) GetObjectDescription() SchemaDescription {
	return SchemaDescription{Name: *s.rootSchema.title, Description: *s.rootSchema.description}
}

// Pass in a property name and see if it is required in the schema
func (s *Schema) IsRequiredProperty(propName string) bool {
	rtn := false

	for _, required := range s.rootSchema.required {
		rtn = rtn || (required == propName)
	}

	return rtn
}

func getRecursiveProperties(s *subSchema, props *[]Property) {
	// Base case
	if s == nil {
		return
	}

	// If type is an object, iterate through child properties
	if s.types.Contains(TYPE_OBJECT) {
	}

	// Check for items?
	if s.types.Contains(TYPE_ARRAY) {
	}

	// if this is a schema reference, dereference this.
	if s.refSchema != nil {
		getRecursiveProperties(s.refSchema, props)
	}

	// Definitions should be independent structs
	if len(s.definitions) > 0 {

	}
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
				fmt.Println(y)
			}
		}
		fmt.Println(*x)
	}

}
