package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"ramah-disabilitas-be/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan. Pastikan key form-data adalah 'file'"})
		return
	}

	// Buka file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka file"})
		return
	}
	defer file.Close()

	// Cek apakah konfigurasi Supabase ada
	if os.Getenv("SUPABASE_URL") != "" && os.Getenv("SUPABASE_KEY") != "" {
		// Gunakan Supabase Storage
		publicURL, err := utils.UploadToSupabase(file, fileHeader.Filename, fileHeader.Header.Get("Content-Type"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload ke cloud: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "File berhasil diupload ke cloud",
			"url":     publicURL,
		})
	} else {
		// Fallback ke Local Storage (Hanya untuk Development Lokal Tanpa internet/credential)
		// Namun di Koyeb ini akan hilang.

		uploadDir := "storage/public"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat direktori penyimpanan"})
			return
		}

		ext := filepath.Ext(fileHeader.Filename)
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(fileHeader, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
			return
		}

		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		host := c.Request.Host
		publicURL := fmt.Sprintf("%s://%s/storage/public/%s", scheme, host, filename)

		c.JSON(http.StatusOK, gin.H{
			"message": "File berhasil diupload (Local - Peringatan: File akan hilang jika restart di cloud)",
			"url":     publicURL,
		})
	}
}
