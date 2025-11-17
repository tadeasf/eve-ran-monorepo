package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type IntArray []int

// Value implements the driver.Valuer interface for database serialization
func (a IntArray) Value() (driver.Value, error) {
	if a == nil {
		return json.Marshal([]int{})
	}
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface for database deserialization
func (a *IntArray) Scan(value interface{}) error {
	if value == nil {
		*a = IntArray([]int{})
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var arr []int
	err := json.Unmarshal(bytes, &arr)
	*a = IntArray(arr)
	return err
}

// Position model
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func (p *Position) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	err := json.Unmarshal(bytes, &p)
	return err
}

// PaginatedResponse model
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Limit      int         `json:"limit"`
	Offset     int         `json:"offset"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalItems int         `json:"totalItems"`
	TotalPages int         `json:"totalPages"`
}

// ErrorResponse model
type ErrorResponse struct {
	Error string `json:"error"`
}
