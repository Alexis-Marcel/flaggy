package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alexis/flaggy/internal/models"
)

func (s *SQLiteStore) CreateFlag(flag *models.Flag) error {
	now := time.Now().UTC()
	flag.CreatedAt = now
	flag.UpdatedAt = now

	_, err := s.db.Exec(
		`INSERT INTO flags (key, type, description, enabled, default_value, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		flag.Key, flag.Type, flag.Description, flag.Enabled,
		string(flag.DefaultValue), flag.CreatedAt, flag.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create flag: %w", err)
	}
	return nil
}

func (s *SQLiteStore) GetFlag(key string) (*models.Flag, error) {
	flag := &models.Flag{}
	var defaultVal string
	err := s.db.QueryRow(
		`SELECT key, type, description, enabled, default_value, created_at, updated_at
		 FROM flags WHERE key = ?`, key,
	).Scan(&flag.Key, &flag.Type, &flag.Description, &flag.Enabled,
		&defaultVal, &flag.CreatedAt, &flag.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get flag: %w", err)
	}
	flag.DefaultValue = json.RawMessage(defaultVal)

	rules, err := s.getRulesForFlag(key)
	if err != nil {
		return nil, err
	}
	flag.Rules = rules
	return flag, nil
}

func (s *SQLiteStore) ListFlags() ([]models.Flag, error) {
	rows, err := s.db.Query(
		`SELECT key, type, description, enabled, default_value, created_at, updated_at
		 FROM flags ORDER BY key`)
	if err != nil {
		return nil, fmt.Errorf("list flags: %w", err)
	}
	defer rows.Close()

	var flags []models.Flag
	for rows.Next() {
		var f models.Flag
		var defaultVal string
		if err := rows.Scan(&f.Key, &f.Type, &f.Description, &f.Enabled,
			&defaultVal, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan flag: %w", err)
		}
		f.DefaultValue = json.RawMessage(defaultVal)
		flags = append(flags, f)
	}
	return flags, rows.Err()
}

func (s *SQLiteStore) UpdateFlag(key string, req *models.UpdateFlagRequest) (*models.Flag, error) {
	flag, err := s.GetFlag(key)
	if err != nil {
		return nil, err
	}
	if flag == nil {
		return nil, nil
	}

	if req.Description != nil {
		flag.Description = *req.Description
	}
	if req.Enabled != nil {
		flag.Enabled = *req.Enabled
	}
	if req.DefaultValue != nil {
		if err := models.ValidateValueForType(flag.Type, req.DefaultValue); err != nil {
			return nil, fmt.Errorf("default_value: %w", err)
		}
		flag.DefaultValue = req.DefaultValue
	}
	flag.UpdatedAt = time.Now().UTC()

	_, err = s.db.Exec(
		`UPDATE flags SET description = ?, enabled = ?, default_value = ?, updated_at = ?
		 WHERE key = ?`,
		flag.Description, flag.Enabled, string(flag.DefaultValue), flag.UpdatedAt, key,
	)
	if err != nil {
		return nil, fmt.Errorf("update flag: %w", err)
	}
	return flag, nil
}

func (s *SQLiteStore) DeleteFlag(key string) error {
	res, err := s.db.Exec(`DELETE FROM flags WHERE key = ?`, key)
	if err != nil {
		return fmt.Errorf("delete flag: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("flag not found")
	}
	return nil
}

func (s *SQLiteStore) ToggleFlag(key string) (*models.Flag, error) {
	now := time.Now().UTC()
	_, err := s.db.Exec(
		`UPDATE flags SET enabled = NOT enabled, updated_at = ? WHERE key = ?`, now, key,
	)
	if err != nil {
		return nil, fmt.Errorf("toggle flag: %w", err)
	}
	return s.GetFlag(key)
}

// --- Rules ---

func (s *SQLiteStore) CreateRule(flagKey string, rule *models.Rule) error {
	now := time.Now().UTC()
	rule.FlagKey = flagKey
	rule.CreatedAt = now
	rule.UpdatedAt = now

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		`INSERT INTO rules (flag_key, description, value, priority, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		rule.FlagKey, rule.Description, string(rule.Value), rule.Priority,
		rule.CreatedAt, rule.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert rule: %w", err)
	}
	ruleID, _ := res.LastInsertId()
	rule.ID = ruleID

	for i := range rule.Conditions {
		c := &rule.Conditions[i]
		c.RuleID = ruleID
		c.CreatedAt = now
		res, err := tx.Exec(
			`INSERT INTO conditions (rule_id, attribute, operator, value, created_at)
			 VALUES (?, ?, ?, ?, ?)`,
			c.RuleID, c.Attribute, c.Operator, string(c.Value), c.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("insert condition: %w", err)
		}
		cID, _ := res.LastInsertId()
		c.ID = cID
	}

	return tx.Commit()
}

func (s *SQLiteStore) UpdateRule(flagKey string, ruleID int64, req *models.CreateRuleRequest) (*models.Rule, error) {
	now := time.Now().UTC()

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// Update the rule
	res, err := tx.Exec(
		`UPDATE rules SET description = ?, value = ?, priority = ?, updated_at = ?
		 WHERE id = ? AND flag_key = ?`,
		req.Description, string(req.Value), req.Priority, now, ruleID, flagKey,
	)
	if err != nil {
		return nil, fmt.Errorf("update rule: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, fmt.Errorf("rule not found")
	}

	// Replace conditions: delete old, insert new
	if _, err := tx.Exec(`DELETE FROM conditions WHERE rule_id = ?`, ruleID); err != nil {
		return nil, fmt.Errorf("delete conditions: %w", err)
	}

	rule := &models.Rule{
		ID:          ruleID,
		FlagKey:     flagKey,
		Description: req.Description,
		Value:       req.Value,
		Priority:    req.Priority,
		UpdatedAt:   now,
	}

	for _, c := range req.Conditions {
		cRes, err := tx.Exec(
			`INSERT INTO conditions (rule_id, attribute, operator, value, created_at)
			 VALUES (?, ?, ?, ?, ?)`,
			ruleID, c.Attribute, c.Operator, string(c.Value), now,
		)
		if err != nil {
			return nil, fmt.Errorf("insert condition: %w", err)
		}
		cID, _ := cRes.LastInsertId()
		rule.Conditions = append(rule.Conditions, models.Condition{
			ID:        cID,
			RuleID:    ruleID,
			Attribute: c.Attribute,
			Operator:  c.Operator,
			Value:     c.Value,
			CreatedAt: now,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	return rule, nil
}

func (s *SQLiteStore) DeleteRule(flagKey string, ruleID int64) error {
	res, err := s.db.Exec(
		`DELETE FROM rules WHERE id = ? AND flag_key = ?`, ruleID, flagKey,
	)
	if err != nil {
		return fmt.Errorf("delete rule: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("rule not found")
	}
	return nil
}

// --- Evaluation ---

func (s *SQLiteStore) GetFlagForEvaluation(key string) (*models.Flag, error) {
	// Query 1: Get the flag
	flag := &models.Flag{}
	var defaultVal string
	err := s.db.QueryRow(
		`SELECT key, type, description, enabled, default_value, created_at, updated_at
		 FROM flags WHERE key = ?`, key,
	).Scan(&flag.Key, &flag.Type, &flag.Description, &flag.Enabled,
		&defaultVal, &flag.CreatedAt, &flag.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get flag: %w", err)
	}
	flag.DefaultValue = json.RawMessage(defaultVal)

	// Query 2: Get rules + conditions via JOIN
	rules, err := s.getRulesForFlag(key)
	if err != nil {
		return nil, err
	}
	flag.Rules = rules
	return flag, nil
}

func (s *SQLiteStore) getRulesForFlag(flagKey string) ([]models.Rule, error) {
	rows, err := s.db.Query(
		`SELECT r.id, r.flag_key, r.description, r.value, r.priority, r.created_at, r.updated_at,
		        c.id, c.rule_id, c.attribute, c.operator, c.value, c.created_at
		 FROM rules r
		 LEFT JOIN conditions c ON c.rule_id = r.id
		 WHERE r.flag_key = ?
		 ORDER BY r.priority, r.id, c.id`, flagKey,
	)
	if err != nil {
		return nil, fmt.Errorf("get rules: %w", err)
	}
	defer rows.Close()

	ruleMap := map[int64]*models.Rule{}
	var ruleOrder []int64

	for rows.Next() {
		var r models.Rule
		var ruleVal string
		var cID sql.NullInt64
		var cRuleID sql.NullInt64
		var cAttr, cOp, cVal sql.NullString
		var cCreated sql.NullTime

		if err := rows.Scan(
			&r.ID, &r.FlagKey, &r.Description, &ruleVal, &r.Priority,
			&r.CreatedAt, &r.UpdatedAt,
			&cID, &cRuleID, &cAttr, &cOp, &cVal, &cCreated,
		); err != nil {
			return nil, fmt.Errorf("scan rule: %w", err)
		}
		r.Value = json.RawMessage(ruleVal)

		if _, exists := ruleMap[r.ID]; !exists {
			ruleMap[r.ID] = &r
			ruleOrder = append(ruleOrder, r.ID)
		}

		if cID.Valid {
			cond := models.Condition{
				ID:        cID.Int64,
				RuleID:    cRuleID.Int64,
				Attribute: cAttr.String,
				Operator:  models.Operator(cOp.String),
				Value:     json.RawMessage(cVal.String),
				CreatedAt: cCreated.Time,
			}
			ruleMap[r.ID].Conditions = append(ruleMap[r.ID].Conditions, cond)
		}
	}

	rules := make([]models.Rule, 0, len(ruleOrder))
	for _, id := range ruleOrder {
		rules = append(rules, *ruleMap[id])
	}
	return rules, rows.Err()
}
