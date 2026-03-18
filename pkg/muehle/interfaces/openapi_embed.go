package interfaces

import "embed"

//go:embed openapi.yaml
var openAPISpec []byte

//go:embed swaggerui/*
var swaggerUIStatic embed.FS
