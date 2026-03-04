package main

import (
	"bytes"
	"embed"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed templates/*.html
var templateFiles embed.FS

func TestEmbed_Template_ParseFS(t *testing.T) {
	tmpl, err := template.ParseFS(templateFiles, "templates/*.html")
	assert.NoError(t, err)

	data := struct {
		Title   string
		Message string
	}{
		Title:   "Hello Embed",
		Message: "This is rendered from an embedded template.",
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "index.html", data)
	assert.NoError(t, err)

	result := buf.String()
	assert.Contains(t, result, "Hello Embed")
	assert.Contains(t, result, "This is rendered from an embedded template.")
}
