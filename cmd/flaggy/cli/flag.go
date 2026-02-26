package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var flagCmd = &cobra.Command{
	Use:   "flag",
	Short: "Manage feature flags",
}

// --- flag list ---

var flagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all flags",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, status, err := doRequest("GET", "/api/v1/flags", nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("server error (%d): %s", status, string(data))
		}

		var flags []struct {
			Key         string `json:"key"`
			Type        string `json:"type"`
			Enabled     bool   `json:"enabled"`
			Description string `json:"description"`
		}
		if err := json.Unmarshal(data, &flags); err != nil {
			return fmt.Errorf("parse response: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "KEY\tTYPE\tENABLED\tDESCRIPTION")
		for _, f := range flags {
			status := "off"
			if f.Enabled {
				status = "on"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", f.Key, f.Type, status, f.Description)
		}
		w.Flush()
		return nil
	},
}

// --- flag get ---

var flagGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a flag by key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, status, err := doRequest("GET", "/api/v1/flags/"+args[0], nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("server error (%d): %s", status, string(data))
		}
		fmt.Println(prettyJSON(data))
		return nil
	},
}

// --- flag create ---

var (
	createType        string
	createDescription string
	createEnabled     bool
	createDefault     string
)

var flagCreateCmd = &cobra.Command{
	Use:   "create <key>",
	Short: "Create a new flag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body := map[string]interface{}{
			"key":           args[0],
			"type":          createType,
			"description":   createDescription,
			"enabled":       createEnabled,
			"default_value": json.RawMessage(createDefault),
		}

		data, status, err := doRequest("POST", "/api/v1/flags", body)
		if err != nil {
			return err
		}
		if status != 201 {
			return fmt.Errorf("server error (%d): %s", status, string(data))
		}
		fmt.Println(prettyJSON(data))
		return nil
	},
}

// --- flag enable ---

var flagEnableCmd = &cobra.Command{
	Use:   "enable <key>",
	Short: "Enable a flag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return setEnabled(args[0], true)
	},
}

// --- flag disable ---

var flagDisableCmd = &cobra.Command{
	Use:   "disable <key>",
	Short: "Disable a flag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return setEnabled(args[0], false)
	},
}

func setEnabled(key string, enabled bool) error {
	body := map[string]interface{}{"enabled": enabled}
	data, status, err := doRequest("PUT", "/api/v1/flags/"+key, body)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("server error (%d): %s", status, string(data))
	}

	state := "disabled"
	if enabled {
		state = "enabled"
	}
	fmt.Printf("Flag %q %s\n", key, state)
	return nil
}

// --- flag delete ---

var flagDeleteCmd = &cobra.Command{
	Use:   "delete <key>",
	Short: "Delete a flag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, status, err := doRequest("DELETE", "/api/v1/flags/"+args[0], nil)
		if err != nil {
			return err
		}
		if status != 204 {
			return fmt.Errorf("server error (%d): %s", status, string(data))
		}
		fmt.Printf("Flag %q deleted\n", args[0])
		return nil
	},
}

func init() {
	flagCreateCmd.Flags().StringVar(&createType, "type", "boolean", "Flag type (boolean, string, number, json)")
	flagCreateCmd.Flags().StringVar(&createDescription, "description", "", "Flag description")
	flagCreateCmd.Flags().BoolVar(&createEnabled, "enabled", false, "Enable the flag on creation")
	flagCreateCmd.Flags().StringVar(&createDefault, "default", "false", "Default value (JSON)")

	flagCmd.AddCommand(flagListCmd, flagGetCmd, flagCreateCmd, flagEnableCmd, flagDisableCmd, flagDeleteCmd)
	rootCmd.AddCommand(flagCmd)
}
