package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var segmentCmd = &cobra.Command{
	Use:   "segment",
	Short: "Manage segments",
}

// --- segment list ---

var segmentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all segments",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, status, err := doRequest("GET", "/api/v1/segments", nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("server error (%d): %s", status, string(data))
		}

		var segments []struct {
			Key         string `json:"key"`
			Description string `json:"description"`
		}
		if err := json.Unmarshal(data, &segments); err != nil {
			return fmt.Errorf("parse response: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "KEY\tDESCRIPTION")
		for _, s := range segments {
			fmt.Fprintf(w, "%s\t%s\n", s.Key, s.Description)
		}
		w.Flush()
		return nil
	},
}

// --- segment get ---

var segmentGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a segment by key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, status, err := doRequest("GET", "/api/v1/segments/"+args[0], nil)
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

// --- segment create ---

var (
	segCreateDescription string
	segCreateConditions  string
)

var segmentCreateCmd = &cobra.Command{
	Use:   "create <key>",
	Short: "Create a new segment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var conditions json.RawMessage
		if err := json.Unmarshal([]byte(segCreateConditions), &conditions); err != nil {
			return fmt.Errorf("invalid conditions JSON: %w", err)
		}

		body := map[string]interface{}{
			"key":         args[0],
			"description": segCreateDescription,
			"conditions":  conditions,
		}

		data, status, err := doRequest("POST", "/api/v1/segments", body)
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

// --- segment delete ---

var segmentDeleteCmd = &cobra.Command{
	Use:   "delete <key>",
	Short: "Delete a segment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, status, err := doRequest("DELETE", "/api/v1/segments/"+args[0], nil)
		if err != nil {
			return err
		}
		if status != 204 {
			return fmt.Errorf("server error (%d): %s", status, string(data))
		}
		fmt.Printf("Segment %q deleted\n", args[0])
		return nil
	},
}

func init() {
	segmentCreateCmd.Flags().StringVar(&segCreateDescription, "description", "", "Segment description")
	segmentCreateCmd.Flags().StringVar(&segCreateConditions, "conditions", "[]", "Conditions as JSON array")

	segmentCmd.AddCommand(segmentListCmd, segmentGetCmd, segmentCreateCmd, segmentDeleteCmd)
	rootCmd.AddCommand(segmentCmd)
}
