package model

import (
	"time"
)

// Priority levels for tasks
type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
)

func (p Priority) String() string {
	switch p {
	case PriorityHigh:
		return "High"
	case PriorityMedium:
		return "Medium"
	default:
		return "Low"
	}
}

func (p Priority) Icon() string {
	switch p {
	case PriorityHigh:
		return "!!!"
	case PriorityMedium:
		return "!!"
	default:
		return "!"
	}
}

// Todo represents an active task in the list
type Todo struct {
	ID          string    `json:"id"`
	Text        string    `json:"text"`
	Description string    `json:"description,omitempty"`
	Priority    Priority  `json:"priority"`
	Created     time.Time `json:"created"`
	Position    int       `json:"position"`
}

// ArchivedTodo represents a completed task
type ArchivedTodo struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Created   time.Time `json:"created"`
	Completed time.Time `json:"completed"`
}

// Stats tracks completion metrics
type Stats struct {
	TotalCompleted int `json:"total_completed"`
	StreakDays     int `json:"streak_days"`
}

// Data is the root structure for persisted data
type Data struct {
	Version int            `json:"version"`
	Items   []Todo         `json:"items"`
	Archive []ArchivedTodo `json:"archive"`
	Stats   Stats          `json:"stats"`
}

// NewData creates an empty data structure
func NewData() *Data {
	return &Data{
		Version: 1,
		Items:   []Todo{},
		Archive: []ArchivedTodo{},
		Stats:   Stats{},
	}
}

// GenerateID creates a unique ID based on timestamp
func GenerateID() string {
	return time.Now().Format("20060102150405.000000000")
}
