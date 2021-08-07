package main

import "io"

type Reader struct {
	File io.Reader
	Path string
}

type Writer struct {
	File io.Writer
	Path string
}
