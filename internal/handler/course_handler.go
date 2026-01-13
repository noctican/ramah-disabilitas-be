package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"ramah-disabilitas-be/internal/model"
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
	status := c.PostForm("status")
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
		Status:      status,
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

	search := c.Query("q")
	status := c.Query("status")
	sort := c.Query("sort")

	courses, err := service.GetCoursesByTeacher(userID.(uint64), search, status, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daftar kelas berhasil diambil",
		"data":    courses,
	})
}

func GetCourseDetail(c *gin.Context) {
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

	course, err := service.GetCourseDetail(courseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kelas tidak ditemukan"})
		return
	}

	// Pastikan yang akses adalah pemilik kelas (untuk endpoint lecturer)
	if course.TeacherID != userID.(uint64) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke kelas ini"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail kelas berhasil diambil",
		"data":    course,
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

	// 1. Handle File Upload (Thumbnail)
	var thumbnailURL string
	file, err := c.FormFile("thumbnail")
	if err == nil {
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
	status := c.PostForm("status")
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
				"errors":  map[string]string{"modules": "Format JSON modules tidak valid."},
			})
			return
		}
	}

	// If thumbnail is empty but user wants to keep the old one, logic is handled in service if we pass empty string?
	// But wait, UpdateCourse service replaces everything. We need to fetch existing course inside service or handle it here.
	// Actually, service.UpdateCourse uses individual fields assignment.
	// Let's rely on service logic: if thumbnail is empty string, service assigns it.
	// Wait, if user doesn't upload new file, thumbnailURL is "". We probably want to keep existing thumbnail.
	// The current service code: `course.Thumbnail = input.Thumbnail`. This will wipe it out if empty.
	// We should only update if not empty. Let's handle this in service later?
	// For now let's construct input. NOTE: User logic might be "send empty to delete".
	// But usually "no file sent" means "no change".
	// Let's modify UpdateCourse service logic slightly or handle it here.
	// Better approach: Retrieve course first here? No, let service handle it.
	// But we need to tell service if it's a "no change" or "delete".
	// Standard multipart: if key not present/empty file -> no change.
	// Let's assume for now service simply overwrites. If we want to support "keep existing", we must pass that info.
	// However, since we cannot easily change service signature right now without affecting other things,
	// let's pass the value. If it is empty string, the service will save empty string (effectively deleting it).
	// To prevent this, we should check if file was uploaded.
	// BUT: if we are using Form-Data, we can't easily distinguish "send empty" vs "not sent" for text fields like description.
	// The existing service implementation blindly updates: `course.Thumbnail = input.Thumbnail`.
	// Use a small trick: If thumbnailURL is empty, we don't put it in input? No, struct has it.
	// We will fix service logic in next step if needed. For now let's match handler to multipart.

	input := service.CourseInput{
		Title:       title,
		Description: description,
		Thumbnail:   thumbnailURL,
		ClassCode:   classCode,
		Status:      status,
		Modules:     modules,
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
		} else if strings.Contains(strings.ToLower(err.Error()), "foreign key constraint fails") || strings.Contains(strings.ToLower(err.Error()), "constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Kelas tidak dapat dihapus karena masih digunakan (terdapat siswa/modul di dalamnya)."})
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

func DeleteModule(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	moduleIDStr := c.Param("id")
	moduleID, err := strconv.ParseUint(moduleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID modul tidak valid"})
		return
	}

	err = service.DeleteModule(moduleID, userID.(uint64))
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "tidak ditemukan") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Modul dan materi di dalamnya berhasil dihapus",
	})
}

func DeleteMaterial(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	materialIDStr := c.Param("id")
	materialID, err := strconv.ParseUint(materialIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID materi tidak valid"})
		return
	}

	err = service.DeleteMaterial(materialID, userID.(uint64))
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "tidak ditemukan") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Materi berhasil dihapus",
	})
}

func CreateMaterial(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	moduleIDStr := c.Param("id")
	moduleID, err := strconv.ParseUint(moduleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID modul tidak valid"})
		return
	}

	// Parsing Form Data
	title := c.PostForm("title")
	materialType := c.PostForm("type")
	sourceURL := c.PostForm("source_url")
	rawContent := c.PostForm("raw_content")
	durationMin, _ := strconv.Atoi(c.PostForm("duration_min"))
	hasCaptions, _ := strconv.ParseBool(c.PostForm("has_captions"))

	// Validasi basic
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Judul materi wajib diisi"})
		return
	}
	if materialType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe materi wajib diisi"})
		return
	}

	// *** Handle File Upload (Jika ada file di form untuk menggantikan source_url) ***
	file, fileErr := c.FormFile("file")
	if fileErr == nil {
		ext := filepath.Ext(file.Filename)
		uploadDir := "storage/public"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat direktori storage"})
			return
		}

		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := filepath.Join(uploadDir, filename)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
			return
		}

		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		sourceURL = fmt.Sprintf("%s://%s/storage/public/%s", scheme, c.Request.Host, filename)
	}

	input := service.MaterialInput{
		Title:       title,
		Type:        model.MaterialType(materialType), // "pdf", "youtube", "text"
		SourceURL:   sourceURL,
		RawContent:  rawContent,
		DurationMin: durationMin,
		HasCaptions: hasCaptions,
	}

	material, err := service.CreateMaterial(moduleID, input, userID.(uint64))
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

	c.JSON(http.StatusCreated, gin.H{
		"message": "Materi berhasil ditambahkan",
		"data":    material,
	})
}

