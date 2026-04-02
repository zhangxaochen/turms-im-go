package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	javaMethodRegex = regexp.MustCompile(`(?m)^\s*(?:@\w+\s*)*public\s+(?:static\s+|final\s+|abstract\s+|default\s+)*(?:<[^>]+>\s+)?(?:[\w<>[\]?,\s]+\s+)+(\w+)\s*\(`)
	goMethodRegex   = regexp.MustCompile(`(?m)^func\s+(?:\([^)]+\)\s+)?(\w+)\s*\(`)
)

func extractPublicMethods(code string) []string {
	matches := javaMethodRegex.FindAllStringSubmatch(code, -1)
	var methods []string
	for _, match := range matches {
		methods = append(methods, match[1])
	}
	return methods
}

func extractGoMethods(code string) []string {
	matches := goMethodRegex.FindAllStringSubmatch(code, -1)
	var methods []string
	for _, match := range matches {
		methods = append(methods, match[1])
	}
	return methods
}

type JavaModule struct {
	Name    string
	Files   []*JavaFile
	Configs []*JavaConfig
}

type JavaFile struct {
	RelativePath string
	ClassName    string
	Methods      []string
}

type JavaConfig struct {
	RelativePath string
	Name         string
}

func main() {
	goMethodMap := make(map[string]string) // method name (lowercase) -> relative go file path

	ignoreMethods := map[string]bool{
		"init":     true,
		"main":     true,
		"tostring": true,
		"equals":   true,
		"hashcode": true,
		"clone":    true,
	}

	goRoots := []string{"internal", "pkg", "tests"}
	for _, root := range goRoots {
		filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
				content, err := os.ReadFile(path)
				if err == nil {
					methods := extractGoMethods(string(content))
					for _, m := range methods {
						lowerM := strings.ToLower(m)
						if !ignoreMethods[lowerM] {
							goMethodMap[lowerM] = path
						}
					}
				}
			}
			return nil
		})
	}

	javaRoots := []string{
		"turms-orig/turms-service/src/main",
		"turms-orig/turms-gateway/src/main",
	}

	modules := map[string]*JavaModule{} // root path -> module

	for _, root := range javaRoots {
		if _, err := os.Stat(root); os.IsNotExist(err) {
			continue // skip if orig repo isn't cloned
		}

		modName := "turms-service"
		if strings.Contains(root, "turms-gateway") {
			modName = "turms-gateway"
		}

		if _, ok := modules[modName]; !ok {
			modules[modName] = &JavaModule{Name: modName}
		}
		mod := modules[modName]

		filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}

			rel, _ := filepath.Rel(root, path)

			if strings.HasSuffix(path, ".java") {
				content, err := os.ReadFile(path)
				if err == nil {
					methods := extractPublicMethods(string(content))
					if len(methods) > 0 {
						mod.Files = append(mod.Files, &JavaFile{
							RelativePath: rel,
							ClassName:    filepath.Base(path),
							Methods:      methods,
						})
					}
				}
			} else if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".properties") {
				mod.Configs = append(mod.Configs, &JavaConfig{
					RelativePath: rel,
					Name:         filepath.Base(path),
				})
			}
			return nil
		})
	}

	// Generate Markdown
	var sb strings.Builder
	sb.WriteString("# Turms Refactoring Progress Report\n\n")

	sb.WriteString("## Modules\n\n")

	for _, modName := range []string{"turms-gateway", "turms-service"} {
		mod := modules[modName]
		if mod == nil {
			continue
		}
		sb.WriteString(fmt.Sprintf("### %s\n\n", mod.Name))
		sb.WriteString("> [简述功能]\n\n")

		if len(mod.Configs) > 0 {
			sb.WriteString("#### Configurations\n\n")
			for _, cfg := range mod.Configs {
				sb.WriteString(fmt.Sprintf("- **%s** (`%s`): [简述功能]\n", cfg.Name, cfg.RelativePath))
			}
			sb.WriteString("\n")
		}

		// Sort files
		sort.Slice(mod.Files, func(i, j int) bool {
			return mod.Files[i].RelativePath < mod.Files[j].RelativePath
		})

		sb.WriteString("#### Java source tracking\n\n")
		for _, f := range mod.Files {
			sb.WriteString(fmt.Sprintf("- **%s** (`%s`)\n", f.ClassName, f.RelativePath))
			sb.WriteString("> [简述功能]\n\n")
			for _, m := range f.Methods {
				goPath, exists := goMethodMap[strings.ToLower(m)]
				if exists {
					sb.WriteString(fmt.Sprintf("  - [x] `%s` -> `%s`\n", m, goPath))
				} else {
					sb.WriteString(fmt.Sprintf("  - [ ] `%s`\n", m))
				}
			}
			sb.WriteString("\n")
		}
	}

	os.MkdirAll("docs", 0755)
	os.WriteFile("docs/refactor_progress_report.md", []byte(sb.String()), 0644)
	fmt.Println("Generated docs/refactor_progress_report.md successfully.")
}
