package cli

import (
	"os"

	"github.com/joho/godotenv"
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
	_ = godotenv.Load() // .env is optional, must run before os.Getenv

	defaultServer := "http://localhost:8080"
	if v := os.Getenv("FLAGGY_SERVER"); v != "" {
		defaultServer = v
	}

	defaultKey := os.Getenv("FLAGGY_MASTER_KEY")

	rootCmd.PersistentFlags().StringVar(&serverURL, "server", defaultServer, "Flaggy server URL")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", defaultKey, "API key or master key for authentication")
}

func Execute() error {
	return rootCmd.Execute()
}
