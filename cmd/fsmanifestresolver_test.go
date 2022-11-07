package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFSManifestResolver(t *testing.T) {

	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("failed to identify working directory: %v\n", err)
		return
	}

	for _, tt := range []struct {
		name     string
		dir      string
		expected *manifest
	}{
		{
			name: "manifest in working directory",
			dir:  "manifest-in-working-directory",
			expected: &manifest{
				projectRoot: "foo/bar",
				Files: map[string]string{
					"file.txt": "https://github.com/wareismymind/res2/blob/main/manifest-in-working-directory.txt",
				},
			},
		},
		{
			name: "manifest in ancestor directory",
			dir:  "manifest-in-ancestor-directory",
			expected: &manifest{
				projectRoot: "",
				Files: map[string]string{
					"file.txt": "https://github.com/wareismymind/res2/blob/main/manifest-in-ancestor-directory.txt",
				},
			},
		},
		{
			name:     "no manifest in directory tree",
			dir:      "no-manifest",
			expected: nil,
		},
	} {

		t.Run(tt.name, func(t *testing.T) {
			fsRoot := filepath.Join(wd, "testdata", "fsmanifestresolvertests", tt.dir)
			workingDir := filepath.Join(fsRoot, "foo", "bar")
			underTest := newFSManifestResolver(workingDir, fsRoot)
			actual, err := underTest.getManifest()

			if tt.expected != nil && err != nil {
				t.Errorf("unexpected error: %v\n", err)
				return
			}

			if tt.expected == nil {
				if err == nil {
					t.Error("expected non-nil err; got nil")
				}
				return
			}

			realExpectedProjectRoot := filepath.Join(fsRoot, tt.expected.projectRoot)
			if actual.projectRoot != realExpectedProjectRoot {
				t.Errorf("expected projectRoot=%s; got %s\n", realExpectedProjectRoot, actual.projectRoot)
				return
			}

			for key, value := range tt.expected.Files {
				value2, _ := actual.Files[key]
				if value2 != value {
					t.Errorf("expected [%s]=%s; got %s\n", key, value, value2)
				}
			}

			for key, value := range actual.Files {
				if _, ok := tt.expected.Files[key]; !ok {
					t.Errorf("got unexpected value [%s]=%s\n", key, value)
				}
			}
		})
	}
}
