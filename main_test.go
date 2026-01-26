package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestGenerationV2(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"protoc-gen-go-scaffold"}

	tmpDir, err := os.MkdirTemp("", "proto-test-v2-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	descPath := filepath.Join(tmpDir, "desc.pb")
	protoDir := filepath.Join(cwd, "test", "testdata")

	cmd := exec.Command("protoc",
		"--descriptor_set_out="+descPath,
		"--include_imports",
		"--proto_path="+protoDir,
		"user.proto",
	)
	cmd.Dir = protoDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "protoc execution failed: %s", output)

	descBytes, err := os.ReadFile(descPath)
	require.NoError(t, err)

	fds := &descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(descBytes, fds)
	require.NoError(t, err)

	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{"user.proto"},
		ProtoFile:      fds.File,
		Parameter:      proto.String("base=github.com/harryosmar/protobuf-go,paths=source_relative"),
	}

	reqBytes, err := proto.Marshal(req)
	require.NoError(t, err)

	origStdin := os.Stdin
	origStdout := os.Stdout
	defer func() {
		os.Stdin = origStdin
		os.Stdout = origStdout
	}()

	stdinR, stdinW, err := os.Pipe()
	require.NoError(t, err)
	stdoutR, stdoutW, err := os.Pipe()
	require.NoError(t, err)

	os.Stdin = stdinR
	os.Stdout = stdoutW

	writeErrCh := make(chan error, 1)
	go func() {
		_, err := stdinW.Write(reqBytes)
		_ = stdinW.Close()
		writeErrCh <- err
	}()

	main()

	_ = stdoutW.Close()
	respBytes, err := io.ReadAll(stdoutR)
	require.NoError(t, err)

	require.NoError(t, <-writeErrCh)

	resp := &pluginpb.CodeGeneratorResponse{}
	err = proto.Unmarshal(respBytes, resp)
	require.NoError(t, err)

	require.Empty(t, resp.Error, "generator returned error: %s", resp.GetError())

	for _, f := range resp.File {
		name := f.GetName()
		require.NotEmpty(t, name)
		outPath := filepath.Join(tmpDir, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(outPath), 0755))
		require.NoError(t, os.WriteFile(outPath, []byte(f.GetContent()), 0644))
	}
}
