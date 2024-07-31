package main

import (
	"embed"

	"github.com/ctroller/goffold/cmd"
)

//go:embed templates
var DefaultTemplates embed.FS

func main() {
	cmd.DefaultTemplates = DefaultTemplates
	cmd.Execute()
}
