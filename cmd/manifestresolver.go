package main

type manifestResolver interface {
	getManifest() (*manifest, error)
}
