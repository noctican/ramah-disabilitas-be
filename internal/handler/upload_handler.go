package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan. Pastikan key form-data adalah 'file'"})
		return
	}

	// Buat folder storage/public jika belum ada
	uploadDir := "storage/public"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat direktori penyimpanan"})
		return
	}

	// Generate nama file unik
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join(uploadDir, filename)

	// Simpan file
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
		return
	}

	// URL untuk akses file (sesuai static route di router)
	// get scheme http or https
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	publicURL := fmt.Sprintf("%s://%s/storage/public/%s", scheme, host, filename)

	c.JSON(http.StatusOK, gin.H{
		"message": "File berhasil diupload",
		"url":     publicURL,
	})
}
