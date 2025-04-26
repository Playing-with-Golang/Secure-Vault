package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"secure-vault/models"
	"secure-vault/utils"
	"time"
)

func downloadHandler(c *gin.Context) {
	token := c.Param("token")

	var meta models.FileMeta
	result := db.First(&meta, "token = ? AND expires_at > ? ", token, time.Now())
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid token"})
		return
	}

	if time.Now().After(meta.ExpiresAt) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Token expired"})
		return
	}

	tempDownloadPath := filepath.Join(os.TempDir(), "download.enc")
	tempDecryptedPath := filepath.Join(os.TempDir(), meta.FileName)

	if err := s3Uploader.Download(meta.S3Key, tempDownloadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download from S3"})
		return
	}
	defer os.Remove(tempDownloadPath)

	if err := utils.DecryptFile(tempDownloadPath, tempDecryptedPath, encryptionKey, meta.EncryptionIV); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt file"})
		return
	}
	defer os.Remove(tempDecryptedPath)

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", meta.FileName))
	c.Header("Content-Type", "application/octet-stream")

	file, err := os.Open(tempDecryptedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	io.Copy(c.Writer, file)
}
