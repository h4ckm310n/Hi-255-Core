package models

type FileInfo struct {
	Filename string
	FileTime int64
	Size     int64
}

type FileItem struct {
	ID       int64
	Filename string
	Sender   string
	Size     int64
	Sendtime int64
}
