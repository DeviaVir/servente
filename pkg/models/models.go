package models

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoRecord            = errors.New("models: no matching record found")
	ErrInvalidCredentials  = errors.New("models: invalid credentials")
	ErrDuplicateEmail      = errors.New("models: duplicate email")
	ErrDuplicateIdentifier = errors.New("models: duplicate service identifier")
	ErrNoOrg               = errors.New("no organization selected")

	StatusInDevelopment  = 1 // active dev
	StatusBeta           = 2 // testing
	StatusProduction     = 3 // prod and active dev
	StatusMaintain       = 4 // prod no longer developed/extended
	StatusDeprecated     = 5 // legacy, to be replaced
	StatusDecommissioned = 6 // it used to be a thing but it's gone now

	AttributesTypes = []string{"link", "webhook", "select"}
	SettingsTypes   = []string{"team-provider"}
)

// @NOTE: add new models to web/main.go AutoMigrate to have them automatically created

// Organization model definition for organizations
type Organization struct {
	gorm.Model
	Identifier             string `gorm:"type:varchar(100) unique"`
	Name                   string
	Active                 bool
	Users                  []*User `gorm:"many2many:user_organizations;"`
	Services               []*Service
	Settings               []*Setting
	ServiceAttributes      []*ServiceAttribute
	OrganizationAttributes []*OrganizationAttribute
	AuditLogs              []*AuditLog
}

// Service model definition of a service
type Service struct {
	gorm.Model
	Identifier        string `gorm:"type:varchar(100)"`
	Title             string
	Description       string
	Status            int
	ServiceAttributes []*ServiceAttribute `gorm:"foreignKey:ServiceID"`
	Organization      Organization
	OrganizationID    uint
}

// User model definition of a user
type User struct {
	gorm.Model
	Name           string
	Email          string `gorm:"unique"`
	HashedPassword []byte `gorm:"type:char(60)"`
	Active         bool
	Organizations  []*Organization `gorm:"many2many:user_organizations;"`
}

// ServiceAttribute model definition of attributes connected to services
type ServiceAttribute struct {
	gorm.Model
	Value          string
	Active         bool
	Setting        Setting
	SettingID      uint
	Organization   *Organization
	OrganizationID uint
	Service        *Service
	ServiceID      int
}

// OrganizationAttribute model definition of attributes specifically for organizations
type OrganizationAttribute struct {
	gorm.Model
	Value          string
	Active         bool
	Setting        Setting
	SettingID      uint
	Organization   *Organization
	OrganizationID uint
}

// Setting model definition of configured settings
type Setting struct {
	gorm.Model
	Key            string // unique id
	Title          string // nice name
	Type           string // link, webhook, select, etc
	Scope          string // "organization" or "service"
	Organization   *Organization
	OrganizationID uint
}

// AuditLog model definition of logs
type AuditLog struct {
	gorm.Model
	User           int
	Message        string
	Organization   *Organization
	OrganizationID uint
}
