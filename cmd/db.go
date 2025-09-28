package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "PostgreSQL database management commands",
	Long:  `Manage PostgreSQL databases with various operations including list, create, drop, backup, and restore`,
}

var dbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all PostgreSQL databases",
	Run: func(cmd *cobra.Command, args []string) {
		listDatabases()
	},
}

var dbCreateCmd = &cobra.Command{
	Use:   "create [database-name]",
	Short: "Create a new PostgreSQL database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createDatabase(args[0])
	},
}

var dbDropCmd = &cobra.Command{
	Use:   "drop [database-name]",
	Short: "Drop a PostgreSQL database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dropDatabase(args[0])
	},
}

var dbBackupCmd = &cobra.Command{
	Use:   "backup [database-name] [output-file]",
	Short: "Backup a PostgreSQL database",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		backupDatabase(args[0], args[1])
	},
}

var dbRestoreCmd = &cobra.Command{
	Use:   "restore [database-name] [backup-file]",
	Short: "Restore a PostgreSQL database from backup",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		restoreDatabase(args[0], args[1])
	},
}

var dbExecCmd = &cobra.Command{
	Use:   "exec [database-name] [sql-query]",
	Short: "Execute SQL query on a PostgreSQL database",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		executeQuery(args[0], args[1])
	},
}

var dbSizeCmd = &cobra.Command{
	Use:   "size [database-name]",
	Short: "Show PostgreSQL database size",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showAllDatabaseSizes()
		} else {
			showDatabaseSize(args[0])
		}
	},
}

var dbTablesCmd = &cobra.Command{
	Use:   "tables [database-name]",
	Short: "List all tables in a PostgreSQL database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		listTables(args[0])
	},
}

var dbTableCmd = &cobra.Command{
	Use:   "table [database-name] [table-name]",
	Short: "Show table structure and recent records",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		showTableDetails(args[0], args[1], limit)
	},
}

func init() {
	dbCmd.PersistentFlags().StringP("host", "H", "", "Database host (env: PGHOST)")
	dbCmd.PersistentFlags().StringP("port", "p", "", "Database port (env: PGPORT)")
	dbCmd.PersistentFlags().StringP("user", "u", "", "Database user (env: PGUSER)")
	dbCmd.PersistentFlags().StringP("password", "P", "", "Database password (env: PGPASSWORD)")
	dbCmd.PersistentFlags().StringP("database", "d", "", "Default database (env: PGDATABASE)")

	// Add limit flag for table command
	dbTableCmd.Flags().IntP("limit", "l", 5, "Number of records to show")

	dbCmd.AddCommand(dbListCmd)
	dbCmd.AddCommand(dbCreateCmd)
	dbCmd.AddCommand(dbDropCmd)
	dbCmd.AddCommand(dbBackupCmd)
	dbCmd.AddCommand(dbRestoreCmd)
	dbCmd.AddCommand(dbExecCmd)
	dbCmd.AddCommand(dbSizeCmd)
	dbCmd.AddCommand(dbTablesCmd)
	dbCmd.AddCommand(dbTableCmd)
}

func getPostgresConfig() (host, port, user, password, database string) {
	host, _ = dbCmd.Flags().GetString("host")
	port, _ = dbCmd.Flags().GetString("port")
	user, _ = dbCmd.Flags().GetString("user")
	password, _ = dbCmd.Flags().GetString("password")
	database, _ = dbCmd.Flags().GetString("database")

	if host == "" {
		host = viper.GetString("postgres.host")
	}
	if port == "" {
		port = viper.GetString("postgres.port")
	}
	if user == "" {
		user = viper.GetString("postgres.user")
	}
	if password == "" {
		password = viper.GetString("postgres.password")
	}
	if database == "" {
		database = viper.GetString("postgres.database")
	}

	return
}

