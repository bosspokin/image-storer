package model

import "mime/multipart"

type File struct {
	Filename string         `json:"filename" form:"filename"`
	File     multipart.File `json:"file" form:"file"`
	URL      string         `json:"url"`
}
