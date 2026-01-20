package handler

import (
	"net/http"
	"ramah-disabilitas-be/internal/service"
	"ramah-disabilitas-be/pkg/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func UpdateAccessibility(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input service.AccessibilityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validasi input gagal.",
			"errors":  utils.FormatValidationError(err),
		})
		return
	}

	profile, err := service.UpdateAccessibilityProfile(userID.(uint64), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var actions []string
	if profile.AISummary {
		actions = append(actions, "Materi akan otomatis diringkas (Mode Fokus)")
	}
	if profile.SubtitlesRequired {
		actions = append(actions, "Video akan selalu menampilkan subtitle")
	}
	if profile.ScreenReaderCompatible {
		actions = append(actions, "Fitur pembaca layar diaktifkan")
	}
	if profile.KeyboardNavigation {
		actions = append(actions, "Navigasi keyboard diaktifkan")
	}
	if profile.TextBasedSubmission {
		actions = append(actions, "Pengumpulan tugas via teks diizinkan")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Preferensi disabilitas berhasil disimpan",
		"data": gin.H{
			"profile":          profile,
			"confirmation_msg": "Terima kasih informasinya! Berdasarkan pilihanmu, kami telah menyiapkan:",
			"active_features":  actions,
		},
	})
}

func CreateStudentByLecturer(c *gin.Context) {
	// Middleware assumed to check lecturer role
	var input service.CreateStudentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validasi input gagal.",
			"errors":  utils.FormatValidationError(err),
		})
		return
	}

	lecturerID, _ := c.Get("userID")
	teacherID := lecturerID.(uint64)
	user, err := service.CreateStudent(input, &teacherID)
	if err != nil {
		status := http.StatusInternalServerError
		// Simple duplicate check if needed, or rely on generic error
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Akun mahasiswa berhasil dibuat",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func ImportStudents(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File CSV wajib diunggah (key: 'file')"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka file"})
		return
	}
	defer f.Close()

	users, err := service.ImportStudentsFromCSV(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses CSV: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Import berhasil",
		"count":   len(users),
	})
}

func UpdateStudentByLecturer(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	studentIDStr := c.Param("id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID siswa tidak valid"})
		return
	}

	var input service.CreateStudentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validasi input gagal: " + err.Error()})
		return
	}

	user, err := service.UpdateStudentByLecturer(studentID, input, userID.(uint64))
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "unauthorized") {
			status = http.StatusForbidden
		} else if strings.Contains(err.Error(), "tidak ditemukan") {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data siswa berhasil diperbarui",
		"data":    user,
	})
}

func DeleteStudentByLecturer(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	studentIDStr := c.Param("id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID siswa tidak valid"})
		return
	}

	err = service.DeleteStudentByLecturer(studentID, userID.(uint64))
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "unauthorized") {
			status = http.StatusForbidden
		} else if strings.Contains(err.Error(), "tidak ditemukan") {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data siswa berhasil dihapus",
	})
}
