package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"secure-vault/models"
	"secure-vault/utils"
	"strconv"
	"time"
)

func uploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	expireHoursStr := c.PostForm("expire_hours")
	expireHours, err := strconv.Atoi(expireHoursStr)
	if err != nil || expireHours < 1 || expireHours > 24 {
		expireHours = 1
	}

	tempInputPath := filepath.Join(os.TempDir(), file.Filename)
	tempEncryptedPath := tempInputPath + ".enc"

	if err := c.SaveUploadedFile(file, tempInputPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	iv, err := utils.EncryptFile(tempInputPath, tempEncryptedPath, encryptionKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt file"})
		return
	}
	defer os.Remove(tempInputPath)
	defer os.Remove(tempEncryptedPath)

	s3Key := "uploads/" + file.Filename + ".enc"
	if err := s3Uploader.Upload(tempEncryptedPath, s3Key); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to S3"})
		return
	}

	token := generateToken(32)

	meta := models.FileMeta{
		FileName:     file.Filename,
		S3Key:        s3Key,
		Token:        token,
		EncryptionIV: iv,
		ExpiresAt:    time.Now().Add(time.Duration(expireHours) * time.Hour),
	}
	db.Create(&meta)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
