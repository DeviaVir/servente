package mock

import (
	"github.com/DeviaVir/servente/pkg/models"
)

var mockOrganization = &models.Organization{
	Model:      *mockGorm,
	Name:       "Alice",
	Identifier: "alice",
	Active:     true,
}

type OrganizationModel struct{}

func (m *OrganizationModel) Insert(identifier, name string) error {
	switch identifier {
	case "dupe":
		return models.ErrDuplicateIdentifier
	default:
		return nil
	}
}

func (m *OrganizationModel) Get(id int) (*models.Organization, error) {
	switch id {
	case 1:
		return mockOrganization, nil
	default:
		return nil, models.ErrNoRecord
	}
}
