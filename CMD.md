# Development Tools - Complete Command Reference

## Overview
A comprehensive CLI tool for managing PostgreSQL, RabbitMQ, and MinIO services.

**Version**: 1.0.1  
**Usage**: `tools [command] [flags]`

---

## Global Commands

### Help & Version
```bash
tools --help              # Show help for tools
tools --version          # Show version information
tools [command] --help   # Show help for specific command
```

### Configuration Management
```bash
tools config all         # Show all service configurations
tools config postgres    # Show PostgreSQL configuration
tools config rabbitmq    # Show RabbitMQ configuration  
tools config minio       # Show MinIO configuration
tools config env         # Show all environment variables
```

### Update & Maintenance
```bash
tools update --check     # Check for available updates
tools update             # Install latest updates
tools update --force     # Force update even if up to date
```

---

## PostgreSQL Database Commands (`db`)

### Database Operations
```bash
# List all databases
tools db list

# Create a new database
tools db create [database-name]

# Drop a database (with confirmation)
tools db drop [database-name]

# Backup database to file
tools db backup [database-name] [output-file.sql]

# Restore database from backup
tools db restore [database-name] [backup-file.sql]

# Execute SQL query
tools db exec [database-name] "[SQL query]"

# Show database size
tools db size                    # Show all database sizes
tools db size [database-name]    # Show specific database size
```

### Table Operations
```bash
# List all tables in database
tools db tables [database-name]

# Show table details and recent records
tools db table [database-name] [table-name]
tools db table [database-name] [table-name] --limit 10  # Custom record limit
```

### Connection Flags
```bash
--host, -H      # Database host (env: PGHOST)
--port, -p      # Database port (env: PGPORT)  
--user, -u      # Database user (env: PGUSER)
--password, -P  # Database password (env: PGPASSWORD)
--database, -d  # Default database (env: PGDATABASE)
```

### Examples
```bash
# Create and setup a new database
tools db create myapp_db
tools db exec myapp_db "CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(100))"
tools db tables myapp_db
tools db table myapp_db users

# Backup and restore
tools db backup production_db prod_backup_2024.sql
tools db restore staging_db prod_backup_2024.sql

# Remote database connection
tools db list --host 192.168.1.100 --port 5433 --user admin
```

---

## RabbitMQ Management Commands (`rabbit`)

### Queue Operations
```bash
# List all queues
tools rabbit queues

# Create a new queue
tools rabbit create-queue [queue-name]
tools rabbit create-queue [queue-name] --durable --auto-delete

# Delete a queue
tools rabbit delete-queue [queue-name]

# Purge messages from queue
tools rabbit purge [queue-name]
```

### Exchange Operations
```bash
# List all exchanges
tools rabbit exchanges

# Create an exchange
tools rabbit create-exchange [exchange-name]
tools rabbit create-exchange [exchange-name] --type topic --durable

# Publish a message
tools rabbit publish [exchange] [routing-key] "[message]"
```

### Monitoring
```bash
# Show RabbitMQ statistics
tools rabbit stats
```

### Connection Flags
```bash
--host, -H      # RabbitMQ host (env: RABBITMQ_HOST)
--port, -p      # Management port (env: RABBITMQ_MANAGEMENT_PORT)
--user, -u      # Username (env: RABBITMQ_DEFAULT_USER)
--password, -P  # Password (env: RABBITMQ_DEFAULT_PASS)
--vhost, -v     # Virtual host (env: RABBITMQ_DEFAULT_VHOST)
```

### Examples
```bash
# Create queue and publish message
tools rabbit create-queue notifications --durable
tools rabbit publish "" notifications "Hello World"

# Monitor queues
tools rabbit queues
tools rabbit stats

# Work with different vhost
tools rabbit queues --vhost /production
```

---

## MinIO Object Storage Commands (`minio`)

### Bucket Operations
```bash
# List all buckets
tools minio buckets

# Create a bucket
tools minio create-bucket [bucket-name]
tools minio create-bucket [bucket-name] --region us-west-1

# Delete a bucket
tools minio delete-bucket [bucket-name]
tools minio delete-bucket [bucket-name] --force  # Delete with contents
```

### Object Operations
```bash
# List objects in bucket
tools minio list [bucket-name]
tools minio list [bucket-name] --prefix uploads/
tools minio list [bucket-name] --recursive

# Upload file to bucket
tools minio upload [bucket] [local-file]
tools minio upload [bucket] [local-file] [remote-name]

# Download object from bucket
tools minio download [bucket] [object-name]
tools minio download [bucket] [object-name] [local-file]

# Delete object
tools minio delete [bucket] [object-name]

# Copy object between buckets
tools minio copy [source-bucket] [object] [dest-bucket] [new-name]
```

### Information & Sync
```bash
# Get bucket/object information
tools minio stat [bucket-name]
tools minio stat [bucket-name] [object-name]

# Mirror local directory to bucket
tools minio mirror [local-dir] [bucket-name]
tools minio mirror [local-dir] [bucket-name] --prefix backup/
```

