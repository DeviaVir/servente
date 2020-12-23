package mysql

import (
	"errors"

	"github.com/DeviaVir/servente/pkg/models"

	"gorm.io/gorm"
)

type OrganizationModel struct {
	DB *gorm.DB
}

func (m *OrganizationModel) Insert(identifier, name string) error {
	organization := models.Organization{Identifier: identifier, Name: name, Active: true}

	if err := m.DB.Create(&organization).Error; err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			return models.ErrDuplicateIdentifier
		}
		return err
	}

	return nil
}

func (m *OrganizationModel) Get(id int) (*models.Organization, error) {
	organization := models.Organization{}

	if err := m.DB.Where("id = ? AND active = ?", id, 1).Select("identifier", "name", "active", "id", "created_at").First(&organization).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return &organization, nil
}
