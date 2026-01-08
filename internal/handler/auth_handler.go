package handler

import (
	"net/http"
	"ramah-disabilitas-be/internal/service"
	"ramah-disabilitas-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var input service.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validasi input gagal.",
			"errors":  utils.FormatValidationError(err),
		})
		return
	}

	user, err := service.Register(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pengguna berhasil terdaftar",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func Login(c *gin.Context) {
	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validasi input gagal.",
			"errors":  utils.FormatValidationError(err),
		})
		return
	}

	user, token, err := service.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
		"user": gin.H{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"role":   user.Role,
			"avatar": user.Avatar,
		},
	})
}

func GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := service.GetMe(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User detail",
		"data": gin.H{
			"id":            user.ID,
			"name":          user.Name,
			"email":         user.Email,
			"role":          user.Role,
			"avatar":        user.Avatar,
			"accessibility": user.Accessibility,
		},
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout berhasil",
	})
}
