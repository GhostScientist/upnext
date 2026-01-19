package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"upnext/internal/cli"
	"upnext/internal/model"
	"upnext/internal/store"
	"upnext/internal/tui"
)

var (
	plainFlag   bool
	jsonFlag    bool
	globalFlag  bool
	allFlag     bool
	priorityStr string
	descFlag    string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "upnext",
		Short: "A beautiful terminal todo app",
		Long:  "upnext - A minimal, beautiful interactive TUI todo app with context-aware task management",
		RunE:  runRoot,
	}

	rootCmd.Flags().BoolVar(&plainFlag, "plain", false, "Output in plain text format")
	rootCmd.Flags().BoolVar(&jsonFlag, "json", false, "Output in JSON format")
	rootCmd.Flags().BoolVar(&allFlag, "all", false, "Show all tasks regardless of context")

	addCmd := &cobra.Command{
		Use:   "add [task]",
		Short: "Add a new task (context-aware by default)",
		Long: `Add a new task to upnext. By default, the task is associated with your
current working directory, so it will only appear when you're in that
directory or its subdirectories.

Use --global to create a task visible from anywhere.`,
		Args: cobra.ExactArgs(1),
		RunE: runAdd,
	}

	addCmd.Flags().BoolVarP(&globalFlag, "global", "g", false, "Create a global task (visible from anywhere)")
	addCmd.Flags().StringVarP(&priorityStr, "priority", "p", "medium", "Priority: high, medium, or low")
	addCmd.Flags().StringVarP(&descFlag, "desc", "d", "", "Task description")

	rootCmd.AddCommand(addCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) error {
	s, err := store.NewJSONStore()
	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}

	// Get current working directory for context filtering
	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}

	// Static output modes
	if plainFlag || jsonFlag {
		data, err := s.Load()
		if err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		if jsonFlag {
			output, err := cli.RenderJSON(data)
			if err != nil {
				return fmt.Errorf("failed to render JSON: %w", err)
			}
			fmt.Println(output)
		} else {
			fmt.Println(cli.RenderPlain(data))
		}
		return nil
	}

	// Default: Launch interactive TUI with context
	return tui.RunWithContext(s, cwd, allFlag)
}

func runAdd(cmd *cobra.Command, args []string) error {
	s, err := store.NewJSONStore()
	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}

	data, err := s.Load()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// Get context (working directory) unless --global is set
	var context string
	if !globalFlag {
		context, err = os.Getwd()
		if err != nil {
			context = ""
		}
	}

	// Parse priority
	priority := model.PriorityMedium
	switch priorityStr {
	case "high", "h":
		priority = model.PriorityHigh
	case "low", "l":
		priority = model.PriorityLow
	}

	// Create new todo at position 0
	todo := model.Todo{
		ID:          model.GenerateID(),
		Text:        args[0],
		Description: descFlag,
		Priority:    priority,
		Created:     time.Now(),
		Position:    0,
		Context:     context,
	}

	// Shift existing items
	for i := range data.Items {
		data.Items[i].Position++
	}

	// Insert at beginning
	data.Items = append([]model.Todo{todo}, data.Items...)

	if err := s.Save(data); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	location := "here"
	if globalFlag {
		location = "globally"
	}
	fmt.Printf("Added %s: %s\n", location, args[0])
	return nil
}
