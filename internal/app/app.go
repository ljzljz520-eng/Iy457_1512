package app

import (
	"parkpatrol/internal/service"
	"parkpatrol/internal/store"
)

type Application struct {
	Patrol *service.Patrol
	Config *service.Configuration
	Export *service.Exporter
	Remote *service.RemoteSteps
}

func New() *Application {
	memory := store.NewMemory()
	return &Application{
		Patrol: service.NewPatrol(memory),
		Config: service.NewConfiguration(),
		Export: service.NewExporter(memory),
		Remote: service.NewRemoteSteps(service.FixedRemoteGateway{}),
	}
}
