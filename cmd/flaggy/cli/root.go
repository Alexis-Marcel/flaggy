package cli

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	serverURL string
	apiKey    string
	Version   = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "flaggy",
	Short: "Flaggy CLI — manage feature flags",
}

func init() {
	_ = godotenv.Load() // .env is optional, must run before os.Getenv

	port := os.Getenv("FLAGGY_PORT")
	if port == "" {
		port = ":8080"
	}
	defaultServer := "http://localhost" + port

	defaultKey := os.Getenv("FLAGGY_MASTER_KEY")

	rootCmd.PersistentFlags().StringVar(&serverURL, "server", defaultServer, "Flaggy server URL")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", defaultKey, "API key or master key for authentication")
}

func Execute() error {
	return rootCmd.Execute()
}
