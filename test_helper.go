package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// normalizeContent removes whitespace differences that don't matter for comparison
func normalizeContent(content string) string {
	// Replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	content = re.ReplaceAllString(content, " ")
	
	// Trim leading/trailing whitespace
	content = strings.TrimSpace(content)
	
	return content
}

// checkMissingElements helps identify what key elements are missing from the generated content
func checkMissingElements(t *testing.T, expected, actual string) {
	// Common key elements to check when there's a mismatch
	keyElements := []string{
		"package ",
		"import",
		"type ",
		"struct {",
		"func New",
		"return &",
	}

	for _, element := range keyElements {
		if strings.Contains(expected, element) && !strings.Contains(actual, element) {
			t.Errorf("Generated code missing key element: %s", element)
		}
	}
}

// loadExpectedFile reads an expected file for comparison
func loadExpectedFile(t *testing.T, filename string) string {
	expectedPath := filepath.Join("test", "expected", filename)
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Logf("Warning: Expected file %s doesn't exist or can't be read: %v", expectedPath, err)
		return ""
	}
	return string(content)
}
