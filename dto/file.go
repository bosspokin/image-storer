package dto

import "mime/multipart"

type File struct {
	ID       uint
	Filename string         `json:"filename" form:"filename"`
	URL      string         `json:"url"`
	File     multipart.File `json:"file,omitempty" form:"file"`
}

type RenameFile struct {
	New string `json:"new"`
}