### Connection Flags
```bash
--endpoint, -e     # MinIO endpoint (env: MINIO_ENDPOINT)
--access-key, -a   # Access key (env: MINIO_ACCESS_KEY)
--secret-key, -s   # Secret key (env: MINIO_SECRET_KEY)
--use-ssl, -S      # Use SSL connection (env: MINIO_USE_SSL)
```

### Examples
```bash
# Create bucket and upload files
tools minio create-bucket my-backups
tools minio upload my-backups backup.tar.gz
tools minio list my-backups

# Mirror entire directory
tools minio mirror ./website static-website

# Work with remote MinIO
tools minio buckets --endpoint s3.example.com:9000 --use-ssl
```

---

## Environment Variables

### PostgreSQL
```bash
PGHOST=localhost
PGPORT=5432
PGUSER=postgres
PGPASSWORD=your_password
PGDATABASE=postgres
```

### RabbitMQ
```bash
RABBITMQ_HOST=localhost
RABBITMQ_MANAGEMENT_PORT=15672
RABBITMQ_DEFAULT_USER=guest
RABBITMQ_DEFAULT_PASS=guest
RABBITMQ_DEFAULT_VHOST=/
```

### MinIO
```bash
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin       # or MINIO_ROOT_USER
MINIO_SECRET_KEY=minioadmin       # or MINIO_ROOT_PASSWORD
MINIO_USE_SSL=false
```

### Alternative Prefix
All environment variables can also use `TOOLS_` prefix:
```bash
TOOLS_POSTGRES_HOST=localhost
TOOLS_RABBITMQ_HOST=localhost
TOOLS_MINIO_ENDPOINT=localhost:9000
```

---

## Configuration File (.env)

Create a `.env` file in the tools directory or current directory:

```env
# PostgreSQL
PGHOST=10.8.0.1
PGPORT=5432
PGUSER=postgres
PGPASSWORD=password
PGDATABASE=postgres

# RabbitMQ
RABBITMQ_HOST=10.8.0.1
RABBITMQ_MANAGEMENT_PORT=15672
RABBITMQ_DEFAULT_USER=admin
RABBITMQ_DEFAULT_PASS=secret
RABBITMQ_DEFAULT_VHOST=/

# MinIO
MINIO_ENDPOINT=10.8.0.1:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
```

---

## Quick Start Examples

### Database Development Workflow
```bash
# Setup development database
tools db create dev_db
tools db exec dev_db "$(cat schema.sql)"
tools db tables dev_db

# Backup before deployment
tools db backup dev_db dev_backup.sql

# Check production
tools db list --host prod.example.com
tools db size production_db
```

### Message Queue Setup
```bash
# Create queues for microservices
tools rabbit create-queue orders --durable
tools rabbit create-queue notifications --durable
tools rabbit create-exchange events --type topic

# Monitor
tools rabbit stats
tools rabbit queues
```

### Object Storage Management
```bash
# Setup buckets
tools minio create-bucket uploads
tools minio create-bucket backups
tools minio create-bucket static-assets

# Upload assets
tools minio upload static-assets logo.png
tools minio mirror ./public static-assets

# List and verify
tools minio list static-assets --recursive
```

---

## Installation & Updates

### Install
```bash
# Using PowerShell (Admin)
.\install.ps1

# Using Command Prompt (Admin)
install.bat
```

### Update
```bash
# Check and install updates
tools update --check
tools update

# Using update scripts
.\update.ps1
update.bat
```

### Uninstall
```bash
uninstall.bat
```

---

## Tips & Best Practices

1. **Use .env file** for permanent configuration instead of flags
2. **Backup databases** regularly using `tools db backup`
3. **Monitor RabbitMQ** with `tools rabbit stats` for performance
4. **Use --force carefully** when deleting buckets with content
5. **Check updates** regularly with `tools update --check`
6. **Use tab completion** (if available) for faster command entry
7. **Combine with scripts** for automation:
   ```bash
   # Backup all databases
   for db in $(tools db list | grep -o '^\w\+'); do
     tools db backup $db backups/$db_$(date +%Y%m%d).sql
   done
   ```

---

## Troubleshooting

### Connection Issues
```bash
# Test configuration
tools config all

# Verify credentials
tools config postgres
tools config rabbitmq
tools config minio

# Check environment variables
tools config env
```

### PostgreSQL Issues
- Ensure PostgreSQL client tools are installed
- Check firewall settings for port 5432
- Verify user permissions

### RabbitMQ Issues
- Management plugin must be enabled (port 15672)
- Check vhost permissions
- Verify management API is accessible

### MinIO Issues
- Check both API (9000) and Console (9001) ports
- Verify access/secret keys
- Ensure bucket names are DNS-compatible

---

## Support

For issues or feature requests:
- Check help: `tools [command] --help`
- View configuration: `tools config all`
- Update to latest: `tools update`
- Report issues: https://github.com/yourusername/tools/issues