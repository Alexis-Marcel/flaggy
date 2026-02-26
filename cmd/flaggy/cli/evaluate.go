package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var evalContext string

var evaluateCmd = &cobra.Command{
	Use:   "evaluate <flag_key>",
	Short: "Evaluate a flag against a context",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var ctx map[string]interface{}
		if evalContext != "" {
			if err := json.Unmarshal([]byte(evalContext), &ctx); err != nil {
				return fmt.Errorf("invalid context JSON: %w", err)
			}
		}

		body := map[string]interface{}{
			"flag_key": args[0],
			"context":  ctx,
		}

		data, status, err := doRequest("POST", "/api/v1/evaluate", body)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("server error (%d): %s", status, string(data))
		}

		var resp struct {
			FlagKey string          `json:"flag_key"`
			Value   json.RawMessage `json:"value"`
			Match   bool            `json:"match"`
			Reason  string          `json:"reason"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return fmt.Errorf("parse response: %w", err)
		}

		fmt.Printf("Flag:   %s\n", resp.FlagKey)
		fmt.Printf("Value:  %s\n", string(resp.Value))
		fmt.Printf("Match:  %v\n", resp.Match)
		fmt.Printf("Reason: %s\n", resp.Reason)
		return nil
	},
}

func init() {
	evaluateCmd.Flags().StringVarP(&evalContext, "context", "c", "", `Evaluation context as JSON (e.g. '{"user":{"plan":"pro"}}')`)
	rootCmd.AddCommand(evaluateCmd)
}
