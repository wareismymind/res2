package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const DEFAULT_MANIFEST_NAME = "res2.yaml"

var errNoManifestFound = errors.New("no manifest found")

type fsManifestResolver struct {
	workingDir string
	fsRoot     string
}

func newFSManifestResolver(workingDir, fsRoot string) *fsManifestResolver {
	return &fsManifestResolver{
		workingDir,
		fsRoot,
	}
}

func (f *fsManifestResolver) getManifest() (*manifest, error) {

	dir := f.workingDir

	for strings.HasPrefix(dir, f.fsRoot) {

		manifestPath := filepath.Join(dir, DEFAULT_MANIFEST_NAME)

		content, err := os.ReadFile(manifestPath)
		if errors.Is(err, os.ErrNotExist) {
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
			continue
		}
		if err != nil {
			return nil, err
		}

		manifest := &manifest{}

		err = yaml.Unmarshal(content, &manifest)
		if err != nil {
			return nil, err
		}

		manifest.projectRoot = dir
		return manifest, nil
	}

	return nil, errNoManifestFound
}
