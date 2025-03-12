package main

import (
	"embed"
	"fmt"
	"testing"

	_ "embed"
)

//go:embed hello.txt
var s string

//go:embed hello.txt
var b []byte

//go:embed hello.txt
var f embed.FS

//go:embed files/*
var files embed.FS

func TestEmbed_AsString(t *testing.T) {
	fmt.Println(s)
}

func TestEmbed_AsByte(t *testing.T) {
	fmt.Println(string(b))
}

func TestEmbed_AsFile(t *testing.T) {
	data, _ := f.ReadFile("hello.txt")
	fmt.Println(string(data))
}

func TestEmbed_AsDir(t *testing.T) {
	file1, _ := files.ReadFile("files/file1.txt")
	fmt.Println(string(file1))

	file2, _ := files.ReadFile("files/file2.txt")
	fmt.Println(string(file2))
}
