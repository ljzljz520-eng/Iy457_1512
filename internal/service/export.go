package service

import (
	"encoding/csv"
	"io"
	"strconv"

	"parkpatrol/internal/domain"
	"parkpatrol/internal/store"
)

type Exporter struct {
	store *store.Memory
}

func NewExporter(memory *store.Memory) *Exporter {
	return &Exporter{store: memory}
}

func (e *Exporter) RecordsCSV(actor domain.Actor, output io.Writer) error {
	if err := permit(actor, domain.RoleCaptain, domain.RoleSupervisor); err != nil {
		return err
	}
	writer := csv.NewWriter(output)
	if err := writer.Write([]string{"record_id", "campus_id", "shift_id", "officer_id", "status", "checkin_count"}); err != nil {
		return err
	}
	for _, record := range e.store.Records(actor.CampusID) {
		if err := writer.Write([]string{
			record.ID,
			record.CampusID,
			record.ShiftID,
			record.OfficerID,
			string(record.Status),
			strconv.Itoa(len(record.CheckIns)),
		}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
