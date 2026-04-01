package main

import "testing"

func TestParseJavaMethods(t *testing.T) {
	javaCode := `
package test;

public class MyClass {
	public void myMethod() {}
	private int internalMethod() { return 1; }
	public String getStatus(String id) { return id; }
	public static final boolean isEnabled() { return true; }
	public <T> T genericMethod(T arg) { return arg; }
}
`
	methods := extractPublicMethods(javaCode)
	if len(methods) != 4 {
		t.Fatalf("Expected 4 public methods, got %d. Parsed: %v", len(methods), methods)
	}

	expectedMethods := map[string]bool{
		"myMethod":      true,
		"getStatus":     true,
		"isEnabled":     true,
		"genericMethod": true,
	}

	for _, m := range methods {
		if !expectedMethods[m] {
			t.Errorf("Unexpected method parsed: %s", m)
		}
	}
}

func TestParseGoMethods(t *testing.T) {
	goCode := `
package domain

func MyGoMethod() {}
func (s *Service) ClassMethod() error { return nil }
func privateMethod() {}
`
	methods := extractGoMethods(goCode)
	if len(methods) != 3 {
		t.Fatalf("Expected 3 Go methods, got %d. Parsed: %v", len(methods), methods)
	}

	expectedMethods := map[string]bool{
		"MyGoMethod":    true,
		"ClassMethod":   true,
		"privateMethod": true,
	}

	for _, m := range methods {
		if !expectedMethods[m] {
			t.Errorf("Unexpected Go method parsed: %s", m)
		}
	}
}

