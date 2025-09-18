// Package main provides NPM security scanner for detecting malware in Node.js projects
package main

// Status constants
const (
	StatusSuccess    = "success"
	StatusFailed     = "failed"
	StatusInProgress = "in_progress"
)

// HTML template constants
const (
	HTMLDivClose   = "</div>"
	HTMLActionCard = "action-card"
	HTMLVulnClass  = "vulnerability"
)

// File permissions (secure defaults)
const (
	FilePermSecure   = 0o600 // Owner read/write only
	DirPermSecure    = 0o755 // Owner full, group/other read/execute
	FilePermReadable = 0o644 // Owner read/write, others read
)

// Severity levels
const (
	SeverityLow      = "low"
	SeverityModerate = "moderate"
	SeverityHigh     = "high"
	SeverityCritical = "critical"
)

// Report constants
const (
	ReportsDirName    = "reports"
	HTMLExtension     = ".html"
	JSONExtension     = ".json"
	ReportSeparator   = 80
	VulnerabilityText = " vulnerabilit"
)

// Bulma CSS class constants
const (
	BulmaSuccess = "is-success"
	BulmaDanger  = "is-danger"
	BulmaWarning = "is-warning"
	BulmaInfo    = "is-info"
)
