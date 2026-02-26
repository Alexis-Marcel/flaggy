package models

import "encoding/json"

type CreateFlagRequest struct {
	Key          string          `json:"key"`
	Type         FlagType        `json:"type"`
	Description  string          `json:"description"`
	Enabled      bool            `json:"enabled"`
	DefaultValue json.RawMessage `json:"default_value"`
}

type UpdateFlagRequest struct {
	Description  *string         `json:"description,omitempty"`
	Enabled      *bool           `json:"enabled,omitempty"`
	DefaultValue json.RawMessage `json:"default_value,omitempty"`
}

type CreateRuleRequest struct {
	Description       string          `json:"description"`
	Conditions        []Condition     `json:"conditions"`
	SegmentKeys       []string        `json:"segment_keys,omitempty"`
	Value             json.RawMessage `json:"value"`
	Priority          int             `json:"priority"`
	RolloutPercentage int             `json:"rollout_percentage"`
}

type BatchEvaluateRequest struct {
	Flags   []string               `json:"flags"`
	Context map[string]interface{} `json:"context"`
}

type BatchEvaluateResponse struct {
	Results []EvaluateResponse `json:"results"`
}

type EvaluateRequest struct {
	FlagKey string                 `json:"flag_key"`
	Context map[string]interface{} `json:"context"`
}

type EvaluateResponse struct {
	FlagKey string          `json:"flag_key"`
	Value   json.RawMessage `json:"value"`
	Match   bool            `json:"match"`
	Reason  string          `json:"reason"`
}
