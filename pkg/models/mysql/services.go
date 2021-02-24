package mysql

import (
	"errors"

	"github.com/DeviaVir/servente/pkg/models"
	"gorm.io/gorm"
)

type ServiceModel struct {
	DB *gorm.DB
}

func (m *ServiceModel) Insert(org *models.Organization, identifier, title, description string, attributes []*models.ServiceAttribute, status int) (int, error) {
	service := models.Service{
		Identifier:        identifier,
		Title:             title,
		Description:       description,
		ServiceAttributes: attributes,
		Status:            status,
		OrganizationID:    org.ID,
	}

	// @TODO: service identifier already exists for this org?

	if err := m.DB.Create(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			return 0, models.ErrDuplicateIdentifier
		}
		return 0, err
	}

	return int(service.ID), nil
}

func (m *ServiceModel) Get(org *models.Organization, id int) (*models.Service, error) {
	service := models.Service{}

	if err := m.DB.Where("id = ?", id, true).Where("organization_id = ?", org.ID).First(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return &service, nil
}

func (m *ServiceModel) Latest(org *models.Organization, start, limit int) ([]*models.Service, error) {
	services := []*models.Service{}

	if err := m.DB.Offset(start).Limit(limit).Where("organization_id = ?", org.ID).Order("created_at desc").Find(&services).Error; err != nil {
		return nil, err
	}

	return services, nil
}
