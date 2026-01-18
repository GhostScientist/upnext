package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"upnext/internal/model"
)

// RenderPlain outputs the todo list in plain text format
func RenderPlain(data *model.Data) string {
	if len(data.Items) == 0 {
		return "No tasks. Add one with: upnext add \"your task\""
	}

	var lines []string
	lines = append(lines, "Tasks:")
	lines = append(lines, strings.Repeat("-", 50))

	for i, item := range data.Items {
		pri := prioritySymbol(item.Priority)
		lines = append(lines, fmt.Sprintf("%d. [%s] %s", i+1, pri, item.Text))
		if item.Description != "" {
			lines = append(lines, fmt.Sprintf("      %s", item.Description))
		}
	}

	lines = append(lines, strings.Repeat("-", 50))
	lines = append(lines, fmt.Sprintf("%d items | %d completed total", len(data.Items), data.Stats.TotalCompleted))

	return strings.Join(lines, "\n")
}

func prioritySymbol(p model.Priority) string {
	switch p {
	case model.PriorityHigh:
		return "!!!"
	case model.PriorityMedium:
		return "!!"
	default:
		return "!"
	}
}

// RenderJSON outputs the todo list in JSON format
func RenderJSON(data *model.Data) (string, error) {
	output := struct {
		Items []struct {
			ID          string `json:"id"`
			Text        string `json:"text"`
			Description string `json:"description,omitempty"`
			Priority    string `json:"priority"`
			Created     string `json:"created"`
			Position    int    `json:"position"`
		} `json:"items"`
		Stats struct {
			TotalCompleted int `json:"total_completed"`
			StreakDays     int `json:"streak_days"`
		} `json:"stats"`
	}{
		Stats: struct {
			TotalCompleted int `json:"total_completed"`
			StreakDays     int `json:"streak_days"`
		}{
			TotalCompleted: data.Stats.TotalCompleted,
			StreakDays:     data.Stats.StreakDays,
		},
	}

	for _, item := range data.Items {
		output.Items = append(output.Items, struct {
			ID          string `json:"id"`
			Text        string `json:"text"`
			Description string `json:"description,omitempty"`
			Priority    string `json:"priority"`
			Created     string `json:"created"`
			Position    int    `json:"position"`
		}{
			ID:          item.ID,
			Text:        item.Text,
			Description: item.Description,
			Priority:    item.Priority.String(),
			Created:     item.Created.Format("2006-01-02T15:04:05Z07:00"),
			Position:    item.Position,
		})
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
