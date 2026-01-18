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
	plainFlag bool
	jsonFlag  bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "upnext",
		Short: "A beautiful terminal todo app",
		Long:  "upnext - A minimal, beautiful interactive TUI todo app",
		RunE:  runRoot,
	}

	rootCmd.Flags().BoolVar(&plainFlag, "plain", false, "Output in plain text format")
	rootCmd.Flags().BoolVar(&jsonFlag, "json", false, "Output in JSON format")

	addCmd := &cobra.Command{
		Use:   "add [task]",
		Short: "Add a new task",
		Args:  cobra.ExactArgs(1),
		RunE:  runAdd,
	}

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

	// Default: Launch interactive TUI
	return tui.Run(s)
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

	// Create new todo at position 0
	todo := model.Todo{
		ID:       model.GenerateID(),
		Text:     args[0],
		Created:  time.Now(),
		Position: 0,
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

	fmt.Printf("Added: %s\n", args[0])
	return nil
}
