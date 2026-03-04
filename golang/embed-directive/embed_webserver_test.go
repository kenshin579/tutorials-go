package main

import (
	"embed"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed static/*
var staticFiles embed.FS

func TestEmbed_FileServer_ServesHTML(t *testing.T) {
	subFS, err := fs.Sub(staticFiles, "static")
	assert.NoError(t, err)

	server := httptest.NewServer(http.FileServer(http.FS(subFS)))
	defer server.Close()

	resp, err := http.Get(server.URL + "/index.html")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "text/html")
}

func TestEmbed_FileServer_ServesCSS(t *testing.T) {
	subFS, err := fs.Sub(staticFiles, "static")
	assert.NoError(t, err)

	server := httptest.NewServer(http.FileServer(http.FS(subFS)))
	defer server.Close()

	resp, err := http.Get(server.URL + "/style.css")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "text/css")
}

func TestEmbed_FileServer_Returns404ForMissing(t *testing.T) {
	subFS, err := fs.Sub(staticFiles, "static")
	assert.NoError(t, err)

	server := httptest.NewServer(http.FileServer(http.FS(subFS)))
	defer server.Close()

	resp, err := http.Get(server.URL + "/notfound.html")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
