// main_test.go
package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

// TestMain sets up the helper process. When the environment variable
// GO_WANT_HELPER_PROCESS is set to "1", we call main() and exit.
// Otherwise, we run the tests.
func TestMain(m *testing.M) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		// In helper mode: call main() as if it were run from the command line.
		main()
		os.Exit(0)
	}
	// Run the normal tests.
	os.Exit(m.Run())
}

// compareOutputWithFile compares command output with the contents of an expected output file
func compareOutputWithFile(t *testing.T, output, expectedOutputPath string) {
	expected, err := os.ReadFile(expectedOutputPath)
	if err != nil {
		t.Fatalf("failed to read expected output file %s: %v", expectedOutputPath, err)
	}

	// Normalize line endings and trim whitespace
	expectedStr := string(bytes.TrimSpace(expected))
	outputStr := string(bytes.TrimSpace([]byte(output)))

	if expectedStr != outputStr {
		t.Errorf("output does not match expected file %s\nExpected:\n%s\nGot:\n%s",
			expectedOutputPath, expectedStr, outputStr)
	}
}

// TestCLICommands runs a table of test cases against our CLI.
func TestCLICommands(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedExitCode int
		expectedOutFile  string // Path to file containing expected output
	}{
		{
			name: "Validate config",
			args: []string{
				"-validateConfig",
				"-configIn=./tests/simple/config.yaml"},
			expectedExitCode: 0,
			expectedOutFile:  "", // Empty string means we don't check output
		},
		{
			name: "Validate query",
			args: []string{
				"-validateQuery",
				"-queryIn=./tests/simple/queries/type_equals/query.yaml"},
			expectedExitCode: 0,
			expectedOutFile:  "",
		},
		{
			name: "Execute query",
			args: []string{
				"-executeQuery",
				"-configIn=./tests/simple/config.yaml",
				"-queryIn=./tests/simple/queries/type_equals/query.yaml"},
			expectedExitCode: 0,
			expectedOutFile:  "./tests/simple/queries/type_equals/config.yaml",
		},
		{
			name: "Validate mermaid",
			args: []string{
				"-validateMermaid",
				"-configIn=./tests/simple/config.yaml",
				"-mermaidIn=./tests/simple/mermaid.yaml"},
			expectedExitCode: 0,
			expectedOutFile:  "",
		},
		{
			name: "Generate mermaid",
			args: []string{
				"-generateMermaid",
				"-configIn=./tests/simple/config.yaml",
				"-mermaidIn=./tests/simple/mermaid.yaml"},
			expectedExitCode: 0,
			expectedOutFile:  "./tests/simple/mermaid.mmd",
		},
		{
			name: "Multiple commands error",
			args: []string{
				"-validateConfig",
				"-validateQuery"},
			expectedExitCode: 1,
			expectedOutFile:  "",
		},
	}

	// For each test case, run the binary as a subprocess.
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], tc.args...)
			cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")

			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out

			err := cmd.Run()

			exitCode := 0
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				} else {
					t.Fatalf("failed to run command: %v", err)
				}
			}

			if exitCode != tc.expectedExitCode {
				t.Errorf("expected exit code %d, got %d", tc.expectedExitCode, exitCode)
			}

			// If an expected output file is specified, compare the output with its contents
			if tc.expectedOutFile != "" {
				compareOutputWithFile(t, out.String(), tc.expectedOutFile)
			}
		})
	}
}
