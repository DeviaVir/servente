package mysql

import (
	"errors"

	"github.com/DeviaVir/servente/pkg/models"
	"gorm.io/gorm"
)

type ServiceModel struct {
	DB *gorm.DB
}

func (m *ServiceModel) Insert(identifier, title, description, attributes string, status int) (int, error) {
	service := models.Service{
		Identifier:  identifier,
		Title:       title,
		Description: description,
		Attributes:  attributes,
		Status:      status,
	}

	if err := m.DB.Create(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			return 0, models.ErrDuplicateIdentifier
		}
		return 0, err
	}

	return int(service.ID), nil
}

func (m *ServiceModel) Get(id int) (*models.Service, error) {
	service := models.Service{}

	if err := m.DB.Where("id = ?", id, true).First(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return &service, nil
}

func (m *ServiceModel) Latest(limit int) ([]*models.Service, error) {
	services := []*models.Service{}

	if err := m.DB.Limit(limit).Order("created_at desc").Find(&services).Error; err != nil {
		return nil, err
	}

	return services, nil
}
