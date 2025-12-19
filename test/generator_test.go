package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceGeneration(t *testing.T) {
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

	// Verify handler filename
	expectedFilename := "user_service.go"
	serviceFilePath := filepath.Join(tmpDir, "handler", expectedFilename)
	_, err = os.Stat(serviceFilePath)
	require.NoError(t, err, "handler file %s not found", expectedFilename)

	// Read generated content
	generatedContent, err := os.ReadFile(serviceFilePath)
	require.NoError(t, err, "handler file not found")

	// Read expected content
	expectedPath := filepath.Join(cwd, "expected", "handler", "user_service.go")
	expectedContent, err := os.ReadFile(expectedPath)
	require.NoError(t, err, "expected file not found")

	// Normalize content by trimming spaces
	normalizedGenerated := strings.TrimSpace(string(generatedContent))
	normalizedExpected := strings.TrimSpace(string(expectedContent))

	// Compare normalized content
	assert.Equal(t, normalizedExpected, normalizedGenerated, "generated handler file does not match expected")
}
