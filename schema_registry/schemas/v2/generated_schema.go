package v2

// Code generated by avro/gen. DO NOT EDIT.

import (
	"time"
)

// Task is a generated struct.
type Task struct {
	TaskID      string    `avro:"task_id"`
	UserID      string    `avro:"user_id"`
	Description string    `avro:"description"`
	JiraID      string    `avro:"jira_id"`
	Completed   bool      `avro:"completed"`
	CreatedAt   time.Time `avro:"created_at"`
}

// TaskAssigned is a generated struct.
type TaskAssigned struct {
	Task      Task    `avro:"task"`
	OldUserID *string `avro:"old_user_id"`
}

// TaskCompleted is a generated struct.
type TaskCompleted struct {
	Task Task `avro:"task"`
}