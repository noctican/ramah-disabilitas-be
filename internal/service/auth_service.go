package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/internal/repository"
	"ramah-disabilitas-be/pkg/utils"
)

type RegisterInput struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	Role            string `json:"role" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(input RegisterInput) (*model.User, error) {
	// Check if email already exists
	_, err := repository.FindUserByEmail(input.Email)
	if err == nil {
		return nil, errors.New("email sudah ada")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	var userRole model.UserRole
	switch input.Role {
	case "dosen", "lecturer":
		userRole = model.RoleLecturer
	case "mahasiswa", "student":
		userRole = model.RoleStudent
	default:
		return nil, errors.New("role tidak valid (pilih 'dosen' atau 'mahasiswa')")
	}

	// Generate verification token
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	user := &model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Role:     userRole,
		// IsVerified:        false,
		IsVerified:        true,
		VerificationToken: token,
	}

	if err := repository.CreateUser(user); err != nil {
		return nil, err
	}

	// Send verification email
	// We run this in a goroutine so it doesn't block the response,
	// OR we run it sync. Given I want to fail if it fails?
	// The prompt implies "given in response... told to check email".
	// If email fails sending, telling them to check it is wrong.
	// So I'll do it sync.
	// if err := utils.SendVerificationEmail(user.Email, token); err != nil {
	// 	// Just log error?
	// 	// For now simple.
	// }

	return user, nil
}

func Login(input LoginInput) (*model.User, string, error) {
	user, err := repository.FindUserByEmail(input.Email)
	if err != nil {
		return nil, "", errors.New("email atau password salah")
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return nil, "", errors.New("email atau password salah")
	}

	// if !user.IsVerified {
	// 	return nil, "", errors.New("email belum diverifikasi. silahkan cek inbox anda")
	// }

	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func VerifyEmail(token string) error {
	user, err := repository.FindUserByVerificationToken(token)
	if err != nil {
		return errors.New("token verifikasi tidak valid")
	}

	if user.IsVerified {
		return errors.New("email sudah terverifikasi")
	}

	user.IsVerified = true
	user.VerificationToken = ""

	return repository.UpdateUser(user)
}

func GetMe(userID uint64) (*model.User, error) {
	return repository.FindUserByID(userID)
}
