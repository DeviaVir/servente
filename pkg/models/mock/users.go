package mock

import (
	"github.com/DeviaVir/servente/pkg/models"
)

var mockUser = &models.User{
	Model:  *mockGorm,
	Name:   "Alice",
	Email:  "alice@example.com",
	Active: true,
}

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	switch email {
	case "alice@example.com":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}

func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *UserModel) GetByEmail(email string) (*models.User, error) {
	switch email {
	case "alice@example.com":
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *UserModel) ChangePassword(id int, currentPassword, newPassword string) error {
	switch id {
	case 1:
		return nil
	default:
		return models.ErrNoRecord
	}
}

func (m *UserModel) Organizations(user *models.User) (orgs []*models.Organization, err error) {
	orgs = append(orgs, mockOrganization)
	switch user.ID {
	case 1:
		return orgs, nil
	default:
		return nil, nil
	}
}
