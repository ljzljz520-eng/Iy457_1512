package service

import (
	"slices"
	"strings"

	"parkpatrol/internal/domain"
	"parkpatrol/internal/store"
)

type Patrol struct {
	store *store.Memory
}

func NewPatrol(memory *store.Memory) *Patrol {
	return &Patrol{store: memory}
}

func (p *Patrol) CreateRoute(actor domain.Actor, route domain.Route) error {
	if err := permit(actor, domain.RoleCaptain); err != nil {
		return err
	}
	if route.ID == "" || strings.TrimSpace(route.Name) == "" || len(route.Checkpoints) == 0 {
		return domain.ErrInvalidInput
	}
	if route.CampusID == "" {
		route.CampusID = actor.CampusID
	}
	if route.CampusID != actor.CampusID {
		return domain.ErrForbidden
	}
	return p.store.PutRoute(route)
}

func (p *Patrol) CreateShift(actor domain.Actor, shift domain.Shift) error {
	if err := permit(actor, domain.RoleCaptain); err != nil {
		return err
	}
	if shift.ID == "" || shift.OfficerID == "" || shift.ServiceDay == "" || shift.Label == "" {
		return domain.ErrInvalidInput
	}
	if shift.CampusID == "" {
		shift.CampusID = actor.CampusID
	}
	if shift.CampusID != actor.CampusID {
		return domain.ErrForbidden
	}
	if _, err := p.store.Route(actor.CampusID, shift.RouteID); err != nil {
		return err
	}
	return p.store.PutShift(shift)
}

func (p *Patrol) Start(actor domain.Actor, shiftID, recordID string) (domain.PatrolRecord, error) {
	if err := permit(actor, domain.RoleOfficer); err != nil {
		return domain.PatrolRecord{}, err
	}
	if recordID == "" {
		return domain.PatrolRecord{}, domain.ErrInvalidInput
	}
	shift, err := p.store.Shift(actor.CampusID, shiftID)
	if err != nil {
		return domain.PatrolRecord{}, err
	}
	if shift.OfficerID != actor.ID {
		return domain.PatrolRecord{}, domain.ErrForbidden
	}
	record := domain.PatrolRecord{
		ID:        recordID,
		CampusID:  actor.CampusID,
		ShiftID:   shift.ID,
		OfficerID: actor.ID,
		Status:    domain.RecordActive,
		CheckIns:  []domain.CheckIn{},
	}
	if err := p.store.PutRecord(record); err != nil {
		return domain.PatrolRecord{}, err
	}
	return record, nil
}

func (p *Patrol) CheckIn(actor domain.Actor, recordID, checkpointID string) (domain.PatrolRecord, error) {
	if err := permit(actor, domain.RoleOfficer); err != nil {
		return domain.PatrolRecord{}, err
	}
	record, err := p.store.Record(actor.CampusID, recordID)
	if err != nil {
		return domain.PatrolRecord{}, err
	}
	if record.OfficerID != actor.ID {
		return domain.PatrolRecord{}, domain.ErrForbidden
	}
	if record.Status != domain.RecordActive {
		return domain.PatrolRecord{}, domain.ErrInvalidState
	}
	shift, err := p.store.Shift(actor.CampusID, record.ShiftID)
	if err != nil {
		return domain.PatrolRecord{}, err
	}
	route, err := p.store.Route(actor.CampusID, shift.RouteID)
	if err != nil {
		return domain.PatrolRecord{}, err
	}
	next := len(record.CheckIns)
	if next >= len(route.Checkpoints) || route.Checkpoints[next] != checkpointID {
		return domain.PatrolRecord{}, domain.ErrInvalidState
	}
	if slices.ContainsFunc(record.CheckIns, func(item domain.CheckIn) bool { return item.CheckpointID == checkpointID }) {
		return domain.PatrolRecord{}, domain.ErrInvalidState
	}
	record.CheckIns = append(record.CheckIns, domain.CheckIn{CheckpointID: checkpointID, Sequence: next + 1})
	if len(record.CheckIns) == len(route.Checkpoints) {
		record.Status = domain.RecordComplete
	}
	if err := p.store.UpdateRecord(record); err != nil {
		return domain.PatrolRecord{}, err
	}
	return record, nil
}

func (p *Patrol) ReportIncident(actor domain.Actor, recordID string, incident domain.Incident) (domain.Incident, error) {
	if err := permit(actor, domain.RoleOfficer); err != nil {
		return domain.Incident{}, err
	}
	record, err := p.store.Record(actor.CampusID, recordID)
	if err != nil {
		return domain.Incident{}, err
	}
	if record.OfficerID != actor.ID {
		return domain.Incident{}, domain.ErrForbidden
	}
	if incident.ID == "" || incident.Description == "" || incident.Severity == "" {
		return domain.Incident{}, domain.ErrInvalidInput
	}
	incident.CampusID = actor.CampusID
	incident.RecordID = record.ID
	incident.ReporterID = actor.ID
	incident.Status = domain.IncidentReported
	if err := p.store.PutIncident(incident); err != nil {
		return domain.Incident{}, err
	}
	return incident, nil
}

func (p *Patrol) ResolveIncident(actor domain.Actor, incidentID, resolution string) (domain.Incident, error) {
	if err := permit(actor, domain.RoleCaptain); err != nil {
		return domain.Incident{}, err
	}
	incident, err := p.store.Incident(actor.CampusID, incidentID)
	if err != nil {
		return domain.Incident{}, err
	}
	if incident.Status != domain.IncidentReported || strings.TrimSpace(resolution) == "" {
		return domain.Incident{}, domain.ErrInvalidState
	}
	incident.Resolution = resolution
	incident.Status = domain.IncidentResolved
	if err := p.store.UpdateIncident(incident); err != nil {
		return domain.Incident{}, err
	}
	return incident, nil
}

func (p *Patrol) ReviewIncident(actor domain.Actor, incidentID, note string, approved bool) (domain.Incident, error) {
	if err := permit(actor, domain.RoleSupervisor); err != nil {
		return domain.Incident{}, err
	}
	incident, err := p.store.Incident(actor.CampusID, incidentID)
	if err != nil {
		return domain.Incident{}, err
	}
	if incident.Status != domain.IncidentResolved || strings.TrimSpace(note) == "" {
		return domain.Incident{}, domain.ErrInvalidState
	}
	incident.ReviewNote = note
	if approved {
		incident.Status = domain.IncidentApproved
	} else {
		incident.Status = domain.IncidentRejected
	}
	if err := p.store.UpdateIncident(incident); err != nil {
		return domain.Incident{}, err
	}
	return incident, nil
}

func (p *Patrol) Routes(actor domain.Actor) []domain.Route {
	return p.store.Routes(actor.CampusID)
}

func (p *Patrol) Records(actor domain.Actor) []domain.PatrolRecord {
	return p.store.Records(actor.CampusID)
}

func permit(actor domain.Actor, roles ...domain.Role) error {
	if actor.ID == "" || actor.CampusID == "" || !slices.Contains(roles, actor.Role) {
		return domain.ErrForbidden
	}
	return nil
}
