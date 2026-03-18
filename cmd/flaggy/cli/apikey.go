package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var apikeyCmd = &cobra.Command{
	Use:     "apikey",
	Aliases: []string{"api-key"},
	Short:   "Manage API keys",
}

// --- apikey list ---

var apikeyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, status, err := doRequest("GET", "/api/v1/api-keys", nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("server error (%d): %s", status, string(data))
		}

		var keys []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Environment string `json:"environment"`
			Prefix      string `json:"prefix"`
			Revoked     bool   `json:"revoked"`
			CreatedAt   string `json:"created_at"`
		}
		if err := json.Unmarshal(data, &keys); err != nil {
			return fmt.Errorf("parse response: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tENVIRONMENT\tPREFIX\tREVOKED")
		for _, k := range keys {
			revoked := "no"
			if k.Revoked {
				revoked = "yes"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", k.ID, k.Name, k.Environment, k.Prefix, revoked)
		}
		w.Flush()
		return nil
	},
}

// --- apikey create ---

var (
	apikeyEnv string
)

var apikeyCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new API key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body := map[string]interface{}{
			"name":        args[0],
			"environment": apikeyEnv,
		}

		data, status, err := doRequest("POST", "/api/v1/api-keys", body)
		if err != nil {
			return err
		}
		if status != 201 {
			return fmt.Errorf("server error (%d): %s", status, string(data))
		}

		var result struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Environment string `json:"environment"`
			Key         string `json:"key"`
		}
		if err := json.Unmarshal(data, &result); err != nil {
			return fmt.Errorf("parse response: %w", err)
		}

		fmt.Printf("API key created:\n")
		fmt.Printf("  ID:          %s\n", result.ID)
		fmt.Printf("  Name:        %s\n", result.Name)
		fmt.Printf("  Environment: %s\n", result.Environment)
		fmt.Printf("  Key:         %s\n", result.Key)
		fmt.Printf("\nSave this key now — it won't be shown again.\n")
		return nil
	},
}

// --- apikey revoke ---

var apikeyRevokeCmd = &cobra.Command{
	Use:   "revoke <id>",
	Short: "Revoke an API key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, status, err := doRequest("DELETE", "/api/v1/api-keys/"+args[0], nil)
		if err != nil {
			return err
		}
		if status != 204 {
			return fmt.Errorf("server error (%d): %s", status, string(data))
		}
		fmt.Printf("API key %q revoked\n", args[0])
		return nil
	},
}

func init() {
	apikeyCreateCmd.Flags().StringVar(&apikeyEnv, "env", "live", "Environment (live, test, staging)")

	apikeyCmd.AddCommand(apikeyListCmd, apikeyCreateCmd, apikeyRevokeCmd)
	rootCmd.AddCommand(apikeyCmd)
}
