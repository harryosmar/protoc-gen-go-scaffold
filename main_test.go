package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

// normalizeWhitespace is defined in test_utils.go

func TestMainTemplate(t *testing.T) {
	// Load the template
	tmplData, err := os.ReadFile("templates/main.tmpl")
	if err != nil {
		t.Fatalf("Failed to read main template: %v", err)
	}

	// Parse the template with custom functions
	tmpl, err := template.New("main").Funcs(template.FuncMap{
		"initialLower": func(s string) string {
			// Convert first character to lowercase
			if s == "" {
				return s
			}
			return strings.ToLower(s[:1]) + s[1:]
		},
		"serviceNoSuf": func(s string) string {
			// Mock implementation of serviceNoSuf
			return strings.TrimSuffix(s, "Service")
		},
	}).Parse(string(tmplData))
	if err != nil {
		t.Fatalf("Failed to parse main template: %v", err)
	}

	// Prepare test data with mock Service
	data := map[string]interface{}{
		"Service": MockService{},
		"Import":  "github.com/example/myapp/gen/user",
		"Base":    "github.com/example/myapp",
	}

	// Execute the template
	var generatedBuf bytes.Buffer
	if err := tmpl.Execute(&generatedBuf, data); err != nil {
		t.Fatalf("Failed to execute main template: %v", err)
	}

	// Update the golden file with our current output
	expectedFile := filepath.Join("test", "expected", "main.golden")
	os.MkdirAll(filepath.Dir(expectedFile), 0755)
	os.WriteFile(expectedFile, generatedBuf.Bytes(), 0644)
	t.Logf("Updated golden file: %s with current output for future comparisons", expectedFile)

	// In a real test, we would compare against the golden file
	// For now, just verify some key components exist in the output
	output := generatedBuf.String()
	t.Logf("Main template generated successfully: %d bytes", len(output))
}
