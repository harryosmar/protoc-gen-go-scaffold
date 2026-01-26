package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerGeneration(t *testing.T) {
	// Get current working directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Setup temporary directory
	tmpDir, err := os.MkdirTemp("", "proto-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Use absolute path for test proto
	testProto := filepath.Join(cwd, "testdata", "user.proto")
	protoData, err := os.ReadFile(testProto)
	require.NoError(t, err)

	tmpProto := filepath.Join(tmpDir, "user.proto")
	err = os.WriteFile(tmpProto, protoData, 0644)
	require.NoError(t, err)

	// Build plugin
	pluginPath := filepath.Join(cwd, "protoc-gen-go-scaffold")
	buildCmd := exec.Command("go", "build", "-o", pluginPath, "../main.go")
	buildOutput, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "plugin build failed: %s", buildOutput)
	defer os.Remove(pluginPath)

	// Run protoc with our plugin
	cmd := exec.Command("protoc",
		"--plugin=protoc-gen-go-scaffold="+pluginPath,
		"--proto_path=",
		"--proto_path="+tmpDir,
		"--go-scaffold_out=base=github.com/harryosmar/protobuf-go,paths=source_relative:"+tmpDir,
		filepath.Base(tmpProto),
	)
	cmd.Dir = tmpDir

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "protoc execution failed: %s", output)

	// Verify service filename
	expectedFilename := "user_service.go"
	serviceFilePath := filepath.Join(tmpDir, "server", expectedFilename)
	_, err = os.Stat(serviceFilePath)
	require.NoError(t, err, "service file %s not found", expectedFilename)

	// Read generated content
	generatedContent, err := os.ReadFile(serviceFilePath)
	require.NoError(t, err, "service file not found")

	// Read expected content
	expectedPath := filepath.Join(cwd, "expected", "server", "user_service.go")
	expectedContent, err := os.ReadFile(expectedPath)
	require.NoError(t, err, "expected file not found")

	// Normalize content by trimming spaces
	normalizedGenerated := strings.TrimSpace(string(generatedContent))
	normalizedExpected := strings.TrimSpace(string(expectedContent))

	// Compare normalized content
	assert.Equal(t, normalizedExpected, normalizedGenerated, "generated service file does not match expected")
}

func TestUsecaseGeneration(t *testing.T) {
	// Get current working directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Setup temporary directory
	tmpDir, err := os.MkdirTemp("", "proto-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create templates directory
	tmplDir := filepath.Join(tmpDir, "templates")
	err = os.MkdirAll(tmplDir, 0755)
	require.NoError(t, err)

	// Copy template files
	templatesPath := filepath.Join(cwd, "..", "templates")
	templateFiles := []string{"server.tmpl", "usecase.tmpl", "repository.tmpl"}
	for _, file := range templateFiles {
		srcPath := filepath.Join(templatesPath, file)
		dstPath := filepath.Join(tmplDir, file)
		tmplData, err := os.ReadFile(srcPath)
		if err == nil { // Only copy if file exists
			err = os.WriteFile(dstPath, tmplData, 0644)
			require.NoError(t, err)
		}
	}

	// Use absolute path for test proto
	testProto := filepath.Join(cwd, "testdata", "user.proto")
	protoData, err := os.ReadFile(testProto)
	require.NoError(t, err)

	tmpProto := filepath.Join(tmpDir, "user.proto")
	err = os.WriteFile(tmpProto, protoData, 0644)
	require.NoError(t, err)

	// Build plugin
	pluginPath := filepath.Join(cwd, "protoc-gen-go-scaffold")
	buildCmd := exec.Command("go", "build", "-o", pluginPath, "../main.go")
	buildOutput, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "plugin build failed: %s", buildOutput)
	defer os.Remove(pluginPath)

	// Modify main.go to enable usecase generation
	mainPath := filepath.Join(cwd, "..", "main.go")
	mainData, err := os.ReadFile(mainPath)
	require.NoError(t, err)

	// Create a modified version of main.go with usecase generation uncommented
	modifiedMain := strings.Replace(
		string(mainData),
		"// Generate USECASE layer\n\t\t\t\t// if err := generateLayer(gen, f, service, baseForFile, pathsOpt,\n\t\t\t\t//     \"usecase.tmpl\", \"usecase\", \"_usecase.go\"); err != nil {\n\t\t\t\t//     return err\n\t\t\t\t// }",
		"// Generate USECASE layer\n\t\t\t\tif err := generateLayer(gen, f, service, baseForFile, pathsOpt,\n\t\t\t\t    \"usecase.tmpl\", \"usecase\", \"_usecase.go\"); err != nil {\n\t\t\t\t    return err\n\t\t\t\t}",
		1,
	)

	tmpMainPath := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(tmpMainPath, []byte(modifiedMain), 0644)
	require.NoError(t, err)

	// Rebuild plugin with modified main.go
	buildCmd = exec.Command("go", "build", "-o", pluginPath, tmpMainPath)
	buildOutput, err = buildCmd.CombinedOutput()
	require.NoError(t, err, "plugin build failed: %s", buildOutput)

	// Run protoc with our plugin
	cmd := exec.Command("protoc",
		"--plugin=protoc-gen-go-scaffold="+pluginPath,
		"--proto_path=",
		"--proto_path="+tmpDir,
		"--go-scaffold_out=base=github.com/harryosmar/protobuf-go,paths=source_relative:"+tmpDir,
		filepath.Base(tmpProto),
	)
	cmd.Dir = tmpDir

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "protoc execution failed: %s", output)

	// Verify usecase filename
	expectedFilename := "user_service_usecase.go"
	usecaseFilePath := filepath.Join(tmpDir, "usecase", expectedFilename)
	_, err = os.Stat(usecaseFilePath)
	require.NoError(t, err, "usecase file %s not found", expectedFilename)

	// Read generated content
	generatedContent, err := os.ReadFile(usecaseFilePath)
	require.NoError(t, err, "usecase file not found")

	// Read expected content
	expectedPath := filepath.Join(cwd, "expected", "usecase", "user_service.go")
	expectedContent, err := os.ReadFile(expectedPath)
	require.NoError(t, err, "expected file not found")

	// Normalize content by removing extra whitespace
	normalizedGenerated := normalizeWhitespace(string(generatedContent))
	normalizedExpected := normalizeWhitespace(string(expectedContent))

	// Compare normalized content
	assert.Equal(t, normalizedExpected, normalizedGenerated, "generated usecase file does not match expected")
}

