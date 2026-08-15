package app_test

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"strings"
	"testing"
	"time"

	"parkpatrol/internal/app"
	"parkpatrol/internal/domain"
	"parkpatrol/internal/fixture"
)

func seededApplication(t *testing.T) *app.Application {
	t.Helper()
	application := app.New()
	if err := application.Patrol.CreateRoute(fixture.NorthCaptain, fixture.NorthRoute()); err != nil {
		t.Fatal(err)
	}
	if err := application.Patrol.CreateShift(fixture.NorthCaptain, fixture.NorthShift()); err != nil {
		t.Fatal(err)
	}
	return application
}

func TestPatrolIncidentWorkflow(t *testing.T) {
	application := seededApplication(t)
	record, err := application.Patrol.Start(fixture.NorthOfficer, fixture.NorthShift().ID, "record-workflow")
	if err != nil {
		t.Fatal(err)
	}
	incident, err := application.Patrol.ReportIncident(fixture.NorthOfficer, record.ID, domain.Incident{
		ID: "incident-workflow", Severity: "high", Description: "Door controller is offline",
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, checkpoint := range fixture.NorthRoute().Checkpoints {
		record, err = application.Patrol.CheckIn(fixture.NorthOfficer, record.ID, checkpoint)
		if err != nil {
			t.Fatal(err)
		}
	}
	incident, err = application.Patrol.ResolveIncident(fixture.NorthCaptain, incident.ID, "Controller power restored")
	if err != nil {
		t.Fatal(err)
	}
	incident, err = application.Patrol.ReviewIncident(fixture.NorthSupervisor, incident.ID, "Recovery verified", true)
	if err != nil {
		t.Fatal(err)
	}
	if record.Status != domain.RecordComplete {
		t.Fatalf("record status = %q, want %q", record.Status, domain.RecordComplete)
	}
	if incident.Status != domain.IncidentApproved {
		t.Fatalf("incident status = %q, want %q", incident.Status, domain.IncidentApproved)
	}
}

func TestCampusIsolation(t *testing.T) {
	application := seededApplication(t)
	if got := application.Patrol.Routes(fixture.SouthCaptain); len(got) != 0 {
		t.Fatalf("south campus routes = %d, want 0", len(got))
	}
	shift := fixture.NorthShift()
	shift.ID = "shift-south-attempt"
	shift.CampusID = fixture.SouthCaptain.CampusID
	if err := application.Patrol.CreateShift(fixture.SouthCaptain, shift); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("cross-campus shift error = %v, want %v", err, domain.ErrNotFound)
	}
	var output bytes.Buffer
	if err := application.Export.RecordsCSV(fixture.SouthCaptain, &output); err != nil {
		t.Fatal(err)
	}
	rows, err := csv.NewReader(strings.NewReader(output.String())).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("south campus CSV rows = %d, want 1", len(rows))
	}
}

func TestMenusRolesDictionariesAndExport(t *testing.T) {
	application := seededApplication(t)
	record, err := application.Patrol.Start(fixture.NorthOfficer, fixture.NorthShift().ID, "record-export")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := application.Patrol.CheckIn(fixture.NorthOfficer, record.ID, "north-gate"); err != nil {
		t.Fatal(err)
	}
	menus := application.Config.Menus(fixture.NorthSupervisor)
	if len(menus) != 4 {
		t.Fatalf("supervisor menus = %d, want 4", len(menus))
	}
	if len(application.Config.Roles()) != 3 {
		t.Fatalf("roles = %d, want 3", len(application.Config.Roles()))
	}
	items, err := application.Config.Dictionary("incident_severity")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 || items[2].Value != "high" {
		t.Fatalf("severity dictionary = %#v", items)
	}
	var output bytes.Buffer
	if err := application.Export.RecordsCSV(fixture.NorthSupervisor, &output); err != nil {
		t.Fatal(err)
	}
	want := "record_id,campus_id,shift_id,officer_id,status,checkin_count\nrecord-export,campus-north,shift-2026-08-15-night,officer-north,active,1\n"
	if output.String() != want {
		t.Fatalf("CSV = %q, want %q", output.String(), want)
	}
}

func TestRemoteStepImmediateDeadline(t *testing.T) {
	application := app.New()
	ctx, cancel := context.WithDeadline(context.Background(), time.Time{})
	defer cancel()
	result, err := application.Remote.Run(ctx, "remote-deadline-check")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("remote error = %v, want %v; status = %q", err, context.DeadlineExceeded, result.Status)
	}
	if result.Status != domain.StepTimedOut {
		t.Fatalf("remote status = %q, want %q", result.Status, domain.StepTimedOut)
	}
}
