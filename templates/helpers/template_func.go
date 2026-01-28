package helpers

import (
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

func TemplateFunc(g *protogen.GeneratedFile) template.FuncMap {
	return template.FuncMap{
		"initialLower": func(s string) string {
			if s == "" {
				return s
			}
			r := []rune(s)
			r[0] = []rune(strings.ToLower(string(r[0])))[0]
			return string(r)
		},
		"trimSuffix": strings.TrimSuffix,
		"lower":      strings.ToLower,
		"join":       func(parts ...string) string { return strings.Join(parts, "/") },
		// Helper function to check if a field exists in a struct
		"hasField": func(fieldName string) bool {
			// For common fields in entities, return true
			return fieldName == "Name" || fieldName == "Email" || fieldName == "Id" ||
				fieldName == "CreatedAt" || fieldName == "UpdatedAt" || fieldName == "Description"
		},
		// Qualify a GoIdent with the correct import alias and register import in this file
		"qual": func(id protogen.GoIdent) string {
			return g.QualifiedGoIdent(id)
		},
		// Return "context.Context" with import registered
		"ctx": func() string {
			return g.QualifiedGoIdent(protogen.GoIdent{
				GoName:       "Context",
				GoImportPath: "context",
			})
		},
	}
}
