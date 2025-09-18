package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// ScanResult represents the result of scanning a single project
type ScanResult struct {
	StartTime       time.Time       `json:"start_time"`
	EndTime         time.Time       `json:"end_time"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	ProjectPath     string          `json:"project_path"`
	Status          string          `json:"status"`
	NodeModules     ActionResult    `json:"node_modules"`
	NpmInstall      ActionResult    `json:"npm_install"`
	SecurityScan    ActionResult    `json:"security_scan"`
	Duration        time.Duration   `json:"duration"`
}

// ActionResult represents the result of a specific action
type ActionResult struct {
	Error   string `json:"error"`
	Output  string `json:"output"`
	Success bool   `json:"success"`
}

// Vulnerability represents a security vulnerability found during scan
type Vulnerability struct {
	Severity    string `json:"severity"`
	Package     string `json:"package"`
	Description string `json:"description"`
	Fixed       bool   `json:"fixed"`
}

// ScanReport represents the complete scan report
type ScanReport struct {
	StartTime       time.Time     `json:"start_time"`
	EndTime         time.Time     `json:"end_time"`
	TotalDuration   time.Duration `json:"total_duration"`
	Results         []ScanResult  `json:"results"`
	ScanID          string        `json:"scan_id"`
	ProjectsScanned int           `json:"projects_scanned"`
	SuccessCount    int           `json:"success_count"`
	ErrorCount      int           `json:"error_count"`
	SafeChainMode   bool          `json:"safe_chain_mode"`
}

var currentReport *ScanReport

// initReport initializes a new scan report
func initReport() {
	currentReport = &ScanReport{
		ScanID:        fmt.Sprintf("scan_%d", time.Now().Unix()),
		StartTime:     time.Now(),
		Results:       make([]ScanResult, 0),
		SafeChainMode: false,
	}
}

// setSafeChainMode sets whether Safe Chain is being used
func setSafeChainMode(enabled bool) {
	if currentReport != nil {
		currentReport.SafeChainMode = enabled
	}
}

// addProjectResult adds a project scan result to the report
func addProjectResult(result *ScanResult) {
	if currentReport == nil {
		return
	}

	currentReport.Results = append(currentReport.Results, *result)
	currentReport.ProjectsScanned++

	if result.Status == StatusSuccess {
		currentReport.SuccessCount++
	} else {
		currentReport.ErrorCount++
	}
}

// finalizeReport completes the scan report
func finalizeReport() {
	if currentReport == nil {
		return
	}

	currentReport.EndTime = time.Now()
	currentReport.TotalDuration = currentReport.EndTime.Sub(currentReport.StartTime)
}

// printTerminalReport prints the scan report to terminal
func printTerminalReport() {
	if currentReport == nil {
		warningColor.Println("⚠️  No scan report available")
		return
	}

	printReportHeader()
	printReportSummary()
	printProjectResults()
	fmt.Println(strings.Repeat("=", ReportSeparator))
}

// printReportHeader prints the report header
func printReportHeader() {
	fmt.Println("\n" + strings.Repeat("=", ReportSeparator))
	infoColor.Printf("📊 SCAN REPORT - %s\n", currentReport.ScanID)
	fmt.Println(strings.Repeat("=", ReportSeparator))
}

// printReportSummary prints the report summary statistics
func printReportSummary() {
	infoColor.Printf("⏱️  Total Duration: %v\n", currentReport.TotalDuration.Round(time.Second))
	infoColor.Printf("📁 Projects Scanned: %d\n", currentReport.ProjectsScanned)
	successColor.Printf("✅ Successful: %d\n", currentReport.SuccessCount)
	if currentReport.ErrorCount > 0 {
		errorColor.Printf("❌ Failed: %d\n", currentReport.ErrorCount)
	}
	infoColor.Printf("🔒 Safe Chain Mode: %v\n", currentReport.SafeChainMode)
	fmt.Println()
}

// printProjectResults prints detailed results for each project
func printProjectResults() {
	for i := range currentReport.Results {
		result := &currentReport.Results[i]
		printProjectHeader(i+1, result)
		printProjectActions(result)
		printProjectVulnerabilities(result)
		fmt.Println()
	}
}

// printProjectHeader prints project basic info
func printProjectHeader(index int, result *ScanResult) {
	fmt.Printf("📦 [%d/%d] %s\n", index, len(currentReport.Results), result.ProjectPath)
	fmt.Printf("    Status: ")
	if result.Status == StatusSuccess {
		successColor.Printf("✅ Success")
	} else {
		errorColor.Printf("❌ %s", result.Status)
	}
	fmt.Printf(" (Duration: %v)\n", result.Duration.Round(time.Second))
}

// printProjectActions prints project action results
func printProjectActions(result *ScanResult) {
	printActionResult("🗑️  Node Modules", result.NodeModules, "Removed")
	printActionResult("📦 NPM Install", result.NpmInstall, "Success")
	printActionResult("🔍 Security Scan", result.SecurityScan, "Completed")
}

// printActionResult prints a single action result
func printActionResult(name string, action ActionResult, successMsg string) {
	if action.Success {
		fmt.Printf("    %s: ✅ %s\n", name, successMsg)
	} else if action.Error != "" {
		fmt.Printf("    %s: ❌ %s\n", name, action.Error)
	}
}

// printProjectVulnerabilities prints vulnerability information
func printProjectVulnerabilities(result *ScanResult) {
	if len(result.Vulnerabilities) > 0 {
		fmt.Printf("    🚨 Vulnerabilities: %d found\n", len(result.Vulnerabilities))
		for _, vuln := range result.Vulnerabilities {
			printSingleVulnerability(vuln)
		}
	} else if result.SecurityScan.Success {
		fmt.Printf("    🛡️  No vulnerabilities detected\n")
	}
}

// printSingleVulnerability prints a single vulnerability
func printSingleVulnerability(vuln Vulnerability) {
	severityColor := getSeverityColor(vuln.Severity)
	severityColor.Printf("      - %s: %s (%s)", vuln.Severity, vuln.Package, vuln.Description)
	if vuln.Fixed {
		successColor.Printf(" - FIXED")
	}
	fmt.Println()
}

// getSeverityColor returns appropriate color for vulnerability severity
func getSeverityColor(severity string) *color.Color {
	switch severity {
	case SeverityHigh, SeverityCritical:
		return errorColor
	case SeverityLow:
		return infoColor
	default:
		return warningColor
	}
}

// generateHTMLReport generates an HTML report file
func generateHTMLReport() error {
	if currentReport == nil {
		return fmt.Errorf("no scan report available")
	}

	htmlContent := generateHTMLContent()

	reportDir := ReportsDirName
	if err := os.MkdirAll(reportDir, DirPermSecure); err != nil {
		return fmt.Errorf("failed to create reports directory: %w", err)
	}

	filename := fmt.Sprintf("%s/%s.html", reportDir, currentReport.ScanID)
	if err := os.WriteFile(filename, []byte(htmlContent), FilePermSecure); err != nil {
		return fmt.Errorf("failed to write HTML report: %w", err)
	}

	successColor.Printf("📄 HTML Report generated: %s\n", filename)
	return nil
}

// generateJSONReport generates a JSON report file
func generateJSONReport() error {
	if currentReport == nil {
		return fmt.Errorf("no scan report available")
	}

	jsonData, err := json.MarshalIndent(currentReport, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON report: %w", err)
	}

	reportDir := ReportsDirName
	if err := os.MkdirAll(reportDir, DirPermSecure); err != nil {
		return fmt.Errorf("failed to create reports directory: %w", err)
	}

	filename := fmt.Sprintf("%s/%s.json", reportDir, currentReport.ScanID)
	if err := os.WriteFile(filename, jsonData, FilePermSecure); err != nil {
		return fmt.Errorf("failed to write JSON report: %w", err)
	}

	successColor.Printf("📄 JSON Report generated: %s\n", filename)
	return nil
}

// generateHTMLContent generates the HTML content for the report using Bulma CSS
func generateHTMLContent() string {
	return generateBulmaHTMLReport()
}

// getSafeChainTagClass returns the appropriate Bulma tag class for Safe Chain status
func getSafeChainTagClass() string {
	if currentReport.SafeChainMode {
		return BulmaSuccess
	}
	return BulmaWarning
}

// getTotalVulnerabilities counts total vulnerabilities across all projects
func getTotalVulnerabilities() int {
	total := 0
	for i := range currentReport.Results {
		total += len(currentReport.Results[i].Vulnerabilities)
	}
	return total
}

// generateBulmaProjectsHTML generates HTML for all projects using Bulma components
func generateBulmaProjectsHTML() string {
	html := ""

	for i := range currentReport.Results {
		html += generateBulmaProjectCard(i+1, &currentReport.Results[i])
	}

	return html
}

// generateBulmaProjectCard generates a Bulma card for a single project
func generateBulmaProjectCard(index int, result *ScanResult) string {
	statusClass, statusIcon := getProjectStatusStyle(result)
	return fmt.Sprintf(`
        <div class="card project-card">
            %s
            <div class="card-content">
                <div class="content">
                    %s
                    <div class="field is-grouped is-grouped-multiline">
                        %s
                    </div>
                    %s
                </div>
            </div>
        </div>`,
		generateProjectCardHeader(index, result, statusClass, statusIcon),
		generateProjectCardMeta(result),
		generateBulmaActionsHTML(result),
		generateBulmaVulnerabilitiesHTML(result.Vulnerabilities, result.SecurityScan.Success))
}

// generateBulmaActionsHTML generates Bulma tags for each action
func generateBulmaActionsHTML(result *ScanResult) string {
	actions := []struct {
		icon   string
		name   string
		action ActionResult
	}{
		{"fas fa-trash", "Node Modules", result.NodeModules},
		{"fas fa-download", "NPM Install", result.NpmInstall},
		{"fas fa-search", "Security Scan", result.SecurityScan},
	}

	html := ""
	for _, action := range actions {
		tagClass := BulmaSuccess
		icon := "fas fa-check"
		if !action.action.Success {
			tagClass = BulmaDanger
			icon = "fas fa-times"
		}

		html += fmt.Sprintf(`
                        <div class="control">
                            <div class="tags has-addons">
                                <span class="tag is-dark">
                                    <i class="%s"></i>&nbsp; %s
                                </span>
                                <span class="tag %s">
                                    <i class="%s"></i>
                                </span>
                            </div>
                        </div>`,
			action.icon, action.name, tagClass, icon)
	}

	return html
}

// generateBulmaVulnerabilitiesHTML generates HTML for vulnerabilities using Bulma
func generateBulmaVulnerabilitiesHTML(vulnerabilities []Vulnerability, scanSuccess bool) string {
	if len(vulnerabilities) == 0 {
		if scanSuccess {
			return `
                    <div class="notification is-success is-light">
                        <i class="fas fa-shield-alt"></i>&nbsp;
                        <strong>No vulnerabilities detected</strong> - Your project is secure!
                    </div>`
		}
		return ""
	}

	html := fmt.Sprintf(`
                    <div class="field">
                        <label class="label">
                            <i class="fas fa-bug"></i>&nbsp; 
                            Security Vulnerabilities (%d found)
                        </label>`, len(vulnerabilities))

	for _, vuln := range vulnerabilities {
		// Determine severity styling
		severityClass, severityIcon := getVulnerabilitySeverityStyle(vuln)

		html += fmt.Sprintf(`
                        <div class="vulnerability-item %s">
                            <div class="level">
                                <div class="level-left">
                                    <div class="level-item">
                                        <span class="tag %s is-medium">
                                            <i class="%s"></i>&nbsp;
                                            %s
                                        </span>
                                    </div>
                                    <div class="level-item">
                                        <div>
                                            <p class="has-text-weight-bold">%s</p>
                                            <p class="is-size-7 has-text-grey">%s</p>
                                        </div>
                                    </div>
                                </div>
                                <div class="level-right">
                                    %s
                                </div>
                            </div>
                        </div>`,
			getVulnBgClass(vuln),
			severityClass,
			severityIcon,
			strings.ToUpper(vuln.Severity),
			vuln.Package,
			vuln.Description,
			getFixedBadgeHTML(vuln.Fixed))
	}

	html += `
                    </div>`

	return html
}

// getVulnBgClass returns the background class for vulnerability severity
func getVulnBgClass(vuln Vulnerability) string {
	if vuln.Fixed {
		return "vuln-fixed"
	}

	switch vuln.Severity {
	case SeverityCritical:
		return "vuln-critical"
	case SeverityHigh:
		return "vuln-high"
	case SeverityModerate:
		return "vuln-moderate"
	case SeverityLow:
		return "vuln-low"
	default:
		return "vuln-moderate"
	}
}

// getVulnerabilitySeverityStyle returns the appropriate Bulma classes for vulnerability severity
func getVulnerabilitySeverityStyle(vuln Vulnerability) (string, string) {
	if vuln.Fixed {
		return BulmaSuccess, "fas fa-check-circle"
	}

	switch vuln.Severity {
	case SeverityCritical:
		return BulmaDanger, "fas fa-skull-crossbones"
	case SeverityHigh:
		return BulmaDanger, "fas fa-exclamation-circle"
	case SeverityModerate:
		return BulmaWarning, "fas fa-exclamation-triangle"
	case SeverityLow:
		return BulmaInfo, "fas fa-info-circle"
	default:
		return BulmaWarning, "fas fa-exclamation-triangle"
	}
}

// getFixedBadgeHTML returns HTML for fixed status badge
func getFixedBadgeHTML(fixed bool) string {
	if fixed {
		return fmt.Sprintf(`<span class="tag %s">
                                        <i class="fas fa-wrench"></i>&nbsp; FIXED
                                    </span>`, BulmaSuccess)
	}
	return `<span class="tag is-light">
                                    <i class="fas fa-clock"></i>&nbsp; PENDING
                                </span>`
}

// generateBulmaHTMLReport generates a modern HTML report using Bulma CSS framework
func generateBulmaHTMLReport() string {
	return fmt.Sprintf(`<!DOCTYPE html>
%s
<body>
%s
%s
<section class="section">
    <div class="container">
        %s
    </div>
</section>
%s
</body>
</html>`,
		generateBulmaHTMLHead(),
		generateBulmaHeroSection(),
		generateBulmaStatsSection(),
		generateBulmaProjectsHTML(),
		generateBulmaFooter())
}

// generateBulmaHTMLHead generates HTML head section
func generateBulmaHTMLHead() string {
	return fmt.Sprintf(`<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NPM Security Scanner Report - %s</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.min.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
    <style>%s</style>
</head>`, currentReport.ScanID, generateBulmaCSS())
}

// generateBulmaCSS generates custom CSS styles
func generateBulmaCSS() string {
	return `
        .hero-gradient { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); }
        .vulnerability-item { border-left: 4px solid; border-radius: 6px; padding: 1rem; margin-bottom: 0.75rem; 
                              transition: all 0.2s ease; }
        .vulnerability-item:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
        .vuln-critical { border-left-color: #ff3860; background: linear-gradient(to right, #feecf0, #fef7f7); }
        .vuln-high { border-left-color: #ff6348; background: linear-gradient(to right, #fef0ef, #fef7f7); }
        .vuln-moderate { border-left-color: #ffdd57; background: linear-gradient(to right, #fffbeb, #fffef7); }
        .vuln-low { border-left-color: #209cee; background: linear-gradient(to right, #eff5ff, #f7faff); }
        .vuln-fixed { border-left-color: #23d160; background: linear-gradient(to right, #f0fff4, #f7fef9); }
        .stats-card { transition: transform 0.2s ease; }
        .stats-card:hover { transform: translateY(-4px); }
        .scan-meta { background: rgba(255,255,255,0.1); border-radius: 8px; padding: 1rem; margin-top: 1rem; }
        .project-card { margin-bottom: 1.5rem; }`
}

// generateBulmaHeroSection generates hero section HTML
func generateBulmaHeroSection() string {
	return fmt.Sprintf(`<section class="hero is-large hero-gradient">
        <div class="hero-body">
            <div class="container">
                <h1 class="title is-1 has-text-white">
                    <i class="fas fa-shield-virus"></i> NPM Security Scanner
                </h1>
                <h2 class="subtitle is-3 has-text-white-ter">
                    Comprehensive malware protection for your Node.js projects
                </h2>
                <div class="scan-meta">
                    <div class="level">
                        <div class="level-left">%s</div>
                        <div class="level-right">
                            <div class="level-item">
                                <span class="tag is-large %s">
                                    <i class="fas fa-link"></i>&nbsp; Safe Chain: %v
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>`, generateMetaItems(), getSafeChainTagClass(), currentReport.SafeChainMode)
}

// generateMetaItems generates metadata items for hero section
func generateMetaItems() string {
	return fmt.Sprintf(`<div class="level-item">
                                <div>
                                    <p class="heading has-text-white-ter">Report ID</p>
                                    <p class="title is-5 has-text-white">%s</p>
                                </div>
                            </div>
                            <div class="level-item">
                                <div>
                                    <p class="heading has-text-white-ter">Generated</p>
                                    <p class="title is-5 has-text-white">%s</p>
                                </div>
                            </div>
                            <div class="level-item">
                                <div>
                                    <p class="heading has-text-white-ter">Total Duration</p>
                                    <p class="title is-5 has-text-white">%v</p>
                                </div>
                            </div>`,
		currentReport.ScanID,
		currentReport.EndTime.Format("2006-01-02 15:04:05"),
		currentReport.TotalDuration.Round(time.Second))
}

// generateBulmaStatsSection generates statistics cards section
func generateBulmaStatsSection() string {
	return fmt.Sprintf(`<section class="section">
        <div class="container">
            <div class="columns">%s</div>
        </div>
    </section>`, generateStatsCards())
}

// generateStatsCards generates individual stat cards
func generateStatsCards() string {
	return fmt.Sprintf(`
                <div class="column is-3">
                    <div class="card stats-card">
                        <div class="card-content has-text-centered">
                            <div class="content">
                                <i class="fas fa-folder-open fa-3x has-text-info mb-3"></i>
                                <p class="title is-1 has-text-info">%d</p>
                                <p class="subtitle is-5">Projects Scanned</p>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="column is-3">
                    <div class="card stats-card">
                        <div class="card-content has-text-centered">
                            <div class="content">
                                <i class="fas fa-check-circle fa-3x has-text-success mb-3"></i>
                                <p class="title is-1 has-text-success">%d</p>
                                <p class="subtitle is-5">Successful</p>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="column is-3">
                    <div class="card stats-card">
                        <div class="card-content has-text-centered">
                            <div class="content">
                                <i class="fas fa-times-circle fa-3x has-text-danger mb-3"></i>
                                <p class="title is-1 has-text-danger">%d</p>
                                <p class="subtitle is-5">Failed</p>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="column is-3">
                    <div class="card stats-card">
                        <div class="card-content has-text-centered">
                            <div class="content">
                                <i class="fas fa-bug fa-3x has-text-warning mb-3"></i>
                                <p class="title is-1 has-text-warning">%d</p>
                                <p class="subtitle is-5">Total Vulnerabilities</p>
                            </div>
                        </div>
                    </div>
                </div>`,
		currentReport.ProjectsScanned,
		currentReport.SuccessCount,
		currentReport.ErrorCount,
		getTotalVulnerabilities())
}

// generateBulmaFooter generates footer section
func generateBulmaFooter() string {
	return fmt.Sprintf(`<footer class="footer">
    <div class="content has-text-centered">
        <p>
            <strong>NPM Security Scanner v1.3.0</strong><br>
            Generated at %s<br>
            <i class="fas fa-shield-alt"></i> Protecting your dependencies from malware
        </p>
    </div>
</footer>`, time.Now().Format("2006-01-02 15:04:05"))
}

// getProjectStatusStyle returns status class and icon for project
func getProjectStatusStyle(result *ScanResult) (string, string) {
	if result.Status == StatusSuccess {
		return BulmaSuccess, "fas fa-check-circle"
	}
	return BulmaDanger, "fas fa-times-circle"
}

// generateProjectCardHeader generates card header for project
func generateProjectCardHeader(index int, result *ScanResult, statusClass, statusIcon string) string {
	return fmt.Sprintf(`<header class="card-header">
                <p class="card-header-title is-size-4">
                    <i class="fas fa-folder-open"></i>&nbsp;
                    %d. %s
                </p>
                <div class="card-header-icon">
                    <span class="tag %s is-medium">
                        <i class="%s"></i>&nbsp;
                        %s
                    </span>
                </div>
            </header>`, index, result.ProjectPath, statusClass, statusIcon, result.Status)
}

// generateProjectCardMeta generates project metadata section
func generateProjectCardMeta(result *ScanResult) string {
	return fmt.Sprintf(`<div class="level">
                        <div class="level-left">
                            <div class="level-item">
                                <div>
                                    <p class="heading">Duration</p>
                                    <p class="title is-6">%v</p>
                                </div>
                            </div>
                            <div class="level-item">
                                <div>
                                    <p class="heading">Started</p>
                                    <p class="title is-6">%s</p>
                                </div>
                            </div>
                        </div>
                    </div>`, result.Duration.Round(time.Second), result.StartTime.Format("15:04:05"))
}
