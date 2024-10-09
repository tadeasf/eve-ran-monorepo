package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

type IntArray []int

func (a *IntArray) Scan(value interface{}) error {
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
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalItems int         `json:"totalItems"`
	TotalPages int         `json:"totalPages"`
}

// ErrorResponse model
type ErrorResponse struct {
	Error string `json:"error"`
}
