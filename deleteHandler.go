package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"secure-vault/models"
)

func deleteHandler(c *gin.Context) {
	token := c.Param("token")

	var meta models.FileMeta
	result := db.First(&meta, "token = ?", token)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid token"})
		return
	}

	if err := s3Uploader.Delete(meta.S3Key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete from S3"})
		return
	}

	db.Delete(&meta)

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
