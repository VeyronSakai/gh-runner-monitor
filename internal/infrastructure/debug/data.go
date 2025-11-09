package debug

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

// Data represents the structure of debug JSON data
type Data struct {
	CurrentTime time.Time        `json:"CurrentTime"`
	Runners     []*entity.Runner `json:"runners"`
	Jobs        []*entity.Job    `json:"jobs"`
}

// LoadDebugData loads debug data from a JSON file
// This is a helper function for creating all debug repositories
func LoadDebugData(jsonPath string) (*Data, error) {
	file, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var data Data
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &data, nil
}
