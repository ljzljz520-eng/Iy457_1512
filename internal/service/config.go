package service

import "parkpatrol/internal/domain"

type Configuration struct {
	menus        []domain.Menu
	roles        []domain.RoleDefinition
	dictionaries map[string][]domain.DictionaryItem
}

func NewConfiguration() *Configuration {
	return &Configuration{
		menus: []domain.Menu{
			{Code: "patrol.routes", Title: "Patrol routes", Roles: []domain.Role{domain.RoleCaptain, domain.RoleSupervisor}},
			{Code: "patrol.shifts", Title: "Patrol shifts", Roles: []domain.Role{domain.RoleCaptain, domain.RoleSupervisor}},
			{Code: "patrol.checkin", Title: "Checkpoint check-in", Roles: []domain.Role{domain.RoleOfficer}},
			{Code: "incident.review", Title: "Incident review", Roles: []domain.Role{domain.RoleSupervisor}},
			{Code: "records.export", Title: "Record export", Roles: []domain.Role{domain.RoleCaptain, domain.RoleSupervisor}},
		},
		roles: []domain.RoleDefinition{
			{Code: domain.RoleCaptain, DisplayName: "Security captain"},
			{Code: domain.RoleOfficer, DisplayName: "Patrol officer"},
			{Code: domain.RoleSupervisor, DisplayName: "Campus supervisor"},
		},
		dictionaries: map[string][]domain.DictionaryItem{
			"incident_severity": {
				{Value: "low", Label: "Low"},
				{Value: "medium", Label: "Medium"},
				{Value: "high", Label: "High"},
			},
			"incident_status": {
				{Value: string(domain.IncidentReported), Label: "Reported"},
				{Value: string(domain.IncidentResolved), Label: "Resolved"},
				{Value: string(domain.IncidentApproved), Label: "Approved"},
				{Value: string(domain.IncidentRejected), Label: "Rejected"},
			},
		},
	}
}

func (c *Configuration) Menus(actor domain.Actor) []domain.Menu {
	result := make([]domain.Menu, 0)
	for _, menu := range c.menus {
		for _, role := range menu.Roles {
			if actor.Role == role {
				copyMenu := menu
				copyMenu.Roles = append([]domain.Role(nil), menu.Roles...)
				result = append(result, copyMenu)
				break
			}
		}
	}
	return result
}

func (c *Configuration) Roles() []domain.RoleDefinition {
	return append([]domain.RoleDefinition(nil), c.roles...)
}

func (c *Configuration) Dictionary(name string) ([]domain.DictionaryItem, error) {
	items, ok := c.dictionaries[name]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return append([]domain.DictionaryItem(nil), items...), nil
}
