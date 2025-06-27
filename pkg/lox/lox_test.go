package lox

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoxFeatures(t *testing.T) {
	tests := []struct {
		name    string
		testDir string
	}{
		{
			name:    "Conditional Execution",
			testDir: "conditional_execution",
		},
		{
			name:    "Logical Operators",
			testDir: "logical_operators",
		},
		{
			name:    "Loops",
			testDir: "loops",
		},
		{
			name:    "Functions",
			testDir: "functions",
		},
		{
			name:    "Resolving & Binding",
			testDir: "resolving_binding",
		},
		{
			name:    "Classes",
			testDir: "classes",
		},
		{
			name:    "Inheritence",
			testDir: "inheritence",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := os.ReadDir(filepath.Join("../../testdata", tt.testDir))
			if err != nil {
				t.Fatal(err)
			}

			for _, f := range files {
				if !strings.HasSuffix(f.Name(), ".lox") {
					continue
				}

				sourcePath := filepath.Join("../../testdata", tt.testDir, f.Name())
				expectedPath := filepath.Join("../../testdata", tt.testDir, strings.TrimSuffix(f.Name(), ".lox")+".expected")

				source, err := os.ReadFile(sourcePath)
				if err != nil {
					t.Fatal(err)
				}

				expected, err := os.ReadFile(expectedPath)
				if err != nil {
					t.Fatal(err)
				}

				// Capture stdout
				old := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w

				// Run the interpreter
				Run(string(source))

				// Restore stdout
				w.Close()
				os.Stdout = old

				// Read captured output
				var output strings.Builder
				_, err = io.Copy(&output, r)
				if err != nil {
					t.Fatal(err)
				}

				// Compare output with expected
				trimmedOutput := strings.ReplaceAll(strings.TrimSpace(output.String()), "\r\n", "\n")
				trimmedExpected := strings.ReplaceAll(strings.TrimSpace(string(expected)), "\r\n", "\n")
				if trimmedOutput != trimmedExpected {
					t.Errorf(
						"Test %s failed.\nExpected (len=%d):\n%q\nGot (len=%d):\n%q",
						f.Name(),
						len(trimmedExpected),
						trimmedExpected,
						len(trimmedOutput),
						trimmedOutput)
					// t.Errorf("Test %s failed.\nExpected:\n%s\n\nGot:\n%s",
					// 	f.Name(), trimmedExpected, trimmedOutput)
				}
			}
		})
	}
}
