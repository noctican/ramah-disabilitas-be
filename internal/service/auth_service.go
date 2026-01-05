package service

import (
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
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(input RegisterInput) (*model.User, error) {
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Role:     model.RoleStudent,
		RankTier: model.RankBronze,
	}

	if err := repository.CreateUser(user); err != nil {
		return nil, err
	}

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

	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func GetMe(userID uint64) (*model.User, error) {
	return repository.FindUserByID(userID)
}