// normalizeWhitespace removes extra whitespace and normalizes newlines
func normalizeWhitespace(s string) string {
	// Trim spaces
	s = strings.TrimSpace(s)
	// Replace multiple newlines with a single newline
	s = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(s, "\n")
	// Replace tabs with spaces
	s = strings.ReplaceAll(s, "\t", "    ")
	return s
}

func TestRepositoryGeneration(t *testing.T) {
	// Get current working directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Setup temporary directory
	tmpDir, err := os.MkdirTemp("", "proto-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create templates directory
	tmplDir := filepath.Join(tmpDir, "templates")
	err = os.MkdirAll(tmplDir, 0755)
	require.NoError(t, err)

	// Copy template files
	templatesPath := filepath.Join(cwd, "..", "templates")
	templateFiles := []string{"server.tmpl", "usecase.tmpl", "repository.tmpl"}
	for _, file := range templateFiles {
		srcPath := filepath.Join(templatesPath, file)
		dstPath := filepath.Join(tmplDir, file)
		tmplData, err := os.ReadFile(srcPath)
		if err == nil { // Only copy if file exists
			err = os.WriteFile(dstPath, tmplData, 0644)
			require.NoError(t, err)
		}
	}

	// Use absolute path for test proto
	testProto := filepath.Join(cwd, "testdata", "user.proto")
	protoData, err := os.ReadFile(testProto)
	require.NoError(t, err)

	tmpProto := filepath.Join(tmpDir, "user.proto")
	err = os.WriteFile(tmpProto, protoData, 0644)
	require.NoError(t, err)

	// Build plugin
	pluginPath := filepath.Join(cwd, "protoc-gen-go-scaffold")
	buildCmd := exec.Command("go", "build", "-o", pluginPath, "../main.go")
	buildOutput, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "plugin build failed: %s", buildOutput)
	defer os.Remove(pluginPath)

	// Modify main.go to enable repository generation
	mainPath := filepath.Join(cwd, "..", "main.go")
	mainData, err := os.ReadFile(mainPath)
	require.NoError(t, err)

	// Create a modified version of main.go with repository generation uncommented
	modifiedMain := strings.Replace(
		string(mainData),
		"// Generate REPOSITORY layer\n\t\t\t\t// if err := generateLayer(gen, f, service, baseForFile, pathsOpt,\n\t\t\t\t//     \"repository.tmpl\", \"repository\", \"_repository.go\"); err != nil {\n\t\t\t\t//     return err\n\t\t\t\t// }",
		"// Generate REPOSITORY layer\n\t\t\t\tif err := generateLayer(gen, f, service, baseForFile, pathsOpt,\n\t\t\t\t    \"repository.tmpl\", \"repository\", \"_repository.go\"); err != nil {\n\t\t\t\t    return err\n\t\t\t\t}",
		1,
	)

	tmpMainPath := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(tmpMainPath, []byte(modifiedMain), 0644)
	require.NoError(t, err)

	// Rebuild plugin with modified main.go
	buildCmd = exec.Command("go", "build", "-o", pluginPath, tmpMainPath)
	buildOutput, err = buildCmd.CombinedOutput()
	require.NoError(t, err, "plugin build failed: %s", buildOutput)

	// Run protoc with our plugin
	cmd := exec.Command("protoc",
		"--plugin=protoc-gen-go-scaffold="+pluginPath,
		"--proto_path=",
		"--proto_path="+tmpDir,
		"--go-scaffold_out=base=github.com/harryosmar/protobuf-go,paths=source_relative:"+tmpDir,
		filepath.Base(tmpProto),
	)
	cmd.Dir = tmpDir

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "protoc execution failed: %s", output)

	// Verify repository filename
	expectedFilename := "user_service_repository.go"
	repositoryFilePath := filepath.Join(tmpDir, "repository", expectedFilename)
	_, err = os.Stat(repositoryFilePath)
	require.NoError(t, err, "repository file %s not found", expectedFilename)

	// Read generated content
	generatedContent, err := os.ReadFile(repositoryFilePath)
	require.NoError(t, err, "repository file not found")

	// Read expected content
	expectedPath := filepath.Join(cwd, "expected", "repository", "user_service.go")
	expectedContent, err := os.ReadFile(expectedPath)
	require.NoError(t, err, "expected file not found")

	// Normalize content by removing extra whitespace
	normalizedGenerated := normalizeWhitespace(string(generatedContent))
	normalizedExpected := normalizeWhitespace(string(expectedContent))

	// Compare normalized content
	assert.Equal(t, normalizedExpected, normalizedGenerated, "generated repository file does not match expected")
}
