package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show configuration for services",
	Long:  `Display current configuration for PostgreSQL, MySQL, RabbitMQ, and MinIO`,
}

var configAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Show all configurations",
	Run: func(cmd *cobra.Command, args []string) {
		showAllConfig()
	},
}

var configPostgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Show PostgreSQL configuration",
	Run: func(cmd *cobra.Command, args []string) {
		showPostgresConfig()
	},
}

var configRabbitCmd = &cobra.Command{
	Use:   "rabbitmq",
	Short: "Show RabbitMQ configuration",
	Run: func(cmd *cobra.Command, args []string) {
		showRabbitConfig()
	},
}

var configMinioCmd = &cobra.Command{
	Use:   "minio",
	Short: "Show MinIO configuration",
	Run: func(cmd *cobra.Command, args []string) {
		showMinioConfig()
	},
}

var configEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Show environment variables",
	Run: func(cmd *cobra.Command, args []string) {
		showEnvironmentVariables()
	},
}

func init() {
	configCmd.AddCommand(configAllCmd)
	configCmd.AddCommand(configPostgresCmd)

	configCmd.AddCommand(configRabbitCmd)
	configCmd.AddCommand(configMinioCmd)
	configCmd.AddCommand(configEnvCmd)
}

func getConfigSource(key string, envVars ...string) string {
	// Check if value comes from environment variable
	for _, env := range envVars {
		if os.Getenv(env) != "" {
			return fmt.Sprintf("(from %s)", env)
		}
	}

	// Check if value comes from TOOLS_ prefixed env
	toolsEnv := "TOOLS_" + strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	if os.Getenv(toolsEnv) != "" {
		return fmt.Sprintf("(from %s)", toolsEnv)
	}

	return "(default)"
}

func maskPassword(password string) string {
	if password == "" {
		return "(not set)"
	}
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + strings.Repeat("*", len(password)-4) + password[len(password)-2:]
}

func showPostgresConfig() {
	color.Green("=== PostgreSQL Configuration ===")
	fmt.Println()

	host := viper.GetString("postgres.host")
	port := viper.GetString("postgres.port")
	user := viper.GetString("postgres.user")
	password := viper.GetString("postgres.password")
	database := viper.GetString("postgres.database")

	fmt.Printf("Host:     %s %s\n", host, getConfigSource("postgres.host", "PGHOST"))
	fmt.Printf("Port:     %s %s\n", port, getConfigSource("postgres.port", "PGPORT"))
	fmt.Printf("User:     %s %s\n", user, getConfigSource("postgres.user", "PGUSER"))
	fmt.Printf("Password: %s %s\n", maskPassword(password), getConfigSource("postgres.password", "PGPASSWORD"))
	fmt.Printf("Database: %s %s\n", database, getConfigSource("postgres.database", "PGDATABASE"))
	fmt.Println()

	// Show connection string
	color.Cyan("Connection String:")
	if password != "" {
		fmt.Printf("postgresql://%s:%s@%s:%s/%s\n", user, maskPassword(password), host, port, database)
	} else {
		fmt.Printf("postgresql://%s@%s:%s/%s\n", user, host, port, database)
	}

	// Test connection command
	color.Yellow("\nTest Connection:")
	fmt.Printf("psql -h %s -p %s -U %s -d %s -c \"SELECT version();\"\n", host, port, user, database)
}

func showRabbitConfig() {
	color.Green("=== RabbitMQ Configuration ===")
	fmt.Println()

	host := viper.GetString("rabbitmq.host")
	port := viper.GetString("rabbitmq.port")
	user := viper.GetString("rabbitmq.user")
	password := viper.GetString("rabbitmq.password")
	vhost := viper.GetString("rabbitmq.vhost")

	fmt.Printf("Host:            %s %s\n", host, getConfigSource("rabbitmq.host", "RABBITMQ_HOST"))
	fmt.Printf("Management Port: %s %s\n", port, getConfigSource("rabbitmq.port", "RABBITMQ_MANAGEMENT_PORT"))
	fmt.Printf("User:            %s %s\n", user, getConfigSource("rabbitmq.user", "RABBITMQ_DEFAULT_USER"))
	fmt.Printf("Password:        %s %s\n", maskPassword(password), getConfigSource("rabbitmq.password", "RABBITMQ_DEFAULT_PASS"))
	fmt.Printf("Virtual Host:    %s %s\n", vhost, getConfigSource("rabbitmq.vhost", "RABBITMQ_DEFAULT_VHOST"))
	fmt.Println()

	// Show Management URL
	color.Cyan("Management URL:")
	fmt.Printf("http://%s:%s\n", host, port)
	fmt.Printf("Login: %s / %s\n", user, maskPassword(password))
	fmt.Println()

	// Show AMQP Connection String
	color.Cyan("AMQP Connection String:")
	encodedVhost := vhost
	if vhost == "/" {
		encodedVhost = "%2F"
	}
	fmt.Printf("amqp://%s:%s@%s:5672/%s\n", user, maskPassword(password), host, encodedVhost)

	// Test connection
	color.Yellow("\nTest Connection:")
	fmt.Printf("curl -u %s:%s http://%s:%s/api/overview\n", user, maskPassword(password), host, port)
}

