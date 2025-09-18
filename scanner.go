package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// findNpmProjects searches for all NPM projects in the given directory recursively
func findNpmProjects(rootDir string) ([]string, error) {
	infoColor.Printf("🔍 Searching for NPM projects in %s...\n", rootDir)

	var projects []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// node_modulesディレクトリはスキップ
		if info.IsDir() && info.Name() == "node_modules" {
			return filepath.SkipDir
		}

		// package.jsonファイルを見つけた場合
		if info.Name() == "package.json" {
			projectDir := filepath.Dir(path)

			// node_modules内のpackage.jsonは除外
			if !strings.Contains(projectDir, "node_modules") {
				projects = append(projects, projectDir)
				infoColor.Printf("  📁 Found: %s\n", projectDir)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	infoColor.Printf("🎯 Found %d NPM project(s)\n\n", len(projects))
	return projects, nil
}

// showProjectsAndConfirm displays the list of projects and asks for user confirmation
func showProjectsAndConfirm(projects []string) bool {
	infoColor.Println("📋 NPM Projects to be scanned:")
	for i, project := range projects {
		fmt.Printf("  %d. %s\n", i+1, project)
	}
	fmt.Println()

	return askForConfirmation("Do you want to proceed with the security scan?")
}

// scanProjects performs security scan on all given projects
func scanProjects(projects []string) error {
	// Initialize report
	initReport()
	setSafeChainMode(isSafeChainAvailable())

	infoColor.Printf("🚀 Starting security scan for %d project(s)...\n\n", len(projects))

	for i, project := range projects {
		infoColor.Printf("📦 [%d/%d] Processing: %s\n", i+1, len(projects), project)

		// Create result for this project
		result := ScanResult{
			ProjectPath: project,
			StartTime:   time.Now(),
			Status:      StatusInProgress,
		}

		// Step 1: node_modules削除
		if err := removeNodeModules(project); err != nil {
			errorColor.Printf("❌ Failed to remove node_modules in %s: %v\n", project, err)
			result.NodeModules.Success = false
			result.NodeModules.Error = err.Error()
			result.Status = StatusFailed
		} else {
			result.NodeModules.Success = true
		}

		// Step 2: npm install実行（node_modules削除が成功した場合のみ）
		if result.NodeModules.Success {
			if err := runNpmInstall(project); err != nil {
				errorColor.Printf("❌ Failed to run npm install in %s: %v\n", project, err)
				result.NpmInstall.Success = false
				result.NpmInstall.Error = err.Error()
				result.Status = StatusFailed
			} else {
				result.NpmInstall.Success = true
			}
		}

		// Step 3: セキュリティスキャン実行（npm installが成功した場合のみ）
		if result.NpmInstall.Success {
			if err := runSecurityScan(project, &result); err != nil {
				errorColor.Printf("❌ Failed to run security scan in %s: %v\n", project, err)
				result.SecurityScan.Success = false
				result.SecurityScan.Error = err.Error()
				result.Status = StatusFailed
			} else {
				result.Status = StatusSuccess
			}
		}

		// Finalize result
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)

		// Add to report
		addProjectResult(&result)

		if result.Status == StatusSuccess {
			successColor.Printf("✅ [%d/%d] Completed: %s\n\n", i+1, len(projects), project)
		} else {
			errorColor.Printf("❌ [%d/%d] Failed: %s\n\n", i+1, len(projects), project)
		}
	}

	// Finalize and show report
	finalizeReport()
	showScanResults()

	return nil
}

// removeNodeModules removes the node_modules directory in the given project directory
func removeNodeModules(projectDir string) error {
	nodeModulesPath := filepath.Join(projectDir, "node_modules")

	// node_modulesディレクトリが存在するかチェック
	if _, err := os.Stat(nodeModulesPath); os.IsNotExist(err) {
		infoColor.Printf("  📂 node_modules not found in %s (skipping)\n", projectDir)
		return nil
	}

	infoColor.Printf("  🗑️  Removing node_modules in %s...\n", projectDir)

	if err := os.RemoveAll(nodeModulesPath); err != nil {
		return fmt.Errorf("failed to remove node_modules: %w", err)
	}

	successColor.Printf("  ✅ node_modules removed from %s\n", projectDir)
	return nil
}

// runNpmInstall executes npm install in the given project directory
func runNpmInstall(projectDir string) error {
	infoColor.Printf("  📦 Running npm install in %s...\n", projectDir)

	cmd := exec.Command("npm", "install")
	cmd.Dir = projectDir

	// 出力をキャプチャ
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm install failed: %w\nOutput: %s", err, string(output))
	}

	successColor.Printf("  ✅ npm install completed in %s\n", projectDir)
	return nil
}

// isSafeChainAvailable checks if Safe Chain is available and properly set up
func isSafeChainAvailable() bool {
	if err := checkCommand("safe-chain"); err != nil {
		return false
	}
	return isSafeChainSetupComplete()
}

// showScanResults displays the final scan results and generates reports
func showScanResults() {
	if currentReport == nil {
		return
	}

	// Show terminal report
	printTerminalReport()

	// Always generate reports automatically
	infoColor.Println("📄 Generating reports...")

	var htmlReportPath string
	if err := generateHTMLReport(); err != nil {
		errorColor.Printf("❌ Failed to generate HTML report: %v\n", err)
	} else {
		htmlReportPath = fmt.Sprintf("%s/%s%s", ReportsDirName, currentReport.ScanID, HTMLExtension)
	}

	if err := generateJSONReport(); err != nil {
		errorColor.Printf("❌ Failed to generate JSON report: %v\n", err)
	}

	successColor.Println("✅ Reports generated successfully!")

	// Show report links
	if htmlReportPath != "" {
		showReportLinks(htmlReportPath)
	}
}

// showReportLinks displays clickable links to the generated reports
func showReportLinks(htmlReportPath string) {
	fmt.Printf("\n")
	infoColor.Println("📄 Generated Reports:")

	// Get absolute path for better user experience
	absPath, err := filepath.Abs(htmlReportPath)
	if err != nil {
		absPath = htmlReportPath
	}

	// Show clickable file:// URL for HTML report
	fileURL := fmt.Sprintf("file://%s", absPath)
	infoColor.Printf("🌐 HTML Report: ")
	successColor.Printf("%s\n", fileURL)

	// Show local file path
	infoColor.Printf("📁 Local Path: ")
	successColor.Printf("%s\n", absPath)

	// Show JSON report path
	jsonPath := strings.Replace(absPath, HTMLExtension, JSONExtension, 1)
	infoColor.Printf("📋 JSON Report: ")
	successColor.Printf("%s\n", jsonPath)

	fmt.Printf("\n")
	infoColor.Println("💡 Tips:")
	infoColor.Println("  - Click the file:// URL above to open in browser")
	infoColor.Println("  - Copy the local path to open manually")
	infoColor.Println("  - Use JSON report for automation/scripting")
}
