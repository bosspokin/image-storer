package model

import "mime/multipart"

type File struct {
	Filename string         `json:"filename" form:"filename"`
	URL      string         `json:"url"`
	File     multipart.File `json:"file,omitempty" form:"file"`
}
