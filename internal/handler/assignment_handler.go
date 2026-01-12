package handler

import (
	"net/http"
	"ramah-disabilitas-be/internal/service"
	"ramah-disabilitas-be/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateAssignment(c *gin.Context) {
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

	var input service.AssignmentInput
	// Gunakan ShouldBind agar bisa handle JSON maupun Form Data
	if err := c.ShouldBind(&input); err != nil {
		errorMessages := utils.FormatValidationError(err)

		// Jika FormatValidationError mengembalikan general error untuk EOF/Unmarshal, kita perjelas
		if err.Error() == "EOF" {
			errorMessages = map[string]string{"general": "Body request kosong. Pastikan mengirim JSON atau Form Data yang valid."}
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validasi input gagal",
			"errors":  errorMessages,
		})
		return
	}

	assignment, err := service.CreateAssignment(courseID, input, userID.(uint64))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "unauthorized: anda tidak memiliki akses ke kelas ini" {
			status = http.StatusForbidden
		} else if err.Error() == "kelas tidak ditemukan" || err.Error() == "modul tidak ditemukan" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tugas berhasil dibuat",
		"data":    assignment,
	})
}

func GetAssignments(c *gin.Context) {
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

	assignments, err := service.GetAssignmentsByCourse(courseID, userID.(uint64))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "unauthorized: anda tidak memiliki akses ke kelas ini" {
			status = http.StatusForbidden
		} else if err.Error() == "kelas tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daftar tugas berhasil diambil",
		"data":    assignments,
	})
}

func GetMyAssignments(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	assignments, err := service.GetStudentAssignments(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daftar tugas saya berhasil diambil",
		"data":    assignments,
	})
}
