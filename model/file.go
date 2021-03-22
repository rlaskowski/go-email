package model

import "time"

type File struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modification_time"`
}

func NewFile(name string) *File {
	return &File{Name: name}
}
