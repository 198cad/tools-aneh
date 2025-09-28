package main

import (
	"fmt"
	"os"
	"tools/cmd"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tools",
	Short: "Development tools for PostgreSQL, RabbitMQ, and MinIO management",
	Long: `A comprehensive CLI tool for managing various development services:
- PostgreSQL database management
- RabbitMQ queue management
- MinIO object storage management

Environment Variables:
  PostgreSQL:
    PGHOST, PGPORT, PGUSER, PGPASSWORD, PGDATABASE
    or TOOLS_POSTGRES_HOST, TOOLS_POSTGRES_PORT, etc.

  RabbitMQ:
    RABBITMQ_HOST, RABBITMQ_MANAGEMENT_PORT, RABBITMQ_DEFAULT_USER,
    RABBITMQ_DEFAULT_PASS, RABBITMQ_DEFAULT_VHOST
    or TOOLS_RABBITMQ_HOST, TOOLS_RABBITMQ_PORT, etc.

  MinIO:
    MINIO_ENDPOINT, MINIO_ACCESS_KEY (or MINIO_ROOT_USER),
    MINIO_SECRET_KEY (or MINIO_ROOT_PASSWORD), MINIO_USE_SSL
    or TOOLS_MINIO_ENDPOINT, TOOLS_MINIO_ACCESS_KEY, etc.`,
	Version: Version,
}

func main() {
	// Initialize configuration
	cmd.InitConfig()

	// Register commands
	rootCmd.AddCommand(cmd.GetDBCommand())
	rootCmd.AddCommand(cmd.GetRabbitMQCommand())
	rootCmd.AddCommand(cmd.GetMinIOCommand())
	rootCmd.AddCommand(cmd.GetConfigCommand())
	rootCmd.AddCommand(cmd.GetUpdateCommand())

	// Disable default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
