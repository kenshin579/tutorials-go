package tasks

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeWelcomeEmail  = "email:welcome"
	TypeReminderEmail = "email:reminder"
	TypeLogging       = "logging"
)

type EmailTaskPayload struct {
	// ID for the email recipient.
	UserID int
}

type LoggingTaskPayload struct {
	ID string
}

// ----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
// ----------------------------------------------

func NewLoggingTask(id string) (*asynq.Task, error) {
	payload, err := json.Marshal(LoggingTaskPayload{ID: id})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeLogging, payload), nil
}

func NewWelcomeEmailTask(id int) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailTaskPayload{UserID: id})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeWelcomeEmail, payload), nil
}

func NewReminderEmailTask(id int) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailTaskPayload{UserID: id})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeReminderEmail, payload), nil
}

func HandleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var p EmailTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	log.Printf(" [*] Send Welcome Email to User %d", p.UserID)
	return nil
}

func HandleReminderEmailTask(ctx context.Context, t *asynq.Task) error {
	var p EmailTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	log.Printf(" [*] Send Reminder Email to User %d", p.UserID)
	return nil
}

func HandleLoggingTask(ctx context.Context, t *asynq.Task) error {
	var p LoggingTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	log.Printf(" [*] Logging Task: %s", p.ID)
	return nil
}
