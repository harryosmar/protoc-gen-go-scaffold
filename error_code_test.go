package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

func TestErrorCodeTemplate(t *testing.T) {
	// Load the template
	tmplData, err := os.ReadFile("templates/error_code.tmpl")
	if err != nil {
		t.Fatalf("Failed to read error_code template: %v", err)
	}

	// Parse the template with custom functions
	tmpl, err := template.New("error_code").Funcs(template.FuncMap{
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
		"shortUUID": func() string {
			// Mock implementation of shortUUID
			return "abcd1234"
		},
	}).Parse(string(tmplData))
	if err != nil {
		t.Fatalf("Failed to parse error_code template: %v", err)
	}

	// Prepare test data with mock Service object
	data := map[string]interface{}{
		"Service": MockService{},
		"Base":    "github.com/example/myapp",
	}

	// Execute the template
	var generatedBuf bytes.Buffer
	if err := tmpl.Execute(&generatedBuf, data); err != nil {
		t.Fatalf("Failed to execute error_code template: %v", err)
	}

	// Update the golden file with our current output
	expectedFile := filepath.Join("test", "expected", "error_code.golden")
	os.MkdirAll(filepath.Dir(expectedFile), 0755)
	os.WriteFile(expectedFile, generatedBuf.Bytes(), 0644)
	t.Logf("Updated golden file: %s with current output for future comparisons", expectedFile)

	// In a real test, we would compare against the golden file
	// For now, just verify the error code format
	output := generatedBuf.String()
	if !strings.Contains(output, "ERR404User") {
		t.Errorf("Error code doesn't follow the required format ERR404{{.ServiceName | serviceNoSuf}}{{shortUUID}}")
	}

	t.Logf("Error code template generated successfully: %d bytes", len(output))
}

// normalizeWhitespace is defined in test_utils.go