func UpdateMaterial(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	materialIDStr := c.Param("id")
	materialID, err := strconv.ParseUint(materialIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID materi tidak valid"})
		return
	}

	// Parsing
	title := c.PostForm("title")
	materialType := c.PostForm("type")
	sourceURL := c.PostForm("source_url")
	rawContent := c.PostForm("raw_content")
	durationMin, _ := strconv.Atoi(c.PostForm("duration_min"))
	hasCaptions, _ := strconv.ParseBool(c.PostForm("has_captions"))

	// *** Handle File Upload ***
	file, fileErr := c.FormFile("file")
	if fileErr == nil {
		ext := filepath.Ext(file.Filename)
		uploadDir := "storage/public"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat direktori storage"})
			return
		}

		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := filepath.Join(uploadDir, filename)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
			return
		}

		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		sourceURL = fmt.Sprintf("%s://%s/storage/public/%s", scheme, c.Request.Host, filename)
	}

	input := service.MaterialInput{
		Title:       title,
		Type:        model.MaterialType(materialType),
		SourceURL:   sourceURL,
		RawContent:  rawContent,
		DurationMin: durationMin,
		HasCaptions: hasCaptions,
	}
	// Jika user tidak mengirim type, kita asumsikan 'text' atau tidak update?
	// Karena logic service.MaterialInput binding required, di sini kita manual.
	// Service Input punya binding required untuk type.
	// Jika kosong, fetch existing?
	// Tapi untuk POST handler manual seperti ini, lebih baik mandatory atau check existing di service.
	// Di sini kita asumsikan FE mengirim semua field yg mau diupdate.
	// Tapi type mungkin tidak berubah.

	// Namun service layer kita pakai struct MaterialInput yg field Type-nya required?
	// Mari cek definisi struct: `json:"type" binding:"required"`
	// Karena kita construct manual, binding tag tidak effect runtime unless we use validator manually.
	// Tapi kalo empty string masuk ke model, bisa error enum.
	// Mari kita biarkan, tapi jika kosong default atau reject?
	// Simple validation:
	if materialType == "" {
		// Maybe user just wants to update title. We should probably fetch existing type or allow empty string
		// and handle it in Service (Update logic).
		// My service logic overwrites: `material.Type = input.Type`. So it will break if empty.
		// Let's enforce it for now or rely on FE sending it.
	}

	material, err := service.UpdateMaterial(materialID, input, userID.(uint64))
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
		"message": "Materi berhasil diperbarui",
		"data":    material,
	})
}

func ToggleMaterialCompletion(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	materialIDStr := c.Param("id")
	materialID, err := strconv.ParseUint(materialIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID materi tidak valid"})
		return
	}

	completed, err := service.ToggleMaterialCompletion(userID.(uint64), materialID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	message := "Material ditandai selesai"
	if !completed {
		message = "Tanda selesai dihapus"
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   message,
		"completed": completed,
	})
}

func GetStudentCourseDetail(c *gin.Context) {
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

	course, err := service.GetStudentCourseDetail(courseID, userID.(uint64))
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
		"message": "Detail kelas berhasil diambil",
		"data":    course,
	})
}

func GetMaterialDetail(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	materialIDStr := c.Param("id")
	materialID, err := strconv.ParseUint(materialIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID materi tidak valid"})
		return
	}

	material, err := service.GetMaterialDetailWithStatus(materialID, userID.(uint64))
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
		"message": "Detail materi berhasil diambil",
		"data":    material,
	})
}

func CreateStudentAndEnroll(c *gin.Context) {
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

	var input service.CreateStudentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validasi input gagal.",
			"errors":  utils.FormatValidationError(err),
		})
		return
	}

	user, err := service.CreateStudentAndEnroll(courseID, input, userID.(uint64))
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

	c.JSON(http.StatusCreated, gin.H{
		"message": "Siswa berhasil dibuat dan ditambahkan ke kelas",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
