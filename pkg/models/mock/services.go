package mock

import (
	"github.com/DeviaVir/servente/pkg/models"
)

var mockService = &models.Service{
	Model:       *mockGorm,
	Identifier:  "servente",
	Title:       "Servente",
	Description: "A fake service owned by golang...",
	Attributes:  "",
	Status:      1,
}

type ServiceModel struct{}

func (m *ServiceModel) Insert(identifier, title, description, attributes string, status int) (int, error) {
	return 2, nil
}

func (m *ServiceModel) Get(id int) (*models.Service, error) {
	switch id {
	case 1:
		return mockService, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *ServiceModel) Latest(limit int) ([]*models.Service, error) {
	return []*models.Service{mockService}, nil
}
