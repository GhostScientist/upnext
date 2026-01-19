package model

import (
	"path/filepath"
	"strings"
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
	Context     string    `json:"context,omitempty"` // Working directory where task was created
}

// ArchivedTodo represents a completed task
type ArchivedTodo struct {
	ID          string    `json:"id"`
	Text        string    `json:"text"`
	Description string    `json:"description,omitempty"`
	Priority    Priority  `json:"priority"`
	Created     time.Time `json:"created"`
	Completed   time.Time `json:"completed"`
	Context     string    `json:"context,omitempty"` // Working directory where task was created
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

// IsContextRelevant checks if a task context is relevant to the current working directory.
// A task is relevant if:
// - The task has no context (global task)
// - The task context is a parent of cwd (task from parent dir applies to subdirs)
// - The task context is the same as cwd
// - The task context is a child of cwd (task from subdir is visible from parent)
func IsContextRelevant(taskContext, cwd string) bool {
	// Global tasks (no context) are always relevant
	if taskContext == "" {
		return true
	}

	// Clean paths for comparison
	taskContext = filepath.Clean(taskContext)
	cwd = filepath.Clean(cwd)

	// Exact match
	if taskContext == cwd {
		return true
	}

	// Task context is parent of cwd (tasks from parent dirs apply to subdirs)
	if strings.HasPrefix(cwd, taskContext+string(filepath.Separator)) {
		return true
	}

	// Task context is child of cwd (tasks from subdirs visible from parent)
	if strings.HasPrefix(taskContext, cwd+string(filepath.Separator)) {
		return true
	}

	return false
}

// GetContextDisplay returns a display-friendly version of the context path
// relative to the current working directory
func GetContextDisplay(taskContext, cwd string) string {
	if taskContext == "" {
		return "global"
	}

	taskContext = filepath.Clean(taskContext)
	cwd = filepath.Clean(cwd)

	if taskContext == cwd {
		return "."
	}

	// Try to make it relative
	rel, err := filepath.Rel(cwd, taskContext)
	if err != nil {
		return taskContext
	}

	return rel
}

// FilterByContext returns items that are relevant to the given context
func (d *Data) FilterByContext(cwd string) []Todo {
	var filtered []Todo
	for _, item := range d.Items {
		if IsContextRelevant(item.Context, cwd) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// FilterArchiveByContext returns archived items relevant to the given context
func (d *Data) FilterArchiveByContext(cwd string) []ArchivedTodo {
	var filtered []ArchivedTodo
	for _, item := range d.Archive {
		if IsContextRelevant(item.Context, cwd) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
