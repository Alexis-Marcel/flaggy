package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

type FlagType string

const (
	FlagTypeBoolean FlagType = "boolean"
	FlagTypeString  FlagType = "string"
	FlagTypeNumber  FlagType = "number"
	FlagTypeJSON    FlagType = "json"
)

type Operator string

const (
	OpEquals     Operator = "equals"
	OpNotEquals  Operator = "not_equals"
	OpIn         Operator = "in"
	OpNotIn      Operator = "not_in"
	OpContains   Operator = "contains"
	OpStartsWith Operator = "starts_with"
	OpGT         Operator = "gt"
	OpGTE        Operator = "gte"
	OpLT         Operator = "lt"
	OpLTE        Operator = "lte"
	OpExists     Operator = "exists"
	OpRegex      Operator = "regex"
)

var validOperators = map[Operator]bool{
	OpEquals: true, OpNotEquals: true,
	OpIn: true, OpNotIn: true,
	OpContains: true, OpStartsWith: true,
	OpGT: true, OpGTE: true, OpLT: true, OpLTE: true,
	OpExists: true, OpRegex: true,
}

type Flag struct {
	Key          string          `json:"key"`
	Type         FlagType        `json:"type"`
	Description  string          `json:"description"`
	Enabled      bool            `json:"enabled"`
	DefaultValue json.RawMessage `json:"default_value"`
	Rules        []Rule          `json:"rules,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type Rule struct {
	ID          int64           `json:"id"`
	FlagKey     string          `json:"flag_key"`
	Description string          `json:"description"`
	Value       json.RawMessage `json:"value"`
	Priority    int             `json:"priority"`
	Conditions  []Condition     `json:"conditions"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type Condition struct {
	ID        int64           `json:"id"`
	RuleID    int64           `json:"rule_id"`
	Attribute string          `json:"attribute"`
	Operator  Operator        `json:"operator"`
	Value     json.RawMessage `json:"value"`
	CreatedAt time.Time       `json:"created_at"`
}

var keyRegex = regexp.MustCompile(`^[a-z][a-z0-9_]{1,62}[a-z0-9]$`)

func ValidateFlag(f *Flag) error {
	if !keyRegex.MatchString(f.Key) {
		return fmt.Errorf("key must match %s", keyRegex.String())
	}
	switch f.Type {
	case FlagTypeBoolean, FlagTypeString, FlagTypeNumber, FlagTypeJSON:
	default:
		return fmt.Errorf("invalid flag type: %q", f.Type)
	}
	if err := ValidateValueForType(f.Type, f.DefaultValue); err != nil {
		return fmt.Errorf("default_value: %w", err)
	}
	return nil
}

func ValidateRule(r *Rule) error {
	if len(r.Conditions) == 0 {
		return fmt.Errorf("rule must have at least one condition")
	}
	for i, c := range r.Conditions {
		if err := ValidateCondition(&c); err != nil {
			return fmt.Errorf("condition[%d]: %w", i, err)
		}
	}
	return nil
}

func ValidateCondition(c *Condition) error {
	if c.Attribute == "" {
		return fmt.Errorf("attribute is required")
	}
	if !validOperators[c.Operator] {
		return fmt.Errorf("invalid operator: %q", c.Operator)
	}
	return nil
}

func ValidateValueForType(ft FlagType, raw json.RawMessage) error {
	if len(raw) == 0 {
		return fmt.Errorf("value is required")
	}
	var v interface{}
	if err := json.Unmarshal(raw, &v); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	switch ft {
	case FlagTypeBoolean:
		if _, ok := v.(bool); !ok {
			return fmt.Errorf("expected boolean value")
		}
	case FlagTypeString:
		if _, ok := v.(string); !ok {
			return fmt.Errorf("expected string value")
		}
	case FlagTypeNumber:
		if _, ok := v.(float64); !ok {
			return fmt.Errorf("expected number value")
		}
	case FlagTypeJSON:
		// any valid JSON is fine
	}
	return nil
}
