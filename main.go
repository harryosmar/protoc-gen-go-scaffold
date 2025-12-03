package main

import (
	"flag"

	"github.com/harryosmar/protoc-gen-go-scaffold/generator"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	var flags flag.FlagSet

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f)
		}
		return nil
	})
}

func generateFile(gen *protogen.Plugin, file *protogen.File) {
	if len(file.Messages) == 0 {
		return
	}

	// Generate handler layer
	generator.GenerateHandlers(gen, file)
	
	// Generate service layer
	generator.GenerateServices(gen, file)
	
	// Generate repository layer
	generator.GenerateRepositories(gen, file)
}
