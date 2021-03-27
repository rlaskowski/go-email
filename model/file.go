package model

import (
	"io"
	"time"
)

type File struct {
	Name    string    `json:"name"`
	Key     string    `json:"key"`
	Size    int64     `json:"size"`
	Data    []byte    `json:"-"`
	Reader  io.Reader `json:"-"`
	ModTime time.Time `json:"modification_time"`
}

func NewFile(name string) *File {
	return &File{Name: name}
}
