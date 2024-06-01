package openapi

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var dist embed.FS
var DistFS, _ = fs.Sub(dist, "dist")

//go:embed openapi.yaml
var OpenAPI []byte
