package service

import (
	"context"
	"errors"

	"parkpatrol/internal/domain"
)

type RemoteGateway interface {
	Execute(context.Context, string) (string, error)
}

type FixedRemoteGateway struct{}

func (FixedRemoteGateway) Execute(ctx context.Context, stepID string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return "verified:" + stepID, nil
	}
}

type RemoteTask struct {
	ID      string
	Context context.Context
}

type RemoteSteps struct {
	gateway RemoteGateway
}

func NewRemoteSteps(gateway RemoteGateway) *RemoteSteps {
	return &RemoteSteps{gateway: gateway}
}

func (r *RemoteSteps) Run(ctx context.Context, stepID string) (domain.StepResult, error) {
	task := RemoteTask{ID: stepID, Context: ctx}
	detail, err := r.gateway.Execute(task.Context, task.ID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return domain.StepResult{StepID: stepID, Status: domain.StepTimedOut, Detail: err.Error()}, err
		}
		return domain.StepResult{StepID: stepID, Status: domain.StepFailed, Detail: err.Error()}, err
	}
	return domain.StepResult{StepID: stepID, Status: domain.StepCompleted, Detail: detail}, nil
}
