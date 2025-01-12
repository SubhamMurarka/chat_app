package helpers

import (
	"fmt"
	"time"

	"github.com/sony/sonyflake"
)

func GenerateID() (uint64, error) {
	// Configure Sonyflake with custom settings
	settings := sonyflake.Settings{
		StartTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), // Custom epoch
		MachineID: func() (uint16, error) {
			return 42, nil // Custom machine ID (e.g., 42)
		},
	}

	// Create a new Sonyflake instance
	sf := sonyflake.NewSonyflake(settings)

	// Generate a unique ID
	id, err := sf.NextID()
	if err != nil {
		fmt.Printf("Failed to generate ID: %v\n", err)
		return 0, err
	}

	// Print the generated ID
	fmt.Printf("Generated ID: %d\n", id)

	return id, nil
}
