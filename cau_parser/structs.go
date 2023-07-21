package cau_parser

import "time"

type CAUAttachment struct {
	Name string
	Url  string
}

type CAUArticle struct {
	Url     string
	Title   string
	Author  string
	Date    time.Time
	Content string
	Files   []CAUAttachment
}
