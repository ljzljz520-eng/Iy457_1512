package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"parkpatrol/internal/app"
	"parkpatrol/internal/domain"
	"parkpatrol/internal/fixture"
)

type demoOutput struct {
	Campus         string                `json:"campus"`
	Route          string                `json:"route"`
	RecordStatus   domain.RecordStatus   `json:"recordStatus"`
	IncidentStatus domain.IncidentStatus `json:"incidentStatus"`
	MenuCodes      []string              `json:"menuCodes"`
	RoleCount      int                   `json:"roleCount"`
	SeverityCount  int                   `json:"severityCount"`
	RemoteStatus   domain.StepStatus     `json:"remoteStatus"`
	CSV            string                `json:"csv"`
}

func main() {
	if err := run(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(output *os.File) error {
	application := app.New()
	route := fixture.NorthRoute()
	if err := application.Patrol.CreateRoute(fixture.NorthCaptain, route); err != nil {
		return err
	}
	shift := fixture.NorthShift()
	if err := application.Patrol.CreateShift(fixture.NorthCaptain, shift); err != nil {
		return err
	}
	record, err := application.Patrol.Start(fixture.NorthOfficer, shift.ID, "record-north-001")
	if err != nil {
		return err
	}
	incident, err := application.Patrol.ReportIncident(fixture.NorthOfficer, record.ID, domain.Incident{
		ID: "incident-north-001", Severity: "high", Description: "Utility cabinet seal is broken",
	})
	if err != nil {
		return err
	}
	for _, checkpoint := range route.Checkpoints {
		record, err = application.Patrol.CheckIn(fixture.NorthOfficer, record.ID, checkpoint)
		if err != nil {
			return err
		}
	}
	incident, err = application.Patrol.ResolveIncident(fixture.NorthCaptain, incident.ID, "Cabinet isolated and seal replaced")
	if err != nil {
		return err
	}
	incident, err = application.Patrol.ReviewIncident(fixture.NorthSupervisor, incident.ID, "Evidence confirmed", true)
	if err != nil {
		return err
	}
	remoteResult, err := application.Remote.Run(context.Background(), "remote-camera-check")
	if err != nil {
		return err
	}
	severity, err := application.Config.Dictionary("incident_severity")
	if err != nil {
		return err
	}
	menus := application.Config.Menus(fixture.NorthSupervisor)
	menuCodes := make([]string, len(menus))
	for i, menu := range menus {
		menuCodes[i] = menu.Code
	}
	var csvOutput bytes.Buffer
	if err := application.Export.RecordsCSV(fixture.NorthSupervisor, &csvOutput); err != nil {
		return err
	}
	result := demoOutput{
		Campus: fixture.NorthSupervisor.CampusID, Route: route.ID, RecordStatus: record.Status,
		IncidentStatus: incident.Status, MenuCodes: menuCodes, RoleCount: len(application.Config.Roles()),
		SeverityCount: len(severity), RemoteStatus: remoteResult.Status, CSV: csvOutput.String(),
	}
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}
