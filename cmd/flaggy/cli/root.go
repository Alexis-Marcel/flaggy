package cli

import (
	"github.com/spf13/cobra"
)

var (
	serverURL string
	apiKey    string
)

var rootCmd = &cobra.Command{
	Use:   "flaggy",
	Short: "Flaggy CLI — manage feature flags",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", "http://localhost:8080", "Flaggy server URL")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key or master key for authentication")
}

func Execute() error {
	return rootCmd.Execute()
}
