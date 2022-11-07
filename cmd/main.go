package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to identify working directory: %v", err)
		os.Exit(1)
	}
	fsRoot := getFSRoot(wd)

	manifestResolver := newFSManifestResolver(wd, fsRoot)

	manifest, err := manifestResolver.getManifest()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to read manifest: %v\n", err)
		os.Exit(1)
	}

	for fileName, url := range manifest.Files {
		if err := download(fileName, manifest.projectRoot, url); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to restore 'fileName': %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("restored %s <- %s\n", fileName, url)
	}
}

func download(fileName, rootDir, url string) error {

	if filepath.IsAbs(fileName) {
		return errors.New("cannot restore to an absoluate path")
	}

	absPath := filepath.Join(rootDir, fileName)
	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	file, err := os.Create(absPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

func getFSRoot(dir string) string {
	parent := filepath.Dir(dir)
	for parent != dir {
		dir = parent
		parent = filepath.Dir(dir)
	}
	return parent
}
