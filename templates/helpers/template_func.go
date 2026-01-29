package helpers

import (
	"crypto/rand"
	"fmt"
	"strings"
	"text/template"
	"unicode"
)

// AddTemplateFuncs adds custom functions to the template
func AddTemplateFuncs(tmpl *template.Template) *template.Template {
	return tmpl.Funcs(template.FuncMap{
		"lower":        strings.ToLower,
		"initialLower": initialLower,
		"upper":        strings.ToUpper,
		"title":        strings.Title,
		"snake":        toSnakeCase,
		"camel":        toCamelCase,
		"lowerCamel":   toLowerCamelCase,
		"kebab":        toKebabCase,
		"genUUID":      generateUUID,
		"shortUUID":    generateShortUUID,
		"noSuffix":     removeSuffix,
		"serviceNoSuf": serviceNoSuffix,
	})
}

// initialLower converts the first character of a string to lowercase
func initialLower(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// generateUUID generates a full UUID string
func generateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "00000000-0000-0000-0000-000000000000"
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// generateShortUUID generates a short UUID-like string (8 chars)
func generateShortUUID() string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "00000000"
	}
	return fmt.Sprintf("%02x", b)
}

// removeSuffix removes "Service" suffix from a service name
func removeSuffix(s string) string {
	return strings.TrimSuffix(s, "Service")
}

// serviceNoSuffix removes "Service" suffix and returns properly cased version
func serviceNoSuffix(s string) string {
	// Remove "Service" suffix
	name := strings.TrimSuffix(s, "Service")
	// Ensure first letter is uppercase for consistent casing
	if len(name) > 0 {
		name = strings.ToUpper(name[:1]) + name[1:]
	}
	return name
}

// toSnakeCase converts a string to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// toCamelCase converts a string to CamelCase
func toCamelCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || unicode.IsSpace(r)
	})

	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}

	return strings.Join(parts, "")
}

// toLowerCamelCase converts a string to lowerCamelCase
func toLowerCamelCase(s string) string {
	result := toCamelCase(s)
	if result == "" {
		return ""
	}
	return initialLower(result)
}

// toKebabCase converts a string to kebab-case
func toKebabCase(s string) string {
	return strings.ReplaceAll(toSnakeCase(s), "_", "-")
}
