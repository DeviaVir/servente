package owners

import (
	"fmt"

	sdk "github.com/DeviaVir/servente-sdk"
	"github.com/DeviaVir/servente/pkg/models"
)

// API struct containing our api data
type API struct {
	Endpoints []string
}

// Init return API objects loaded with functions and powers
func Init(org *models.Organization) (*API, error) {
	// discover the team-provider endpoints from the organization settings
	var endpoints []string
	for _, attr := range org.OrganizationAttributes {
		if attr.Setting.Scope == "organization" && attr.Setting.Type == "team-provider" {
			endpoints = append(endpoints, attr.Value)
		}
	}
	if len(endpoints) < 1 {
		return nil, fmt.Errorf("no endpoints discovered")
	}

	api := &API{
		Endpoints: endpoints,
	}

	return api, nil
}

// GetTeamsList show all teams available to pick from
func (api *API) GetTeamsList() ([]sdk.JSONTeam, error) {
	data, err := api.handle("teams/list")
	if err != nil {
		return nil, err
	}

	var teams []sdk.JSONTeam
	for _, d := range data {
		teams = append(teams, d.Teams...)
	}
	return teams, nil
}

// UserPartOfTeams show all teams an email is part of
func (api *API) UserPartOfTeams(user *models.User) ([]sdk.JSONTeam, error) {
	data, err := api.handle(fmt.Sprintf("teams/membership/%s", user.Email))
	if err != nil {
		return nil, err
	}

	var teams []sdk.JSONTeam
	for _, d := range data {
		teams = append(teams, d.Teams...)
	}
	return teams, nil
}

// UserPartOfTeam check if a user is part of a team
func (api *API) UserPartOfTeam(user *models.User, userTeam string) (bool, error) {
	teams, err := api.UserPartOfTeams(user)
	if err != nil {
		return false, err
	}

	for _, team := range teams {
		if team.Name == userTeam {
			return true, nil
		}
	}

	return false, nil
}
