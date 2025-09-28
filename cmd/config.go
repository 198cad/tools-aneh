package cmd

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func InitConfig() {
	// Load .env file if it exists
	if execPath, err := os.Executable(); err == nil {
		envPath := filepath.Join(filepath.Dir(execPath), ".env")
		godotenv.Load(envPath)
	}
	// Also try loading from current directory
	godotenv.Load(".env")

	// PostgreSQL environment variables
	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", "5432")
	viper.SetDefault("postgres.user", "postgres")
	viper.SetDefault("postgres.password", "")
	viper.SetDefault("postgres.database", "postgres")

	// RabbitMQ environment variables
	viper.SetDefault("rabbitmq.host", "localhost")
	viper.SetDefault("rabbitmq.port", "15672")
	viper.SetDefault("rabbitmq.user", "guest")
	viper.SetDefault("rabbitmq.password", "guest")
	viper.SetDefault("rabbitmq.vhost", "/")

	// MinIO environment variables
	viper.SetDefault("minio.endpoint", "localhost:9000")
	viper.SetDefault("minio.access_key", "minioadmin")
	viper.SetDefault("minio.secret_key", "minioadmin")
	viper.SetDefault("minio.use_ssl", false)

	// Bind environment variables
	viper.SetEnvPrefix("TOOLS")
	viper.AutomaticEnv()

	// Check for specific environment variables
	// PostgreSQL
	if pgHost := os.Getenv("PGHOST"); pgHost != "" {
		viper.Set("postgres.host", pgHost)
	}
	if pgPort := os.Getenv("PGPORT"); pgPort != "" {
		viper.Set("postgres.port", pgPort)
	}
	if pgUser := os.Getenv("PGUSER"); pgUser != "" {
		viper.Set("postgres.user", pgUser)
	}
	if pgPass := os.Getenv("PGPASSWORD"); pgPass != "" {
		viper.Set("postgres.password", pgPass)
	}
	if pgDB := os.Getenv("PGDATABASE"); pgDB != "" {
		viper.Set("postgres.database", pgDB)
	}

	// RabbitMQ
	if rmqHost := os.Getenv("RABBITMQ_HOST"); rmqHost != "" {
		viper.Set("rabbitmq.host", rmqHost)
	}
	if rmqPort := os.Getenv("RABBITMQ_MANAGEMENT_PORT"); rmqPort != "" {
		viper.Set("rabbitmq.port", rmqPort)
	}
	if rmqUser := os.Getenv("RABBITMQ_DEFAULT_USER"); rmqUser != "" {
		viper.Set("rabbitmq.user", rmqUser)
	}
	if rmqPass := os.Getenv("RABBITMQ_DEFAULT_PASS"); rmqPass != "" {
		viper.Set("rabbitmq.password", rmqPass)
	}
	if rmqVHost := os.Getenv("RABBITMQ_DEFAULT_VHOST"); rmqVHost != "" {
		viper.Set("rabbitmq.vhost", rmqVHost)
	}

	// MinIO
	if minioEndpoint := os.Getenv("MINIO_ENDPOINT"); minioEndpoint != "" {
		viper.Set("minio.endpoint", minioEndpoint)
	}
	if minioAccess := os.Getenv("MINIO_ACCESS_KEY"); minioAccess != "" {
		viper.Set("minio.access_key", minioAccess)
	}
	if minioSecret := os.Getenv("MINIO_SECRET_KEY"); minioSecret != "" {
		viper.Set("minio.secret_key", minioSecret)
	}
	if minioSSL := os.Getenv("MINIO_USE_SSL"); minioSSL == "true" {
		viper.Set("minio.use_ssl", true)
	}

	// Alternative MinIO env vars
	if minioAccess := os.Getenv("MINIO_ROOT_USER"); minioAccess != "" {
		viper.Set("minio.access_key", minioAccess)
	}
	if minioSecret := os.Getenv("MINIO_ROOT_PASSWORD"); minioSecret != "" {
		viper.Set("minio.secret_key", minioSecret)
	}
}
