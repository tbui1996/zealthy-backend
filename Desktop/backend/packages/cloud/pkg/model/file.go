package model

import "time"

type File struct {
	ID               int       `json:"id"`
	FileID           string    `json:"fileId"`
	FileName         string    `json:"fileName"`
	FileMimetype     string    `json:"fileMimetype"`
	SendUserID       string    `json:"sendUserId"`
	ChatID           string    `json:"chatId"`
	DateUploaded     time.Time `json:"dateUploaded"`
	DateLastAccessed time.Time `json:"dateLastAccessed"`
	FilePath         string    `json:"filePath"`
	// a null member_id in the database means that
	// the file has not been associated
	MemberID  string     `json:"memberId"`
	DeletedAt *time.Time `json:"DeletedAt"`
}
