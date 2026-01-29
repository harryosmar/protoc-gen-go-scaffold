package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// normalizeWhitespace standardizes whitespace to allow fair comparison
func normalizeWhitespace(s string) string {
	// Replace all whitespace with single spaces
	space := strings.Fields(s)
	s = strings.Join(space, " ")
	// Remove spaces around common punctuation
	s = strings.ReplaceAll(s, " {", "{")
	s = strings.ReplaceAll(s, "{ ", "{")
	s = strings.ReplaceAll(s, " }", "}")
	s = strings.ReplaceAll(s, "} ", "}")
	return s
}

// MockService is a mock implementation for protoc-gen-star Service
type MockService struct{}

// mockStringValue is a simple wrapper for string with a String method
type mockStringValue struct {
	s string
}

// String implements the stringer interface
func (m mockStringValue) String() string {
	return m.s
}

// Name returns a string value that when used with String() gives "UserService"
func (s MockService) Name() mockStringValue {
	return mockStringValue{s: "UserService"}
}

// MockMethod represents a service method for testing
type MockMethod struct {
	NameString string
	InputName  string
	OutputName string
}

// Name returns a struct that satisfies .Name.String in templates
func (m MockMethod) Name() struct{ String string } {
	return struct{ String string }{String: m.NameString}
}

// Input returns a struct that satisfies .Input.Name in templates
func (m MockMethod) Input() struct{ Name string } {
	return struct{ Name string }{Name: m.InputName}
}

// Output returns a struct that satisfies .Output.Name in templates
func (m MockMethod) Output() struct{ Name string } {
	return struct{ Name string }{Name: m.OutputName}
}

// mockService is a more complex mock implementation for template tests
type mockService struct {
	methods []MockMethod
}

// Name implements the Name() method to satisfy Service interface
func (s *mockService) Name() mockStringValue {
	return mockStringValue{s: "UserService"}
}

// Methods returns the service methods for template iteration
func (s *mockService) Methods() []MockMethod {
	return s.methods
}

// Common test data for all templates
func getTestData() map[string]interface{} {
	return map[string]interface{}{
		"ServiceName":  "User",
		"ServiceLower": "user",
		"Import":       "github.com/example/myapp/gen/user",
		"Base":         "github.com/example/myapp",
	}
}

// compareTemplateOutput compares template output with expected content
func compareTemplateOutput(t *testing.T, templateName, generatedOutput string) {
	// Get expected file path
	expectedFile := filepath.Join("test", "expected", templateName+".golden")

	// Check if we're in update mode (create/update golden files)
	if os.Getenv("UPDATE_GOLDEN") != "" {
		err := os.MkdirAll(filepath.Dir(expectedFile), 0755)
		if err != nil {
			t.Fatalf("Failed to create expected directory: %v", err)
		}
		err = os.WriteFile(expectedFile, []byte(generatedOutput), 0644)
		if err != nil {
			t.Fatalf("Failed to write golden file: %v", err)
		}
		t.Logf("Updated golden file: %s", expectedFile)
		return
	}

	// Read the expected output
	expectedBytes, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read expected file %s: %v", expectedFile, err)
	}
	expected := string(expectedBytes)

	// Normalize both for comparison
	normalizedExpected := normalizeOutput(expected)
	normalizedGenerated := normalizeOutput(generatedOutput)

	// Compare full files
	if normalizedExpected != normalizedGenerated {
		// Show detailed diff
		diffs := showDiff(normalizedExpected, normalizedGenerated)
		t.Errorf("Template %s output doesn't match expected golden file.\n%s", templateName, diffs)
	} else {
		t.Logf("Template %s validation successful - output matches expected golden file exactly", templateName)
	}
}

// normalizeOutput removes whitespace differences that don't matter for comparison
func normalizeOutput(content string) string {
	// Replace Windows line endings with Unix line endings
	content = strings.ReplaceAll(content, "\r\n", "\n")

	// Remove trailing whitespace from each line
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	content = strings.Join(lines, "\n")

	// Remove empty lines at the start and end
	content = strings.Trim(content, "\n")

	return content
}

// showDiff shows differences between expected and generated output
func showDiff(expected, generated string) string {
	// Split into lines for line-by-line comparison
	expectedLines := strings.Split(expected, "\n")
	generatedLines := strings.Split(generated, "\n")

	var diff strings.Builder
	diff.WriteString("\nDifferences found:\n")

	// Find the maximum number of lines to compare
	maxLines := len(expectedLines)
	if len(generatedLines) > maxLines {
		maxLines = len(generatedLines)
	}

	// Compare each line
	for i := 0; i < maxLines; i++ {
		if i >= len(expectedLines) {
			diff.WriteString(fmt.Sprintf("Line %d: Missing in expected: %s\n", i+1, generatedLines[i]))
			continue
		}
		if i >= len(generatedLines) {
			diff.WriteString(fmt.Sprintf("Line %d: Missing in generated: %s\n", i+1, expectedLines[i]))
			continue
		}
		if expectedLines[i] != generatedLines[i] {
			diff.WriteString(fmt.Sprintf("Line %d:\n  Expected: %s\n  Generated: %s\n", i+1, expectedLines[i], generatedLines[i]))
		}
	}

	return diff.String()
}
