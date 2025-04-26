package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type FileMeta struct {
	gorm.Model
	FileName     string
	S3Key        string
	Token        string `gorm:"uniqueIndex"`
	EncryptionIV []byte
	ExpiresAt    time.Time
}
