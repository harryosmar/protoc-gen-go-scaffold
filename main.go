package main

import (
	"embed"
	"fmt"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			for _, service := range f.Services {
				// Generate service layer
				if err := generateLayer(gen, f, service, "service.tmpl",
					fmt.Sprintf("service/%s_service.go", strings.ToLower(service.GoName))); err != nil {
					return err
				}

				// Generate usecase layer
				if err := generateLayer(gen, f, service, "usecase.tmpl",
					fmt.Sprintf("usecase/%s_usecase.go", strings.ToLower(service.GoName))); err != nil {
					return err
				}

				// Generate repository layer
				if err := generateLayer(gen, f, service, "repository.tmpl",
					fmt.Sprintf("repository/%s_repository.go", strings.ToLower(service.GoName))); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func generateLayer(gen *protogen.Plugin, f *protogen.File, service *protogen.Service, tmplFile, outputPath string) error {
	g := gen.NewGeneratedFile(outputPath, f.GoImportPath)

	// Read embedded template
	tmplContent, err := templateFS.ReadFile("templates/" + tmplFile)
	if err != nil {
		return fmt.Errorf("error reading embedded template: %w", err)
	}

	tmpl, err := template.New(tmplFile).Funcs(template.FuncMap{
		"initialLower": func(s string) string {
			return strings.ToLower(s[:1]) + s[1:]
		},
	}).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	data := struct {
		Package  string
		Service  *protogen.Service
		ProtoPkg string
	}{
		Package:  string(f.GoPackageName),
		Service:  service,
		ProtoPkg: string(f.GoImportPath),
	}

	if err := tmpl.Execute(g, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}
