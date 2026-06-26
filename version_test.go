package main

import (
	"os/exec"
	"strings"
	"testing"
)

// TestVersionDefault checks that the Version variable is non-empty.
func TestVersionDefault(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
}

// TestVersionFlagPrintsVersion builds the binary and runs it with -version,
// verifying the output matches the injected Version value.
func TestVersionFlagPrintsVersion(t *testing.T) {
	const injected = "v9.9.9-test"
	const binary = "/tmp/gohip-test-version"

	build := exec.Command(
		"go", "build",
		"-ldflags", "-X main.Version="+injected,
		"-o", binary,
		".",
	)
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	out, err := exec.Command(binary, "-version").Output()
	if err != nil {
		t.Fatalf("-version flag returned error: %v", err)
	}

	got := strings.TrimSpace(string(out))
	if got != injected {
		t.Errorf("expected %q, got %q", injected, got)
	}
}
