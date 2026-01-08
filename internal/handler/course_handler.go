package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"ramah-disabilitas-be/internal/service"
	"ramah-disabilitas-be/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateCourse(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 1. Handle File Upload (Thumbnail)
	var thumbnailURL string
	file, err := c.FormFile("thumbnail")
	if err == nil {
		// Validasi ekstensi
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Validasi input gagal.",
				"errors":  map[string]string{"thumbnail": "Format file harus jpg, jpeg, atau png."},
			})
			return
		}

		uploadDir := "storage/public"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat direktori storage"})
			return
		}

		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := filepath.Join(uploadDir, filename)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file thumbnail"})
			return
		}

		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		thumbnailURL = fmt.Sprintf("%s://%s/storage/public/%s", scheme, c.Request.Host, filename)
	}

	// 2. Handle Text Fields
	title := c.PostForm("title")
	description := c.PostForm("description")
	classCode := c.PostForm("class_code")
	modulesStr := c.PostForm("modules") // JSON String

	// 3. Manual Validation
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validasi input gagal.",
			"errors":  map[string]string{"title": "Judul wajib diisi."},
		})
		return
	}

	// 4. Parse Modules JSON
	var modules []service.ModuleInput
	if modulesStr != "" {
		if err := json.Unmarshal([]byte(modulesStr), &modules); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Validasi input gagal.",
				"errors":  map[string]string{"modules": "Format JSON modules tidak valid. Pastikan format array JSON benar."},
			})
			return
		}
	}

	input := service.CourseInput{
		Title:       title,
		Description: description,
		Thumbnail:   thumbnailURL,
		ClassCode:   classCode,
		Modules:     modules,
	}

	course, err := service.CreateCourse(input, userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kelas berhasil dibuat",
		"data":    course,
	})
}

func GetMyCourses(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	courses, err := service.GetCoursesByTeacher(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daftar kelas berhasil diambil",
		"data":    courses,
	})
}

func UpdateCourse(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	courseIDStr := c.Param("id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kelas tidak valid"})
		return
	}

	var input service.CourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validasi input gagal.",
			"errors":  utils.FormatValidationError(err),
		})
		return
	}

	course, err := service.UpdateCourse(courseID, input, userID.(uint64))
	if err != nil {
		if err.Error() == "unauthorized: you do not own this course" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kelas berhasil diperbarui",
		"data":    course,
	})
}

func DeleteCourse(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	courseIDStr := c.Param("id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kelas tidak valid"})
		return
	}

	err = service.DeleteCourse(courseID, userID.(uint64))
	if err != nil {
		if err.Error() == "unauthorized: you do not own this course" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kelas berhasil dihapus",
	})
}

func JoinCourse(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		ClassCode string `json:"class_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validasi input gagal",
			"errors":  utils.FormatValidationError(err),
		})
		return
	}

	err := service.JoinCourse(input.ClassCode, userID.(uint64))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "kelas tidak ditemukan" || err.Error() == "anda sudah bergabung di kelas ini" || err.Error() == "anda adalah pengajar di kelas ini" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil bergabung ke kelas",
	})
}

func GetMyJoinedCourses(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	courses, err := service.GetStudentCourses(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daftar kelas yang diikuti berhasil diambil",
		"data":    courses,
	})
}
