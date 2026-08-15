package fixture

import "parkpatrol/internal/domain"

var (
	NorthCaptain    = domain.Actor{ID: "captain-north", CampusID: "campus-north", Role: domain.RoleCaptain}
	NorthOfficer    = domain.Actor{ID: "officer-north", CampusID: "campus-north", Role: domain.RoleOfficer}
	NorthSupervisor = domain.Actor{ID: "supervisor-north", CampusID: "campus-north", Role: domain.RoleSupervisor}
	SouthCaptain    = domain.Actor{ID: "captain-south", CampusID: "campus-south", Role: domain.RoleCaptain}
)

func NorthRoute() domain.Route {
	return domain.Route{
		ID:          "route-north-night",
		CampusID:    "campus-north",
		Name:        "North perimeter",
		Checkpoints: []string{"north-gate", "warehouse", "utility-room"},
	}
}

func NorthShift() domain.Shift {
	return domain.Shift{
		ID:         "shift-2026-08-15-night",
		CampusID:   "campus-north",
		RouteID:    "route-north-night",
		OfficerID:  "officer-north",
		ServiceDay: "2026-08-15",
		Label:      "Night shift",
	}
}
