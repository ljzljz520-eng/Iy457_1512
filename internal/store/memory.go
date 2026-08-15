package store

import (
	"sort"
	"sync"

	"parkpatrol/internal/domain"
)

type Memory struct {
	mu        sync.RWMutex
	routes    map[string]domain.Route
	shifts    map[string]domain.Shift
	records   map[string]domain.PatrolRecord
	incidents map[string]domain.Incident
}

func NewMemory() *Memory {
	return &Memory{
		routes:    make(map[string]domain.Route),
		shifts:    make(map[string]domain.Shift),
		records:   make(map[string]domain.PatrolRecord),
		incidents: make(map[string]domain.Incident),
	}
}

func (m *Memory) PutRoute(route domain.Route) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.routes[route.ID]; exists {
		return domain.ErrInvalidInput
	}
	route.Checkpoints = append([]string(nil), route.Checkpoints...)
	m.routes[route.ID] = route
	return nil
}

func (m *Memory) Route(campusID, id string) (domain.Route, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	route, ok := m.routes[id]
	if !ok || route.CampusID != campusID {
		return domain.Route{}, domain.ErrNotFound
	}
	route.Checkpoints = append([]string(nil), route.Checkpoints...)
	return route, nil
}

func (m *Memory) Routes(campusID string) []domain.Route {
	m.mu.RLock()
	defer m.mu.RUnlock()
	routes := make([]domain.Route, 0)
	for _, route := range m.routes {
		if route.CampusID == campusID {
			route.Checkpoints = append([]string(nil), route.Checkpoints...)
			routes = append(routes, route)
		}
	}
	sort.Slice(routes, func(i, j int) bool { return routes[i].ID < routes[j].ID })
	return routes
}

func (m *Memory) PutShift(shift domain.Shift) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.shifts[shift.ID]; exists {
		return domain.ErrInvalidInput
	}
	m.shifts[shift.ID] = shift
	return nil
}

func (m *Memory) Shift(campusID, id string) (domain.Shift, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	shift, ok := m.shifts[id]
	if !ok || shift.CampusID != campusID {
		return domain.Shift{}, domain.ErrNotFound
	}
	return shift, nil
}

func (m *Memory) PutRecord(record domain.PatrolRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.records[record.ID]; exists {
		return domain.ErrInvalidInput
	}
	record.CheckIns = append([]domain.CheckIn(nil), record.CheckIns...)
	m.records[record.ID] = record
	return nil
}

func (m *Memory) Record(campusID, id string) (domain.PatrolRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	record, ok := m.records[id]
	if !ok || record.CampusID != campusID {
		return domain.PatrolRecord{}, domain.ErrNotFound
	}
	record.CheckIns = append([]domain.CheckIn(nil), record.CheckIns...)
	return record, nil
}

func (m *Memory) UpdateRecord(record domain.PatrolRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	current, ok := m.records[record.ID]
	if !ok || current.CampusID != record.CampusID {
		return domain.ErrNotFound
	}
	record.CheckIns = append([]domain.CheckIn(nil), record.CheckIns...)
	m.records[record.ID] = record
	return nil
}

func (m *Memory) Records(campusID string) []domain.PatrolRecord {
	m.mu.RLock()
	defer m.mu.RUnlock()
	records := make([]domain.PatrolRecord, 0)
	for _, record := range m.records {
		if record.CampusID == campusID {
			record.CheckIns = append([]domain.CheckIn(nil), record.CheckIns...)
			records = append(records, record)
		}
	}
	sort.Slice(records, func(i, j int) bool { return records[i].ID < records[j].ID })
	return records
}

func (m *Memory) PutIncident(incident domain.Incident) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.incidents[incident.ID]; exists {
		return domain.ErrInvalidInput
	}
	m.incidents[incident.ID] = incident
	return nil
}

func (m *Memory) Incident(campusID, id string) (domain.Incident, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	incident, ok := m.incidents[id]
	if !ok || incident.CampusID != campusID {
		return domain.Incident{}, domain.ErrNotFound
	}
	return incident, nil
}

func (m *Memory) UpdateIncident(incident domain.Incident) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	current, ok := m.incidents[incident.ID]
	if !ok || current.CampusID != incident.CampusID {
		return domain.ErrNotFound
	}
	m.incidents[incident.ID] = incident
	return nil
}

func (m *Memory) Incidents(campusID string) []domain.Incident {
	m.mu.RLock()
	defer m.mu.RUnlock()
	incidents := make([]domain.Incident, 0)
	for _, incident := range m.incidents {
		if incident.CampusID == campusID {
			incidents = append(incidents, incident)
		}
	}
	sort.Slice(incidents, func(i, j int) bool { return incidents[i].ID < incidents[j].ID })
	return incidents
}
