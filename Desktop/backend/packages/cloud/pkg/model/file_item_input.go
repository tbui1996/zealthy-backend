package model

import "time"

type FileItemInput struct {
	FileID    string
	FileName  string
	MimeType  string
	SenderID  string
	SessionID string
	FilePath  string
	Bucket    string
	Time      time.Time
}
