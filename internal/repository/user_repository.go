package repository

import (
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/pkg/database"
)

func CreateUser(user *model.User) error {
	return database.DB.Create(user).Error
}

func FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := database.DB.Preload("Accessibility").Where("email = ?", email).First(&user).Error
	return &user, err
}

func FindUserByID(id uint64) (*model.User, error) {
	var user model.User
	err := database.DB.Preload("Accessibility").First(&user, id).Error
	return &user, err
}

func FindUserByVerificationToken(token string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("verification_token = ?", token).First(&user).Error
	return &user, err
}

func UpdateUser(user *model.User) error {
	return database.DB.Save(user).Error
}

func DeleteUser(id uint64) error {
	return database.DB.Delete(&model.User{}, id).Error
}
