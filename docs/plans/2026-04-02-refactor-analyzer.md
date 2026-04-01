# Refactor Analyzer Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create a Go script (`tools/refactor_analyzer`) to automatically analyze `turms-orig` (Java) and verify what has been re-implemented in `turms-go` (Go), generating a detailed JSON/Markdown progress report.

**Architecture:** A simple Go command line tool utilizing standard library `regexp` to parse Java interfaces and class logic for method names, then scanning `.go` files to check code coverage mappings. It outputs a MarkDown file with checked `[x]` and unchecked `[ ]` boxes. After generating, an AI agent enhances the MarkDown with descriptions.

**Tech Stack:** Go `regexp`, `os`, `path/filepath`.

---

### Task 1: Scaffolding and Java Parser Test

**Files:**
- Create: `tools/refactor_analyzer/analyzer_test.go`
- Create: `tools/refactor_analyzer/analyzer.go`

**Step 1: Write the failing test**

```go
package main

import "testing"

func TestParseJavaMethods(t *testing.T) {
	javaCode := `
package test;

public class MyClass {
	public void myMethod() {}
	private int internalMethod() { return 1; }
	public String getStatus(String id) { return id; }
}
`
	methods := extractPublicMethods(javaCode)
	if len(methods) != 2 {
		t.Fatalf("Expected 2 public methods, got %d", len(methods))
	}
	if methods[0] != "myMethod" || methods[1] != "getStatus" {
		t.Errorf("Unexpected methods parsed: %v", methods)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test -v ./tools/refactor_analyzer`
Expected: FAIL, missing `extractPublicMethods`.

**Step 3: Write minimal implementation**

```go
package main

import (
	"regexp"
)

func extractPublicMethods(code string) []string {
	// Simple regex to catch generic public method names
	regex := regexp.MustCompile(`(?m)^\s*(?:@\w+\s*)*public\s+(?:static\s+)?(?:<[^>]+>\s+)?(?:[\w<>[\]]+\s+)+(\w+)\s*\(`)
	matches := regex.FindAllStringSubmatch(code, -1)
	
	var methods []string
	for _, match := range matches {
		methods = append(methods, match[1])
	}
	return methods
}
```

**Step 4: Run test to verify it passes**

Run: `go test -v ./tools/refactor_analyzer`
Expected: PASS.

**Step 5: Commit**

```bash
git add tools/refactor_analyzer/
git commit -m "feat(tools): add basic java parser for refactor_analyzer"
```

### Task 2: Go Parser and Matcher

**Files:**
- Modify: `tools/refactor_analyzer/analyzer.go`
- Modify: `tools/refactor_analyzer/analyzer_test.go`

**Step 1: Write failing test for Go parser**

Add `TestCheckGoImplementation` mapping a list of target methods to real occurrences in Go source files via string matching or regexp.

**Step 2: Run test**

Run: `go test -v ./tools/refactor_analyzer`

**Step 3: Implement matching logic**

Create a `GoScanner` that indexes all declared funcs in `turms-go` source code via regexp like `func\s+(?:\([^)]+\)\s+)?(\w+)`.

**Step 4: Run test**

Expected: PASS.

**Step 5: Commit**

### Task 3: Directory Walker and Markdown Generator

**Step 1 & 2: Tests**
Test generating the Markdown tree format layout based on extracted data.

**Step 3: Impement Generator**
Implement file walking over `turms-orig/turms-service/...` and `turms-orig/turms-gateway/...`. Match each parsed file against `GoScanner`. Outputs a file `docs/refactor_progress_report.md`.

**Step 4: End-to-end Run**
Build the script. Run `go run ./tools/refactor_analyzer`. Verify `docs/refactor_progress_report.md` exists and contains correct mappings.

**Step 5: Commit**

### Task 4: AI Context Enrichment

**Step 1:** AI reads the output MarkDown and manually edits it using file editing tools to add functionality descriptions (`[简述功能]`) for config files and domain modules directly directly into the new `docs/refactor_progress_report.md`.
