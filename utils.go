package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// checkCommand checks if a command is available in PATH
func checkCommand(command string) error {
	_, err := exec.LookPath(command)
	return err
}

// isValidPackageJSON checks if the package.json file is valid and contains dependencies
func isValidPackageJSON(packageJSONPath string) bool {
	// 基本的な存在チェック
	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
		return false
	}

	// TODO: より詳細な検証を実装する場合
	// - JSONの妥当性チェック
	// - dependenciesやdevDependenciesの存在チェック
	// - npmプロジェクトとしての妥当性チェック

	return true
}

// getAbsolutePath returns the absolute path for the given path
func getAbsolutePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}
	return absPath, nil
}
