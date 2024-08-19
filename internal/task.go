package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Situation string

const (
	SituationInProgress = "in progress"
	SituationCompleted  = "completed"
	SituationNotStarted = "not started"
)

var validSituations = map[Situation]bool{
	SituationInProgress: true,
	SituationCompleted:  true,
	SituationNotStarted: true,
}

func IsValidSituation(s Situation) bool {
	return validSituations[s]
}

type Task struct {
	ID          uuid.UUID
	Name        string
	Description string
	Situation   Situation
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTask(
	name string,
	description string,
	situation Situation,
) Task {
	return Task{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Situation:   situation,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewUpdateTask(
	name string,
	description string,
	situation Situation,
) Task {
	return Task{
		Name:        name,
		Description: description,
		Situation:   situation,
		UpdatedAt:   time.Now(),
	}
}

func (t *Task) ValidateFields() error {
	if t.Name == "" {
		return errors.New("name cannot be empty")
	}
	if len(t.Name) < 3 {
		return errors.New("name must be at least 3 characters long")
	}
	if len(t.Description) > 255 {
		return errors.New("description must have a maximum of 255 characters")
	}
	if len(t.Name) > 32 {
		return errors.New("name must have a maximum of 32 characters")
	}
	if !IsValidSituation(t.Situation) {
		return errors.New("invalid situation value")
	}
	return nil
}
