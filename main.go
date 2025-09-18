package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	appName    = "npm-security-scanner"
	appVersion = "1.3.0"
)

var (
	// カラー出力用
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	warningColor = color.New(color.FgYellow, color.Bold)
	infoColor    = color.New(color.FgCyan, color.Bold)
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   appName + " [target-directory]",
		Short: "NPM Security Scanner using Safe Chain",
		Long: `NPMパッケージのマルウェア感染対策のためのセキュリティスキャナーツール。
指定されたディレクトリ配下のすべてのNPMプロジェクトを再帰的に検索し、
Safe Chainを使用して一括でセキュリティスキャンを実行します。`,
		Version: appVersion,
		Args:    cobra.MaximumNArgs(1),
		Run:     runScanner,
	}

	if err := rootCmd.Execute(); err != nil {
		errorColor.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runScanner(_ *cobra.Command, args []string) {
	// ターゲットディレクトリの決定
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	infoColor.Printf("🔍 NPM Security Scanner v%s\n", appVersion)
	infoColor.Printf("Target directory: %s\n\n", targetDir)

	// Step 1: Safe Chainのインストール確認
	if err := checkSafeChainInstallation(); err != nil {
		if err.Error() == "terminal restart required" {
			infoColor.Println("👋 Exiting for terminal restart...")
			os.Exit(0)
		}
		errorColor.Printf("❌ Safe Chain setup failed: %v\n", err)
		os.Exit(1)
	}

	// Step 2: NPMプロジェクトの検索
	projects, err := findNpmProjects(targetDir)
	if err != nil {
		errorColor.Printf("❌ Failed to find NPM projects: %v\n", err)
		os.Exit(1)
	}

	if len(projects) == 0 {
		warningColor.Println("⚠️  No NPM projects found in the specified directory")
		return
	}

	// Step 3: プロジェクトリストの表示と確認
	if !showProjectsAndConfirm(projects) {
		infoColor.Println("🚫 Scan canceled by user")
		return
	}

	// Step 4: スキャン実行
	if err := scanProjects(projects); err != nil {
		errorColor.Printf("❌ Scan failed: %v\n", err)
		os.Exit(1)
	}

	successColor.Println("✅ All projects scanned successfully!")
}

// askForConfirmation prompts the user for yes/no confirmation
func askForConfirmation(question string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		infoColor.Printf("%s [y/N]: ", question)
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v", err)
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))

		switch response {
		case "y", "yes":
			return true
		case "n", "no", "":
			return false
		default:
			warningColor.Println("Please answer 'y' or 'n'")
		}
	}
}
