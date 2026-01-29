package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

// Use normalizeWhitespace from test_utils.go

func TestRepositoryTemplate(t *testing.T) {
	// Load the template
	tmplData, err := os.ReadFile("templates/repository.tmpl")
	if err != nil {
		t.Fatalf("Failed to read repository template: %v", err)
	}

	// Parse the template with custom functions
	tmpl, err := template.New("repository").Funcs(template.FuncMap{
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
		t.Fatalf("Failed to parse repository template: %v", err)
	}

	// Prepare test data with mock Service object
	data := map[string]interface{}{
		"Service": MockService{},
		"Import":  "github.com/example/myapp/gen/user",
		"Base":    "github.com/example/myapp",
	}

	// Execute the template
	var generatedBuf bytes.Buffer
	if err := tmpl.Execute(&generatedBuf, data); err != nil {
		t.Fatalf("Failed to execute repository template: %v", err)
	}

	// Read the expected golden file
	expectedFile := filepath.Join("test", "expected", "repository.golden")
	expectedBytes, err := os.ReadFile(expectedFile)
	if err != nil {
		// If golden file doesn't exist, create it (useful for first run)
		if os.IsNotExist(err) {
			t.Logf("Golden file doesn't exist, creating it: %s", expectedFile)
			os.MkdirAll(filepath.Dir(expectedFile), 0755)
			os.WriteFile(expectedFile, generatedBuf.Bytes(), 0644)
			t.Logf("Test skipped - golden file created for future comparison")
			return
		}
		t.Fatalf("Failed to read golden file: %v", err)
	}

	// Compare full file contents with normalized whitespace
	generated := normalizeWhitespace(generatedBuf.String())
	expected := normalizeWhitespace(string(expectedBytes))

	if generated != expected {
		t.Errorf("Repository template output doesn't match golden file\n\nGolden file: %s\n\nExpected:\n%s\n\nGenerated:\n%s\n",
			expectedFile, string(expectedBytes), generatedBuf.String())
	} else {
		t.Logf("Repository template output exactly matches golden file")
	}
}
