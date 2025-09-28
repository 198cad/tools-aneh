# Development Tools CLI

A comprehensive CLI tool for managing PostgreSQL, RabbitMQ, and MinIO services.

## Installation

```bash
go build -o tools.exe .
```

## Configuration

The tool supports configuration through environment variables. Create a `.env` file in the same directory as the executable:

```bash
cp .env.example .env
# Edit .env with your configuration
```

### Supported Environment Variables

#### PostgreSQL
- `PGHOST` - Database host (default: localhost)
- `PGPORT` - Database port (default: 5432)
- `PGUSER` - Database user (default: postgres)
- `PGPASSWORD` - Database password
- `PGDATABASE` - Default database (default: postgres)



#### RabbitMQ
- `RABBITMQ_HOST` - RabbitMQ host (default: localhost)
- `RABBITMQ_MANAGEMENT_PORT` - Management API port (default: 15672)
- `RABBITMQ_DEFAULT_USER` - Username (default: guest)
- `RABBITMQ_DEFAULT_PASS` - Password (default: guest)
- `RABBITMQ_DEFAULT_VHOST` - Virtual host (default: /)

#### MinIO
- `MINIO_ENDPOINT` - MinIO endpoint (default: localhost:9000)
- `MINIO_ACCESS_KEY` or `MINIO_ROOT_USER` - Access key
- `MINIO_SECRET_KEY` or `MINIO_ROOT_PASSWORD` - Secret key
- `MINIO_USE_SSL` - Use SSL connection (default: false)

## Usage

### Database Commands

```bash
# List all databases
./tools.exe db list

# Create a database
./tools.exe db create mydatabase

# Drop a database
./tools.exe db drop mydatabase

# Backup a database
./tools.exe db backup mydatabase backup.sql

# Restore a database
./tools.exe db restore mydatabase backup.sql


```

### RabbitMQ Commands

```bash
# List all queues
./tools.exe rabbit queues

# Create a queue
./tools.exe rabbit create-queue myqueue

# Delete a queue
./tools.exe rabbit delete-queue myqueue

# Purge messages from a queue
./tools.exe rabbit purge myqueue

# List exchanges
./tools.exe rabbit exchanges

# Create an exchange
./tools.exe rabbit create-exchange myexchange --type topic

# Publish a message
./tools.exe rabbit publish exchange routing-key "message content"

# Show statistics
./tools.exe rabbit stats
```

### MinIO Commands

```bash
# List all buckets
./tools.exe minio buckets

# Create a bucket
./tools.exe minio create-bucket mybucket

# Delete a bucket
./tools.exe minio delete-bucket mybucket
./tools.exe minio delete-bucket mybucket --force  # Delete with all contents

# List objects in a bucket
./tools.exe minio list mybucket
./tools.exe minio list mybucket --recursive

# Upload a file
./tools.exe minio upload mybucket localfile.txt
./tools.exe minio upload mybucket localfile.txt remote-name.txt

# Download a file
./tools.exe minio download mybucket remote-file.txt
./tools.exe minio download mybucket remote-file.txt local-name.txt

# Delete an object
./tools.exe minio delete mybucket object-name.txt

# Copy objects between buckets
./tools.exe minio copy source-bucket file.txt dest-bucket new-file.txt

# Get object/bucket information
./tools.exe minio stat mybucket
./tools.exe minio stat mybucket object.txt

# Mirror a local directory to a bucket
./tools.exe minio mirror ./local-dir mybucket
```

## Command-Line Flags

All commands support overriding environment variables with command-line flags:

```bash
# Database with custom connection
./tools.exe db list --host 192.168.1.10 --port 5433 --user admin --password secret

# RabbitMQ with custom connection
./tools.exe rabbit queues --host rabbitmq.example.com --port 15673

# MinIO with custom endpoint
./tools.exe minio buckets --endpoint s3.example.com:9000 --access-key mykey --secret-key mysecret
```

## Development

### Requirements
- Go 1.21 or later
- Access to PostgreSQL/MySQL (for database commands)
- Access to RabbitMQ Management API (for RabbitMQ commands)
- Access to MinIO server (for MinIO commands)

### Building from Source

```bash
git clone <repository>
cd tools
go mod download
go build -o tools.exe .
```

### Running with Air (Hot Reload)

```bash
air
```

## License

MIT