func listDatabases() {
	host, port, user, password, _ := getPostgresConfig()

	cmd := exec.Command("psql", "-h", host, "-p", port, "-U", user, "-l")
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red("Error: %v", err)
		color.Yellow("Make sure PostgreSQL client tools are installed and accessible.")
		color.Yellow("Install with: apt-get install postgresql-client (Linux) or download from postgresql.org (Windows)")
		return
	}
	fmt.Println(string(output))
}

func createDatabase(dbName string) {
	host, port, user, password, _ := getPostgresConfig()

	cmd := exec.Command("createdb", "-h", host, "-p", port, "-U", user, dbName)
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red("Error creating database: %v\n%s", err, output)
		return
	}
	color.Green("✓ Database '%s' created successfully", dbName)
}

func dropDatabase(dbName string) {
	host, port, user, password, _ := getPostgresConfig()

	// Safety check
	if dbName == "postgres" || dbName == "template0" || dbName == "template1" {
		color.Red("Cannot drop system database '%s'", dbName)
		return
	}

	// Confirm deletion
	color.Yellow("WARNING: This will permanently delete database '%s'", dbName)
	fmt.Print("Type the database name to confirm deletion: ")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != dbName {
		color.Yellow("Deletion cancelled")
		return
	}

	cmd := exec.Command("dropdb", "-h", host, "-p", port, "-U", user, dbName)
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red("Error dropping database: %v\n%s", err, output)
		return
	}
	color.Green("✓ Database '%s' dropped successfully", dbName)
}

func backupDatabase(dbName, outputFile string) {
	host, port, user, password, _ := getPostgresConfig()

	// Add .sql extension if not present
	if len(outputFile) < 4 || outputFile[len(outputFile)-4:] != ".sql" {
		outputFile += ".sql"
	}

	color.Yellow("Backing up database '%s' to '%s'...", dbName, outputFile)

	cmd := exec.Command("pg_dump",
		"-h", host,
		"-p", port,
		"-U", user,
		"-d", dbName,
		"-f", outputFile,
		"--verbose",
		"--no-owner",
		"--no-acl")
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red("Error backing up database: %v\n%s", err, output)
		return
	}

	// Get file size
	if fileInfo, err := os.Stat(outputFile); err == nil {
		size := fileInfo.Size()
		sizeStr := formatBytes(size)
		color.Green("✓ Database '%s' backed up to '%s' (%s)", dbName, outputFile, sizeStr)
	} else {
		color.Green("✓ Database '%s' backed up to '%s'", dbName, outputFile)
	}
}

func restoreDatabase(dbName, backupFile string) {
	host, port, user, password, _ := getPostgresConfig()

	// Check if backup file exists
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		color.Red("Backup file '%s' not found", backupFile)
		return
	}

	color.Yellow("Restoring database '%s' from '%s'...", dbName, backupFile)

	// First, create the database if it doesn't exist
	createCmd := exec.Command("createdb", "-h", host, "-p", port, "-U", user, dbName)
	createCmd.Env = os.Environ()
	if password != "" {
		createCmd.Env = append(createCmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}
	createCmd.CombinedOutput() // Ignore error if database already exists

	// Restore the database
	cmd := exec.Command("psql", "-h", host, "-p", port, "-U", user, "-d", dbName, "-f", backupFile)
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red("Error restoring database: %v\n%s", err, output)
		return
	}
	color.Green("✓ Database '%s' restored from '%s'", dbName, backupFile)
}

func executeQuery(dbName, query string) {
	host, port, user, password, _ := getPostgresConfig()

	cmd := exec.Command("psql", "-h", host, "-p", port, "-U", user, "-d", dbName, "-c", query)
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red("Error executing query: %v", err)
		return
	}
	fmt.Println(string(output))
}

func showDatabaseSize(dbName string) {
	query := fmt.Sprintf("SELECT pg_database_size('%s'), pg_size_pretty(pg_database_size('%s'));", dbName, dbName)

	host, port, user, password, _ := getPostgresConfig()

	cmd := exec.Command("psql", "-h", host, "-p", port, "-U", user, "-d", "postgres", "-t", "-c", query)
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red("Error getting database size: %v", err)
		return
	}

	color.Green("Database '%s' size: %s", dbName, string(output))
}

