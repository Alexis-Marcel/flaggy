package models

import (
	"fmt"
	"time"
)

type Segment struct {
	Key         string      `json:"key"`
	Description string      `json:"description"`
	Conditions  []Condition `json:"conditions"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type CreateSegmentRequest struct {
	Key         string      `json:"key"`
	Description string      `json:"description"`
	Conditions  []Condition `json:"conditions"`
}

type UpdateSegmentRequest struct {
	Description *string     `json:"description,omitempty"`
	Conditions  []Condition `json:"conditions,omitempty"`
}

func ValidateSegment(s *Segment) error {
	if !keyRegex.MatchString(s.Key) {
		return fmt.Errorf("key must match %s", keyRegex.String())
	}
	if len(s.Conditions) == 0 {
		return fmt.Errorf("segment must have at least one condition")
	}
	for i, c := range s.Conditions {
		if err := ValidateCondition(&c); err != nil {
			return fmt.Errorf("condition[%d]: %w", i, err)
		}
	}
	return nil
}
