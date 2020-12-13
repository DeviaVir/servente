package mock

import (
	"time"

	"github.com/DeviaVir/servente/pkg/models"
)

var mockService = &models.Service{
	ID:      1,
	Title:   "Servente",
	Content: "A fake service owned by golang...",
	Created: time.Now(),
	Expires: time.Now(),
}

type ServiceModel struct{}

func (m *ServiceModel) Insert(title, content, expires string) (int, error) {
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

func (m *ServiceModel) Latest() ([]*models.Service, error) {
	return []*models.Service{mockService}, nil
}
