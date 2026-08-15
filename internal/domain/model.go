package domain

import "errors"

var (
	ErrForbidden    = errors.New("operation forbidden")
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidState = errors.New("invalid state transition")
	ErrInvalidInput = errors.New("invalid input")
)

type Role string

const (
	RoleCaptain    Role = "captain"
	RoleOfficer    Role = "officer"
	RoleSupervisor Role = "supervisor"
)

type Actor struct {
	ID       string `json:"id"`
	CampusID string `json:"campusId"`
	Role     Role   `json:"role"`
}

type Route struct {
	ID          string   `json:"id"`
	CampusID    string   `json:"campusId"`
	Name        string   `json:"name"`
	Checkpoints []string `json:"checkpoints"`
}

type Shift struct {
	ID         string `json:"id"`
	CampusID   string `json:"campusId"`
	RouteID    string `json:"routeId"`
	OfficerID  string `json:"officerId"`
	ServiceDay string `json:"serviceDay"`
	Label      string `json:"label"`
}

type CheckIn struct {
	CheckpointID string `json:"checkpointId"`
	Sequence     int    `json:"sequence"`
}

type RecordStatus string

const (
	RecordActive   RecordStatus = "active"
	RecordComplete RecordStatus = "complete"
)

type PatrolRecord struct {
	ID        string       `json:"id"`
	CampusID  string       `json:"campusId"`
	ShiftID   string       `json:"shiftId"`
	OfficerID string       `json:"officerId"`
	Status    RecordStatus `json:"status"`
	CheckIns  []CheckIn    `json:"checkIns"`
}

type IncidentStatus string

const (
	IncidentReported IncidentStatus = "reported"
	IncidentResolved IncidentStatus = "resolved"
	IncidentApproved IncidentStatus = "approved"
	IncidentRejected IncidentStatus = "rejected"
)

type Incident struct {
	ID          string         `json:"id"`
	CampusID    string         `json:"campusId"`
	RecordID    string         `json:"recordId"`
	ReporterID  string         `json:"reporterId"`
	Severity    string         `json:"severity"`
	Description string         `json:"description"`
	Resolution  string         `json:"resolution"`
	ReviewNote  string         `json:"reviewNote"`
	Status      IncidentStatus `json:"status"`
}

type Menu struct {
	Code  string `json:"code"`
	Title string `json:"title"`
	Roles []Role `json:"roles"`
}

type RoleDefinition struct {
	Code        Role   `json:"code"`
	DisplayName string `json:"displayName"`
}

type DictionaryItem struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type StepStatus string

const (
	StepCompleted StepStatus = "completed"
	StepTimedOut  StepStatus = "timed_out"
	StepFailed    StepStatus = "failed"
)

type StepResult struct {
	StepID string     `json:"stepId"`
	Status StepStatus `json:"status"`
	Detail string     `json:"detail"`
}
