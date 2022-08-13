package model

import "mime/multipart"

type File struct {
	ID       uint
	Filename string         `json:"filename" form:"filename"`
	URL      string         `json:"url"`
	File     multipart.File `json:"file,omitempty" form:"file"`
}

type RenameFile struct {
	Old string `json:"old"`
	New string `json:"new"`
}
