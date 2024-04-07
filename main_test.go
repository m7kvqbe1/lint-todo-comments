package main

import (
	"os/exec"
	"strings"
	"testing"
)

func runTodoChecker(t *testing.T, dir string) string {
	t.Helper()

	cmd := exec.Command("../lint-todo-comments", dir)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("todochecker failed: %s", err)
	}

	return string(output)
}

func TestWithTodoComments(t *testing.T) {
	output := runTodoChecker(t, "testdata/with_todo")

	if !strings.Contains(output, "TODO") {
		t.Error("Expected to find TODO comments, but none were found")
	}
}

func TestWithoutTodoComments(t *testing.T) {
	output := runTodoChecker(t, "testdata/without_todo")

	if strings.Contains(output, "TODO") {
		t.Error("Expected no TODO comments, but some were found")
	}
}

func TestMixedContents(t *testing.T) {
	output := runTodoChecker(t, "testdata/mixed_contents")

	if !strings.Contains(output, "TODO") {
		t.Error("Expected to find TODO comments, but none were found")
	}
}
