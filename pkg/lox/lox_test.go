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
			name:    "if_statement",
			testDir: "if_statement",
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
				trimmedOutput := strings.TrimSpace(output.String())
				trimmedExpected := strings.TrimSpace(string(expected))
				if strings.Compare(trimmedOutput, trimmedExpected) == 0 {
					t.Errorf(
						"Test %s failed.\nExpected:\n%s\nGot:\n%s",
						f.Name(), trimmedExpected, trimmedOutput)
				}
			}
		})
	}
}
