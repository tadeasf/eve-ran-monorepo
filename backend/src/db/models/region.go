package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Region model
type Region struct {
	RegionID       int             `gorm:"primaryKey" json:"region_id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Constellations json.RawMessage `gorm:"type:jsonb" json:"constellations"`
}

func (r *Region) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := &struct {
		RegionID       int             `json:"region_id"`
		Name           string          `json:"name"`
		Description    string          `json:"description"`
		Constellations json.RawMessage `json:"constellations"`
	}{}

	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	r.RegionID = result.RegionID
	r.Name = result.Name
	r.Description = result.Description
	r.Constellations = result.Constellations

	return nil
}
