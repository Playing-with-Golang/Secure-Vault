package main

import (
	"github.com/joho/godotenv"
	"math/rand"
	"os"
	"secure-vault/models"
	"secure-vault/storage"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	db            *gorm.DB
	s3Uploader    *storage.S3Uploader
	tokenTTL      = time.Hour                                  // 1 hour expiration
	encryptionKey = []byte("0123456789abcdef0123456789abcdef") // Server-side only
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	db, err = gorm.Open("postgres", "host=localhost user=user-name port=5432 dbname=secure sslmode=disable password=strong-password")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.FileMeta{})

	//defer db.Close()
	db.LogMode(true)

	s3Uploader, err = storage.NewS3Uploader("playing-with-go")
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())
}

func generateToken(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {

	//go cleanupExpiredFiles(client)

	// Gin Server
	r := gin.Default()
	r.StaticFile("/", "./ui/index.html")
	r.POST("/upload", uploadHandler)
	r.GET("/download/:token", downloadHandler)
	r.DELETE("/delete/:token", deleteHandler)

	r.Run(":" + os.Getenv("PORT"))
}