func showAllDatabaseSizes() {
	query := `SELECT
		datname AS database_name,
		pg_size_pretty(pg_database_size(datname)) AS size
	FROM pg_database
	WHERE datistemplate = false
	ORDER BY pg_database_size(datname) DESC;`

	host, port, user, password, _ := getPostgresConfig()

	cmd := exec.Command("psql", "-h", host, "-p", port, "-U", user, "-d", "postgres", "-c", query)
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red("Error getting database sizes: %v", err)
		return
	}

	fmt.Println(string(output))
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func listTables(dbName string) {
	host, port, user, password, _ := getPostgresConfig()

	// Query to get all tables with additional information
	query := `
		SELECT
			schemaname AS schema,
			tablename AS table_name,
			pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
			obj_description((schemaname||'.'||tablename)::regclass, 'pg_class') AS comment
		FROM pg_tables
		WHERE schemaname NOT IN ('pg_catalog', 'information_schema')
		ORDER BY schemaname, tablename;`

	cmd := exec.Command("psql", "-h", host, "-p", port, "-U", user, "-d", dbName, "-c", query)
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		// If the first query fails, try a simpler one
		simpleQuery := `SELECT schemaname, tablename FROM pg_tables WHERE schemaname NOT IN ('pg_catalog', 'information_schema') ORDER BY schemaname, tablename;`
		cmd = exec.Command("psql", "-h", host, "-p", port, "-U", user, "-d", dbName, "-c", simpleQuery)
		cmd.Env = os.Environ()
		if password != "" {
			cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
		}

		output, err = cmd.CombinedOutput()
		if err != nil {
			color.Red("Error listing tables: %v", err)
			return
		}
	}

	color.Green("Tables in database '%s':", dbName)
	fmt.Println(string(output))

	// Also show count
	countQuery := `SELECT COUNT(*) as table_count FROM pg_tables WHERE schemaname NOT IN ('pg_catalog', 'information_schema');`
	cmd = exec.Command("psql", "-h", host, "-p", port, "-U", user, "-d", dbName, "-t", "-c", countQuery)
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	if countOutput, err := cmd.CombinedOutput(); err == nil {
		count := strings.TrimSpace(string(countOutput))
		color.Yellow("Total tables: %s", count)
	}
}