func showMinioConfig() {
	color.Green("=== MinIO Configuration ===")
	fmt.Println()

	endpoint := viper.GetString("minio.endpoint")
	accessKey := viper.GetString("minio.access_key")
	secretKey := viper.GetString("minio.secret_key")
	useSSL := viper.GetBool("minio.use_ssl")

	fmt.Printf("Endpoint:    %s %s\n", endpoint, getConfigSource("minio.endpoint", "MINIO_ENDPOINT"))
	fmt.Printf("Access Key:  %s %s\n", accessKey, getConfigSource("minio.access_key", "MINIO_ACCESS_KEY", "MINIO_ROOT_USER"))
	fmt.Printf("Secret Key:  %s %s\n", maskPassword(secretKey), getConfigSource("minio.secret_key", "MINIO_SECRET_KEY", "MINIO_ROOT_PASSWORD"))
	fmt.Printf("Use SSL:     %v %s\n", useSSL, getConfigSource("minio.use_ssl", "MINIO_USE_SSL"))
	fmt.Println()

	// Show Console URL
	protocol := "http"
	if useSSL {
		protocol = "https"
	}

	// Extract host from endpoint (remove port if exists)
	host := endpoint
	if idx := strings.LastIndex(endpoint, ":"); idx != -1 {
		host = endpoint[:idx]
	}

	color.Cyan("Console URL:")
	fmt.Printf("%s://%s:9001\n", protocol, host)
	fmt.Printf("Login: %s / %s\n", accessKey, maskPassword(secretKey))
	fmt.Println()

	// Show S3 Endpoint
	color.Cyan("S3 API Endpoint:")
	fmt.Printf("%s://%s\n", protocol, endpoint)

	// Test connection with mc (MinIO Client)
	color.Yellow("\nTest Connection (using mc):")
	fmt.Printf("mc alias set myminio %s://%s %s %s\n", protocol, endpoint, accessKey, maskPassword(secretKey))
	fmt.Printf("mc ls myminio\n")
}

func showAllConfig() {
	showPostgresConfig()
	fmt.Println()
	showRabbitConfig()
	fmt.Println()
	showMinioConfig()
}

func showEnvironmentVariables() {
	color.Green("=== Environment Variables ===")
	fmt.Println()

	envVars := []struct {
		category string
		vars     []string
	}{
		{
			"PostgreSQL",
			[]string{"PGHOST", "PGPORT", "PGUSER", "PGPASSWORD", "PGDATABASE",
				"TOOLS_POSTGRES_HOST", "TOOLS_POSTGRES_PORT", "TOOLS_POSTGRES_USER",
				"TOOLS_POSTGRES_PASSWORD", "TOOLS_POSTGRES_DATABASE"},
		},

		{
			"RabbitMQ",
			[]string{"RABBITMQ_HOST", "RABBITMQ_MANAGEMENT_PORT", "RABBITMQ_DEFAULT_USER",
				"RABBITMQ_DEFAULT_PASS", "RABBITMQ_DEFAULT_VHOST",
				"TOOLS_RABBITMQ_HOST", "TOOLS_RABBITMQ_PORT", "TOOLS_RABBITMQ_USER",
				"TOOLS_RABBITMQ_PASSWORD", "TOOLS_RABBITMQ_VHOST"},
		},
		{
			"MinIO",
			[]string{"MINIO_ENDPOINT", "MINIO_ACCESS_KEY", "MINIO_SECRET_KEY", "MINIO_USE_SSL",
				"MINIO_ROOT_USER", "MINIO_ROOT_PASSWORD",
				"TOOLS_MINIO_ENDPOINT", "TOOLS_MINIO_ACCESS_KEY", "TOOLS_MINIO_SECRET_KEY",
				"TOOLS_MINIO_USE_SSL"},
		},
	}

	for _, group := range envVars {
		color.Yellow("%s:\n", group.category)
		hasValue := false
		for _, env := range group.vars {
			value := os.Getenv(env)
			if value != "" {
				// Mask passwords
				if strings.Contains(strings.ToLower(env), "password") ||
					strings.Contains(strings.ToLower(env), "secret") ||
					strings.Contains(strings.ToLower(env), "pass") {
					value = maskPassword(value)
				}
				fmt.Printf("  %s = %s\n", env, value)
				hasValue = true
			}
		}
		if !hasValue {
			fmt.Printf("  (none set)\n")
		}
		fmt.Println()
	}

	// Show .env file status
	color.Yellow(".env File:")
	if _, err := os.Stat(".env"); err == nil {
		fmt.Println("  Found in current directory")
	} else {
		fmt.Println("  Not found in current directory")
	}

	if execPath, err := os.Executable(); err == nil {
		envPath := fmt.Sprintf("%s\\.env", strings.TrimSuffix(execPath, "\\tools.exe"))
		if _, err := os.Stat(envPath); err == nil {
			fmt.Printf("  Found at: %s\n", envPath)
		}
	}
}

func GetConfigCommand() *cobra.Command {
	return configCmd
}
