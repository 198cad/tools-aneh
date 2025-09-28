package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update tools to the latest version",
	Long:  `Check for updates and install the latest version of tools`,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		check, _ := cmd.Flags().GetBool("check")

		if check {
			checkForUpdates()
		} else {
			performUpdate(force)
		}
	},
}

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update the tools binary itself",
	Long:  `Download and install the latest version of the tools binary`,
	Run: func(cmd *cobra.Command, args []string) {
		performSelfUpdate()
	},
}

func init() {
	updateCmd.Flags().BoolP("force", "f", false, "Force update even if up to date")
	updateCmd.Flags().BoolP("check", "c", false, "Only check for updates, don't install")
}

func checkForUpdates() {
	color.Yellow("Checking for updates...")

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		color.Red("Error getting executable path: %v", err)
		return
	}

	installDir := filepath.Dir(execPath)

	// Check if it's a git repository
	gitDir := filepath.Join(installDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		color.Yellow("Not a git repository. Cannot check for updates.")
		color.Yellow("To enable updates, clone the repository:")
		fmt.Printf("  git clone <repository-url> %s\n", installDir)
		return
	}

	// Fetch latest changes
	fetchCmd := exec.Command("git", "fetch", "origin")
	fetchCmd.Dir = installDir
	if output, err := fetchCmd.CombinedOutput(); err != nil {
		color.Red("Error fetching updates: %v", err)
		if len(output) > 0 {
			fmt.Println(string(output))
		}
		return
	}

	// Check if updates are available
	statusCmd := exec.Command("git", "status", "-uno")
	statusCmd.Dir = installDir
	output, err := statusCmd.CombinedOutput()
	if err != nil {
		color.Red("Error checking status: %v", err)
		return
	}

	if strings.Contains(string(output), "Your branch is behind") {
		color.Green("✓ Updates available!")
		fmt.Println("Run 'tools update' to install the latest version")

		// Show what's new
		logCmd := exec.Command("git", "log", "--oneline", "HEAD..origin/main", "--max-count=10")
		logCmd.Dir = installDir
		if changes, err := logCmd.CombinedOutput(); err == nil && len(changes) > 0 {
			color.Cyan("\nRecent changes:")
			fmt.Println(string(changes))
		}
	} else if strings.Contains(string(output), "Your branch is up to date") {
		color.Green("✓ You are running the latest version!")
	} else {
		fmt.Println(string(output))
	}
}

func performUpdate(force bool) {
	color.Yellow("Starting update process...")

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		color.Red("Error getting executable path: %v", err)
		return
	}

	installDir := filepath.Dir(execPath)

	// Check if it's a git repository
	gitDir := filepath.Join(installDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		// Try to initialize git and add remote
		color.Yellow("Initializing git repository...")
		if err := initGitRepo(installDir); err != nil {
			color.Red("Failed to initialize git repository: %v", err)
			performSelfUpdate()
			return
		}
	}

	// Save current version info
	oldVersionFile := filepath.Join(installDir, "version.go")
	oldVersion := getFileVersion(oldVersionFile)

	// Stash any local changes
	color.Yellow("Saving local changes...")
	stashCmd := exec.Command("git", "stash", "save", fmt.Sprintf("Auto-stash before update at %s", time.Now().Format("2006-01-02 15:04:05")))
	stashCmd.Dir = installDir
	stashOutput, _ := stashCmd.CombinedOutput()
	hasStash := !strings.Contains(string(stashOutput), "No local changes")

	// Pull latest changes
	color.Yellow("Downloading updates...")
	pullCmd := exec.Command("git", "pull", "origin", "main")
	pullCmd.Dir = installDir
	pullOutput, err := pullCmd.CombinedOutput()
	if err != nil {
		// Try to handle merge conflicts
		if strings.Contains(string(pullOutput), "conflict") {
			color.Yellow("Merge conflicts detected. Resolving...")

			// Reset to origin/main
			resetCmd := exec.Command("git", "reset", "--hard", "origin/main")
			resetCmd.Dir = installDir
			if _, err := resetCmd.CombinedOutput(); err != nil {
				color.Red("Error resolving conflicts: %v", err)
				return
			}
			color.Green("✓ Conflicts resolved by using latest version")
		} else {
			color.Red("Error pulling updates: %v", err)
			fmt.Println(string(pullOutput))
			return
		}
	} else {
		fmt.Println(string(pullOutput))
	}

	// Check if version changed
	newVersion := getFileVersion(oldVersionFile)
	if !force && oldVersion == newVersion {
		color.Green("✓ Already up to date!")

		// Restore stashed changes if any
		if hasStash {
			color.Yellow("Restoring local changes...")
			popCmd := exec.Command("git", "stash", "pop")
			popCmd.Dir = installDir
			popCmd.CombinedOutput()
		}
		return
	}

	// Install dependencies
	color.Yellow("Installing dependencies...")
	modCmd := exec.Command("go", "mod", "download")
	modCmd.Dir = installDir
	if output, err := modCmd.CombinedOutput(); err != nil {
		color.Red("Error downloading dependencies: %v", err)
		if len(output) > 0 {
			fmt.Println(string(output))
		}
	}

	// Build new version
	color.Yellow("Building new version...")
	buildCmd := exec.Command("go", "build", "-o", "tools.exe.new", ".")
	buildCmd.Dir = installDir
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		color.Red("Error building new version: %v", err)
		fmt.Println(string(buildOutput))

		// Restore stashed changes if any
		if hasStash {
			popCmd := exec.Command("git", "stash", "pop")
			popCmd.Dir = installDir
			popCmd.CombinedOutput()
		}
		return
	}

	// Replace old executable
	newExe := filepath.Join(installDir, "tools.exe.new")
	oldExe := filepath.Join(installDir, "tools.exe")
	backupExe := filepath.Join(installDir, "tools.exe.backup")

	// Backup current executable
	if err := copyFile(oldExe, backupExe); err != nil {
		color.Red("Error creating backup: %v", err)
		return
	}

	// Windows-specific: Schedule replacement using batch file
	if runtime.GOOS == "windows" {
		if err := scheduleReplacement(oldExe, newExe, backupExe); err != nil {
			color.Red("Error scheduling update: %v", err)
			return
		}

		color.Green("✓ Update downloaded successfully!")
		color.Yellow("The update will be applied when you restart the tool.")
		fmt.Println("\nPlease close and restart the tools to complete the update.")
	} else {
		// Direct replacement for non-Windows systems
		if err := os.Rename(newExe, oldExe); err != nil {
			color.Red("Error replacing executable: %v", err)
			// Restore from backup
			os.Rename(backupExe, oldExe)
			return
		}

		color.Green("✓ Update completed successfully!")
		color.Green("Version: %s", newVersion)

		// Clean up backup
		os.Remove(backupExe)
	}

	// Restore stashed changes if any
	if hasStash {
		color.Yellow("Restoring local changes...")
		popCmd := exec.Command("git", "stash", "pop")
		popCmd.Dir = installDir
		if output, err := popCmd.CombinedOutput(); err != nil {
			color.Yellow("Could not restore local changes: %v", err)
			color.Yellow("Your changes are saved in git stash")
		} else if len(output) > 0 {
			fmt.Println(string(output))
		}
	}

	// Show what's new
	showChangelog(installDir)
}

