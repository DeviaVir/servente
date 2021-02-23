package mysql

import (
	"errors"

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

func (m *OrganizationModel) Update(user *models.User, organization *models.Organization, name string) (*models.Organization, error) {
	if err := m.DB.Save(&organization).Error; err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			return nil, models.ErrDuplicateIdentifier
		}
		return nil, err
	}

	return organization, nil
}

func (m *OrganizationModel) UpdateSetting(updateSetting *models.Setting) (*models.Setting, error) {
	setting := models.Setting{}
	found := true
	if err := m.DB.Where("`key` = ? AND `scope` = ? AND `organization_id` = ?", updateSetting.Key, updateSetting.Scope, updateSetting.OrganizationID).First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			found = false
		} else {
			return nil, err
		}
	}

	if found {
		setting.Type = updateSetting.Type
		setting.Scope = updateSetting.Scope
		setting.Title = updateSetting.Title
		if err := m.DB.Save(&setting).Error; err != nil {
			return nil, err
		}
	} else {
		setting = *updateSetting
		if err := m.DB.Create(&setting).Error; err != nil {
			return nil, err
		}
	}

	return &setting, nil
}

func (m *OrganizationModel) UpdateAttribute(setting *models.Setting, val string) (*models.OrganizationAttribute, error) {
	attr := models.OrganizationAttribute{}

	found := true
	if err := m.DB.Where("setting_id = ? AND organization_id = ?", setting.ID, setting.OrganizationID).First(&attr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			found = false
			attr = models.OrganizationAttribute{
				Active:         true,
				Value:          val,
				SettingID:      setting.ID,
				OrganizationID: setting.OrganizationID,
			}
		} else {
			return nil, err
		}
	}

	if found {
		attr.Value = val
		if err := m.DB.Save(&attr).Error; err != nil {
			return nil, err
		}
	} else {
		if err := m.DB.Create(&attr).Error; err != nil {
			return nil, err
		}
	}

	return &attr, nil
}

func (m *OrganizationModel) Get(id string) (*models.Organization, error) {
	organization := models.Organization{}

	if err := m.DB.Where("identifier = ? AND active = ?", id, 1).First(&organization).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return &organization, nil
}

func (m *OrganizationModel) GetSettings(organization *models.Organization) ([]*models.Setting, error) {
	if err := m.DB.Model(&organization).Association("Settings").Error; err != nil {
		return nil, err
	}

	var settings []*models.Setting
	m.DB.Model(&organization).Association("Settings").Find(&settings)

	return settings, nil
}

func (m *OrganizationModel) GetAttributes(organization *models.Organization) ([]*models.OrganizationAttribute, error) {
	if err := m.DB.Model(&organization).Association("OrganizationAttributes").Error; err != nil {
		return nil, err
	}

	var settings []*models.OrganizationAttribute
	m.DB.Model(&organization).Association("OrganizationAttributes").Find(&settings)

	return settings, nil
}

func (m *OrganizationModel) GetServices(organization *models.Organization, limit, start int) ([]*models.Service, error) {
	if err := m.DB.Model(&organization).Association("Services").Error; err != nil {
		return nil, err
	}

	var services []*models.Service
	m.DB.Model(&organization).Offset(start).Limit(limit).Order("created_at desc").Association("Services").Find(&services)

	return services, nil
}

func (m *OrganizationModel) GetService(organization *models.Organization, id int) (*models.Service, error) {
	if err := m.DB.Model(&organization).Association("Services").Error; err != nil {
		return nil, err
	}

	var service *models.Service
	m.DB.Model(&organization).First(&service, id).Association("Services")

	return service, nil
}

func (m *OrganizationModel) GetServiceAttributes(organization *models.Organization) ([]*models.ServiceAttribute, error) {
	if err := m.DB.Model(&organization).Association("ServiceAttributes").Error; err != nil {
		return nil, err
	}

	var settings []*models.ServiceAttribute
	m.DB.Model(&organization).Association("ServiceAttributes").Find(&settings)

	return settings, nil
}
