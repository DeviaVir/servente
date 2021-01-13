package mysql

import (
	"errors"
	"fmt"

	"github.com/DeviaVir/servente/pkg/models"

	"gorm.io/gorm"
)

type OrganizationModel struct {
	DB *gorm.DB
}

func (m *OrganizationModel) Insert(user *models.User, identifier, name string) (*models.Organization, error) {
	organization := models.Organization{
		Identifier: identifier,
		Name:       name,
		Active:     true,
	}

	fmt.Println(organization)

	if err := m.DB.Create(&organization).Error; err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			return nil, models.ErrDuplicateIdentifier
		}
		return nil, err
	}

	if err := m.DB.Model(&user).Association("Organizations").Error; err != nil {
		return nil, err
	}

	m.DB.Model(&user).Association("Organizations").Append(&organization)

	return &organization, nil
}

func (m *OrganizationModel) Get(id string) (*models.Organization, error) {
	organization := models.Organization{}

	if err := m.DB.Where("identifier = ? AND active = ?", id, 1).Select("identifier", "name", "active", "id", "created_at").First(&organization).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return &organization, nil
}