func showTableDetails(dbName, tableName string, limit int) {
	host, port, user, password, _ := getPostgresConfig()

	// Parse schema and table name
	schema := "public"
	table := tableName
	if strings.Contains(tableName, ".") {
		parts := strings.Split(tableName, ".")
		schema = parts[0]
		table = parts[1]
	}

	color.Green("=== Table: %s.%s in database: %s ===", schema, table, dbName)
	fmt.Println()

	// 1. Show table structure
	color.Cyan("Table Structure:")
	structureQuery := fmt.Sprintf(`
		SELECT
			column_name AS "Column",
			data_type AS "Type",
			character_maximum_length AS "Max Length",
			is_nullable AS "Nullable",
			column_default AS "Default"
		FROM information_schema.columns
		WHERE table_schema = '%s' AND table_name = '%s'
		ORDER BY ordinal_position;`, schema, table)

	cmd := exec.Command("psql", "-h", host, "-p", port, "-U", user, "-d", dbName, "-c", structureQuery)
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red("Error getting table structure: %v", err)
		return
	}
	fmt.Println(string(output))

	// 2. Show indexes
	color.Cyan("Indexes:")
	indexQuery := fmt.Sprintf(`
		SELECT
			indexname AS "Index Name",
			indexdef AS "Definition"
		FROM pg_indexes
		WHERE schemaname = '%s' AND tablename = '%s';`, schema, table)

	cmd = exec.Command("psql", "-h", host, "-p", port, "-U", user, "-d", dbName, "-c", indexQuery)
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	if output, err := cmd.CombinedOutput(); err == nil {
		fmt.Println(string(output))
	}

	// 3. Show foreign keys
	color.Cyan("Foreign Keys:")
	fkQuery := fmt.Sprintf(`
		SELECT
			tc.constraint_name AS "Constraint",
			kcu.column_name AS "Column",
			ccu.table_name AS "References Table",
			ccu.column_name AS "References Column"
		FROM information_schema.table_constraints AS tc
		JOIN information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
		JOIN information_schema.constraint_column_usage AS ccu
			ON ccu.constraint_name = tc.constraint_name
			AND ccu.table_schema = tc.table_schema
		WHERE tc.constraint_type = 'FOREIGN KEY'
			AND tc.table_schema = '%s'
			AND tc.table_name = '%s';`, schema, table)

	cmd = exec.Command("psql", "-h", host, "-p", port, "-U", user, "-d", dbName, "-c", fkQuery)
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	if output, err := cmd.CombinedOutput(); err == nil {
		outputStr := string(output)
		if strings.Contains(outputStr, "(0 rows)") {
			fmt.Println(" No foreign keys")
		} else {
			fmt.Println(outputStr)
		}
	}

	// 4. Show row count
	color.Cyan("Statistics:")
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s.%s;`, schema, table)
	cmd = exec.Command("psql", "-h", host, "-p", port, "-U", user, "-d", dbName, "-t", "-c", countQuery)
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	if output, err := cmd.CombinedOutput(); err == nil {
		count := strings.TrimSpace(string(output))
		fmt.Printf(" Total rows: %s\n", count)
	}

	// Show table size
	sizeQuery := fmt.Sprintf(`SELECT pg_size_pretty(pg_total_relation_size('%s.%s'));`, schema, table)
	cmd = exec.Command("psql", "-h", host, "-p", port, "-U", user, "-d", dbName, "-t", "-c", sizeQuery)
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	if output, err := cmd.CombinedOutput(); err == nil {
		size := strings.TrimSpace(string(output))
		fmt.Printf(" Table size: %s\n\n", size)
	}

	// 5. Show recent records
	color.Cyan("Recent %d Records:", limit)

	// Try to find a column to order by (prefer id, created_at, or primary key)
	orderColumn := ""
	orderQuery := fmt.Sprintf(`
		SELECT column_name
		FROM information_schema.columns
		WHERE table_schema = '%s' AND table_name = '%s'
		AND column_name IN ('id', 'created_at', 'updated_at', 'created', 'timestamp')
		ORDER BY
			CASE column_name
				WHEN 'id' THEN 1
				WHEN 'created_at' THEN 2
				WHEN 'created' THEN 3
				WHEN 'updated_at' THEN 4
				WHEN 'timestamp' THEN 5
			END
		LIMIT 1;`, schema, table)

	cmd = exec.Command("psql", "-h", host, "-p", port, "-U", user, "-d", dbName, "-t", "-c", orderQuery)
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	if output, err := cmd.CombinedOutput(); err == nil {
		orderColumn = strings.TrimSpace(string(output))
	}

	// Build the data query
	dataQuery := fmt.Sprintf(`SELECT * FROM %s.%s`, schema, table)
	if orderColumn != "" {
		dataQuery += fmt.Sprintf(` ORDER BY %s DESC`, orderColumn)
	}
	dataQuery += fmt.Sprintf(` LIMIT %d;`, limit)

	cmd = exec.Command("psql", "-h", host, "-p", port, "-U", user, "-d", dbName, "-x", "-c", dataQuery)
	cmd.Env = os.Environ()
	if password != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}

	output, err = cmd.CombinedOutput()
	if err != nil {
		// If expanded display fails, try normal display
		cmd = exec.Command("psql", "-h", host, "-p", port, "-U", user, "-d", dbName, "-c", dataQuery)
		cmd.Env = os.Environ()
		if password != "" {
			cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
		}
		output, err = cmd.CombinedOutput()
		if err != nil {
			color.Red("Error getting table data: %v", err)
			return
		}
	}

	fmt.Println(string(output))
}

func GetDBCommand() *cobra.Command {
	return dbCmd
}