func performSelfUpdate() {
	color.Yellow("Performing self-update...")

	// Get current executable
	execPath, err := os.Executable()
	if err != nil {
		color.Red("Error getting executable path: %v", err)
		return
	}

	installDir := filepath.Dir(execPath)

	// Create update script
	updateScript := filepath.Join(installDir, "self_update.bat")
	scriptContent := fmt.Sprintf(`@echo off
echo Updating tools...
timeout /t 2 /nobreak > nul
move /Y "%s\tools.exe.new" "%s\tools.exe"
echo Update completed!
del "%%~f0"
`, installDir, installDir)

	if err := os.WriteFile(updateScript, []byte(scriptContent), 0755); err != nil {
		color.Red("Error creating update script: %v", err)
		return
	}

	// Download new version (this would typically download from a release URL)
	color.Yellow("Building latest version...")
	buildCmd := exec.Command("go", "build", "-o", "tools.exe.new", ".")
	buildCmd.Dir = installDir
	if output, err := buildCmd.CombinedOutput(); err != nil {
		color.Red("Error building new version: %v", err)
		fmt.Println(string(output))
		return
	}

	// Execute update script
	cmd := exec.Command("cmd", "/c", "start", "/b", updateScript)
	if err := cmd.Start(); err != nil {
		color.Red("Error starting update: %v", err)
		return
	}

	color.Green("✓ Update in progress...")
	fmt.Println("The application will restart with the new version.")
	os.Exit(0)
}

func initGitRepo(dir string) error {
	// Initialize git
	initCmd := exec.Command("git", "init")
	initCmd.Dir = dir
	if _, err := initCmd.CombinedOutput(); err != nil {
		return err
	}

	// Add remote (you should replace with your actual repository URL)
	remoteCmd := exec.Command("git", "remote", "add", "origin", "https://github.com/yourusername/tools.git")
	remoteCmd.Dir = dir
	remoteCmd.CombinedOutput()

	// Fetch
	fetchCmd := exec.Command("git", "fetch", "origin")
	fetchCmd.Dir = dir
	fetchCmd.CombinedOutput()

	return nil
}

func getFileVersion(versionFile string) string {
	content, err := os.ReadFile(versionFile)
	if err != nil {
		return "unknown"
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Version") && strings.Contains(line, "=") {
			parts := strings.Split(line, "\"")
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	}

	return "unknown"
}

func copyFile(src, dst string) error {
	sourceFile, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, sourceFile, 0755)
}

func scheduleReplacement(oldExe, newExe, backupExe string) error {
	// Create a batch file to replace the executable
	batchFile := filepath.Join(filepath.Dir(oldExe), "update_tools.bat")
	batchContent := fmt.Sprintf(`@echo off
:WAIT
timeout /t 1 /nobreak > nul
tasklist | find /i "tools.exe" > nul
if errorlevel 1 goto :REPLACE
goto :WAIT

:REPLACE
move /Y "%s" "%s"
move /Y "%s" "%s"
echo Update completed successfully!
del "%s"
del "%%~f0"
`, backupExe, oldExe+".old", newExe, oldExe, oldExe+".old")

	if err := os.WriteFile(batchFile, []byte(batchContent), 0755); err != nil {
		return err
	}

	// Start the batch file in background
	cmd := exec.Command("cmd", "/c", "start", "/b", batchFile)
	return cmd.Start()
}

func showChangelog(dir string) {
	// Show recent commits
	logCmd := exec.Command("git", "log", "--oneline", "--max-count=5")
	logCmd.Dir = dir
	if output, err := logCmd.CombinedOutput(); err == nil && len(output) > 0 {
		color.Cyan("\nRecent changes:")
		fmt.Println(string(output))
	}
}

func GetUpdateCommand() *cobra.Command {
	return updateCmd
}
