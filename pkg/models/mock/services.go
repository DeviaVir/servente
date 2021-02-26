package mock

import (
	"github.com/DeviaVir/servente/pkg/models"
)

var mockService = &models.Service{
	Model:       *mockGorm,
	Identifier:  "servente",
	Title:       "Servente",
	Description: "A fake service owned by golang...",
	Status:      1,
}

type ServiceModel struct{}

func (m *ServiceModel) Insert(org *models.Organization, identifier, title, description string, attributes []*models.ServiceAttribute, status int, owner string) (int, error) {
	return 2, nil
}

func (m *ServiceModel) Update(service *models.Service) (int, error) {
	return 2, nil
}

func (m *ServiceModel) Get(org *models.Organization, id int) (*models.Service, error) {
	switch id {
	case 1:
		return mockService, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *ServiceModel) Latest(org *models.Organization, start, limit int) ([]*models.Service, error) {
	return []*models.Service{mockService}, nil
}
