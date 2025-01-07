package auth

import (
	"auther/internal/db/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AuthenticateUser(db *gorm.DB, login, password string) (*models.User, error) {
	var user models.User
	err := db.Where("login = ?", login).First(&user).Error
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}

func CreateUser(db *gorm.DB, user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return db.Create(user).Error
}

func DeleteUser(db *gorm.DB, user *models.User) error {

	return db.Unscoped().Delete(user).Error
}

func DeleteUserByID(db *gorm.DB, userId string) error {

	return db.Unscoped().Delete(&models.User{}, userId).Error
}

func DeleteUserByLogin(db *gorm.DB, login string) error {

	return db.Exec("DELETE FROM users WHERE login = ?", login).Error
}

func IsAdminToken(token string, adminTokens []string) bool {
	for _, adminToken := range adminTokens {
		if token == adminToken {
			return true
		}
	}
	return false
}

func SaveRefreshToken(db *gorm.DB, userID uint, refreshToken string) error {
	return db.Model(&models.User{}).Where("id = ?", userID).Update("refresh_token", refreshToken).Error
}

func ValidateRefreshToken(db *gorm.DB, userID uint, refreshToken string) error {
	var user models.User
	err := db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return err
	}

	if user.RefreshToken != refreshToken {
		return errors.New("invalid refresh token")
	}

	return nil
}

func GetUserByID(db *gorm.DB, userID uint) (*models.User, error) {
	var user models.User
	err := db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
