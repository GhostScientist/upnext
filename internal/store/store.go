package store

import "upnext/internal/model"

// Store defines the interface for persisting todo data
type Store interface {
	// Load reads the data from storage
	Load() (*model.Data, error)
	// Save writes the data to storage
	Save(data *model.Data) error
}
