package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"time"

	storage_go "github.com/supabase-community/storage-go"
)

func UploadToSupabase(file multipart.File, filename string, contentType string) (string, error) {
	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	bucketName := "uploads" // Pastikan bucket ini ada dan public di Supabase

	// Pastikan variable env diset
	if supabaseUrl == "" || supabaseKey == "" {
		return "", fmt.Errorf("SUPABASE_URL atau SUPABASE_KEY belum diset")
	}

	// Inisialisasi client storage
	// Format URL storage: https://<project_id>.supabase.co/storage/v1
	storageClient := storage_go.NewClient(supabaseUrl+"/storage/v1", supabaseKey, nil)

	// Buat nama file unik dengan timestamp
	uniqueFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)

	// Upload file
	_, err := storageClient.UploadFile(bucketName, uniqueFilename, file)
	if err != nil {
		return "", fmt.Errorf("gagal upload ke supabase: %v", err)
	}

	// Generate Public URL
	// Format public URL: https://<project_id>.supabase.co/storage/v1/object/public/<bucket>/<filename>
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseUrl, bucketName, uniqueFilename)

	return publicURL, nil
}
