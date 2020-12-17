package mysql

import (
	"errors"

	"github.com/DeviaVir/servente/pkg/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserModel struct {
	DB *gorm.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	// 12 represents iterations (4096^12), also accepts a SALT to combat rainbow-table attacks
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user := models.User{Name: name, Email: email, HashedPassword: hashedPassword, Active: true}

	if err := m.DB.Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			return models.ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	user := models.User{}
	if err := m.DB.Where("email = ?", email).Where("active = ?", true).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	return int(user.ID), nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	user := models.User{}

	if err := m.DB.Where("id = ? AND active = ?", id, 1).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return &user, nil
}

// ChangePassword allows easily changing a user's password
func (m *UserModel) ChangePassword(id int, currentPassword, newPassword string) error {
	user := models.User{}

	if err := m.DB.Where("id = ?", id).Select("hashed_password").First(&user).Error; err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(currentPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return models.ErrInvalidCredentials
		}
		return err
	}

	// generate and set new password
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	if err := m.DB.Model(&user).Update("hashed_password", string(newHashedPassword)).Error; err != nil {
		return err
	}

	return nil
}
