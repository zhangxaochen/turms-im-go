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
	javaMethodRegex = regexp.MustCompile(`(?m)^\s*(?:@\w+(?:\([^)]*\))?\s*)*public\s+(?:static\s+|final\s+|abstract\s+|default\s+)*(?:<[^>]+>\s+)?(?:[\w<>[\]?,\s]+\s+)+(\w+)\s*\(([^)]*)\)`)
	goMethodRegex   = regexp.MustCompile(`(?m)((?:^\s*//\s*@MappedFrom\s+[^\n]+\n)+)?^func\s+(?:\([^)]+\)\s+)?(\w+)\s*\(([^)]*)\)`)
)

type MethodDef struct {
	Name string
	Args string
}

func cleanWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func extractPublicMethods(code string) []MethodDef {
	matches := javaMethodRegex.FindAllStringSubmatch(code, -1)
	var methods []MethodDef
	for _, match := range matches {
		methods = append(methods, MethodDef{Name: match[1], Args: cleanWhitespace(match[2])})
	}
	return methods
}

func extractGoMethods(code string) []GoMethod {
	matches := goMethodRegex.FindAllStringSubmatch(code, -1)
	var methods []GoMethod
	mappedFromRegex := regexp.MustCompile(`@MappedFrom\s+([^\n]+)`)

	for _, match := range matches {
		mappedFromBlock := match[1]
		name := match[2]
		args := cleanWhitespace(match[3])

		if mappedFromBlock != "" {
			lines := strings.Split(mappedFromBlock, "\n")
			hasAtLeastOne := false
			for _, line := range lines {
				if strings.TrimSpace(line) == "" {
					continue
				}
				subMatch := mappedFromRegex.FindStringSubmatch(line)
				if len(subMatch) > 1 {
					hasAtLeastOne = true
					methods = append(methods, GoMethod{
						MappedFrom: strings.TrimSpace(subMatch[1]),
						Name:       name,
						Args:       args,
					})
				}
			}
			if !hasAtLeastOne {
				methods = append(methods, GoMethod{Name: name, Args: args})
			}
		} else {
			methods = append(methods, GoMethod{
				MappedFrom: "",
				Name:       name,
				Args:       args,
			})
		}
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
	Methods      []MethodDef
}

type JavaConfig struct {
	RelativePath string
	Name         string
}

type GoMethod struct {
	Path       string
	Name       string
	Args       string
	MappedFrom string
}

func main() {
	goMethodMap := make(map[string][]GoMethod) // method name (lowercase) -> all its implementations

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
						if m.MappedFrom != "" {
							idx := strings.Index(m.MappedFrom, "(")
							if idx > 0 {
								javaName := strings.TrimSpace(m.MappedFrom[:idx])
								lowerM := strings.ToLower(javaName)
								m.Path = path
								goMethodMap[lowerM] = append(goMethodMap[lowerM], m)
								continue
							}
						}
						
						lowerM := strings.ToLower(m.Name)
						if !ignoreMethods[lowerM] {
							m.Path = path
							goMethodMap[lowerM] = append(goMethodMap[lowerM], m)
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
				goMethods, exists := goMethodMap[strings.ToLower(m.Name)]
				if exists && len(goMethods) > 0 {
					var exactMatches []GoMethod
					var fuzzyMatches []GoMethod
					javaSig := fmt.Sprintf("%s(%s)", m.Name, m.Args)

					for _, gm := range goMethods {
						if gm.MappedFrom != "" {
							if cleanWhitespace(gm.MappedFrom) == cleanWhitespace(javaSig) {
								exactMatches = append(exactMatches, gm)
							}
						} else {
							fuzzyMatches = append(fuzzyMatches, gm)
						}
					}

					var matchesToUse []GoMethod
					if len(exactMatches) > 0 {
						matchesToUse = exactMatches
					} else {
						matchesToUse = fuzzyMatches
					}

					if len(matchesToUse) > 0 {
						var goStrs []string
						for _, gm := range matchesToUse {
							goStrs = append(goStrs, fmt.Sprintf("`%s:%s(%s)`", gm.Path, gm.Name, gm.Args))
						}
						sb.WriteString(fmt.Sprintf("  - [x] `%s(%s)` -> %s\n", m.Name, m.Args, strings.Join(goStrs, ", ")))
					} else {
						sb.WriteString(fmt.Sprintf("  - [ ] `%s(%s)`\n", m.Name, m.Args))
					}
				} else {
					sb.WriteString(fmt.Sprintf("  - [ ] `%s(%s)`\n", m.Name, m.Args))
				}
			}
			sb.WriteString("\n")
		}
	}

	os.MkdirAll("docs", 0755)
	os.WriteFile("docs/refactor_progress_report.md", []byte(sb.String()), 0644)
	fmt.Println("Generated docs/refactor_progress_report.md successfully.")
}
