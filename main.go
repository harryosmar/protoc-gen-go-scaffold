package main

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

func main() {
	// Parameters passed via --go-scaffold_out=...
	var base string     // base module path, e.g., github.com/harryosmar/protobuf-go
	var pathsOpt string // "import" or "source_relative" (default empty -> source_relative)

	// Parse plugin parameters
	opts := protogen.Options{
		ParamFunc: func(name, value string) error {
			switch name {
			case "base":
				base = value
			case "paths":
				pathsOpt = value
			}
			return nil
		},
	}

	// Run the generator
	opts.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		// Default to source_relative if not provided
		if pathsOpt == "" {
			pathsOpt = "source_relative"
		}

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			// If base is missing, try to infer it from proto import path
			baseForFile := base
			if baseForFile == "" {
				baseForFile = inferBaseFromGoImportPath(f)
			}

			for _, service := range f.Services {
				// Generate SERVICE layer
				if err := generateLayer(gen, f, service, baseForFile, pathsOpt,
					"server.tmpl", "server", ".go"); err != nil {
					return err
				}

				// Uncomment when templates ready:
				// Generate USECASE layer
				// if err := generateLayer(gen, f, service, baseForFile, pathsOpt,
				//     "usecase.tmpl", "usecase", "_usecase.go"); err != nil {
				//     return err
				// }
				// Generate REPOSITORY layer
				// if err := generateLayer(gen, f, service, baseForFile, pathsOpt,
				//     "repository.tmpl", "repository", "_repository.go"); err != nil {
				//     return err
				// }
			}
		}
		return nil
	})
}

// inferBaseFromGoImportPath tries to derive the base module from a typical go_package import path.
// Example: "github.com/harryosmar/protobuf-go/gen/hello" -> "github.com/harryosmar/protobuf-go"
func inferBaseFromGoImportPath(f *protogen.File) string {
	path := string(f.GoImportPath)
	if i := strings.Index(path, "/gen/"); i != -1 {
		return path[:i]
	}
	// Best-effort fallback: trim last segment
	if j := strings.LastIndex(path, "/"); j != -1 {
		return path[:j]
	}
	return path
}

// generateLayer chooses output placement based on pathsOpt:
// - "import": writes under proto package dir (gen/<pkg>/...)
// - "source_relative": writes to repo-root relative paths (e.g., service/*.go, usecase/*.go)
func generateLayer(
	gen *protogen.Plugin,
	f *protogen.File,
	service *protogen.Service,
	base string,
	pathsOpt string, // "import" or "source_relative"
	tmplFile string,
	subdir string, // e.g., "service" | "usecase" | "repository"
	suffix string, // e.g., ".go"
) error {
	lowerName := snakeCase(service.GoName)
	fileName := lowerName + suffix

	var finalPath string
	var importPathArg protogen.GoImportPath

	if pathsOpt == "import" {
		// Place under proto package dir to satisfy protogen prefix check
		pkgRel := string(f.GoImportPath) // e.g., "github.com/.../gen/hello"
		if base != "" && strings.HasPrefix(pkgRel, base+"/") {
			pkgRel = strings.TrimPrefix(pkgRel, base+"/") // -> "gen/hello"
		} else if i := strings.Index(pkgRel, "/gen/"); i != -1 {
			pkgRel = pkgRel[i+1:] // -> "gen/hello"
		} else if j := strings.LastIndex(pkgRel, "/"); j != -1 {
			pkgRel = pkgRel[j+1:] // last segment
		}
		finalPath = filepath.Join(pkgRel, subdir, fileName)
		importPathArg = f.GoImportPath // enforce prefix constraint
	} else {
		// paths=source_relative -> write to repo root (relative) and disable prefix check
		finalPath = filepath.Join(subdir, fileName)
		importPathArg = "" // IMPORTANT: disable prefix check
	}

	// Register generated file
	g := gen.NewGeneratedFile(finalPath, importPathArg)

	// Load template
	tmplContent, err := templateFS.ReadFile("templates/" + tmplFile)
	if err != nil {
		return fmt.Errorf("error reading embedded template: %w", err)
	}

	// Template helpers
	funcs := template.FuncMap{
		"initialLower": func(s string) string {
			if s == "" {
				return s
			}
			r := []rune(s)
			r[0] = []rune(strings.ToLower(string(r[0])))[0]
			return string(r)
		},
		"trimSuffix": strings.TrimSuffix,
		"lower":      strings.ToLower,
		"join":       func(parts ...string) string { return strings.Join(parts, "/") },
		// Qualify a GoIdent with the correct import alias and register import in this file
		"qual": func(id protogen.GoIdent) string {
			return g.QualifiedGoIdent(id)
		},
		// Return "context.Context" with import registered
		"ctx": func() string {
			return g.QualifiedGoIdent(protogen.GoIdent{
				GoName:       "Context",
				GoImportPath: "context",
			})
		},
	}

	// Parse template
	tmpl, err := template.New(tmplFile).Funcs(funcs).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("error parsing template (%s): %w", tmplFile, err)
	}

	// Package name:
	// - For "import": use the proto package name (e.g., "hello")
	// - For "source_relative": use the layer folder name (service/server/usecase/repository)
	pkgName := string(f.GoPackageName)
	if pathsOpt != "import" {
		switch subdir {
		case "service":
			pkgName = "service"
		case "server":
			pkgName = "server"
		case "usecase":
			pkgName = "usecase"
		case "repository":
			pkgName = "repository"
		}
	}

	// Template data
	data := struct {
		Filename     string
		Package      string
		Service      *protogen.Service
		ProtoPkg     string // e.g., "github.com/.../gen/hello"
		Base         string // base module: e.g., "github.com/harryosmar/protobuf-go"
		ServiceLower string
		ServiceNoSuf string // e.g., "Hello" if "HelloHandler"
	}{
		Package:      pkgName,
		Service:      service,
		ProtoPkg:     string(f.GoImportPath),
		Base:         base,
		ServiceLower: lowerName,
		ServiceNoSuf: strings.TrimSuffix(service.GoName, "Service"),
	}

	// Execute template
	if err := tmpl.Execute(g, data); err != nil {
		return fmt.Errorf("error executing template (%s): %w", tmplFile, err)
	}
	return nil
}

func snakeCase(s string) string {
	var buf bytes.Buffer
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			buf.WriteRune('_')
		}
		buf.WriteRune(unicode.ToLower(r))
	}
	return buf.String()
}
