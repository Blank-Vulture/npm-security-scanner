package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindNpmProjects(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// テストケース1: package.jsonがあるプロジェクト
	project1 := filepath.Join(tempDir, "project1")
	if err := os.MkdirAll(project1, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	packageJson1 := filepath.Join(project1, "package.json")
	if err := os.WriteFile(packageJson1, []byte(`{"name": "test-project1"}`), 0644); err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// テストケース2: ネストされたプロジェクト
	project2 := filepath.Join(tempDir, "nested", "project2")
	if err := os.MkdirAll(project2, 0755); err != nil {
		t.Fatalf("Failed to create nested test directory: %v", err)
	}

	packageJson2 := filepath.Join(project2, "package.json")
	if err := os.WriteFile(packageJson2, []byte(`{"name": "test-project2"}`), 0644); err != nil {
		t.Fatalf("Failed to create nested package.json: %v", err)
	}

	// テストケース3: package.jsonがないディレクトリ
	emptyDir := filepath.Join(tempDir, "empty")
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatalf("Failed to create empty directory: %v", err)
	}

	// テスト実行
	projects, err := findNpmProjects(tempDir)
	if err != nil {
		t.Fatalf("findNpmProjects failed: %v", err)
	}

	// 結果の検証
	expectedCount := 2
	if len(projects) != expectedCount {
		t.Errorf("Expected %d projects, got %d", expectedCount, len(projects))
	}

	// プロジェクトパスの検証
	foundProject1 := false
	foundProject2 := false

	for _, project := range projects {
		switch project {
		case project1:
			foundProject1 = true
		case project2:
			foundProject2 = true
		}
	}

	if !foundProject1 {
		t.Errorf("project1 not found in results")
	}
	if !foundProject2 {
		t.Errorf("project2 not found in results")
	}
}

func TestCheckCommand(t *testing.T) {
	// 存在するコマンドのテスト（lsはほぼすべてのUnix系システムに存在）
	if err := checkCommand("ls"); err != nil {
		t.Errorf("Expected 'ls' command to be available, got error: %v", err)
	}

	// 存在しないコマンドのテスト
	if err := checkCommand("non-existent-command-12345"); err == nil {
		t.Errorf("Expected error for non-existent command, but got nil")
	}
}

func TestIsValidPackageJson(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// 有効なpackage.jsonファイル
	validPackageJSON := filepath.Join(tempDir, "valid-package.json")
	if err := os.WriteFile(validPackageJSON, []byte(`{"name": "test"}`), 0644); err != nil {
		t.Fatalf("Failed to create valid package.json: %v", err)
	}

	// 存在しないファイル
	nonExistentFile := filepath.Join(tempDir, "non-existent.json")

	// テスト実行
	if !isValidPackageJSON(validPackageJSON) {
		t.Errorf("Expected valid package.json to return true")
	}

	if isValidPackageJSON(nonExistentFile) {
		t.Errorf("Expected non-existent file to return false")
	}
}

func TestGetAbsolutePath(t *testing.T) {
	// 相対パスのテスト
	relativePath := "."
	absPath, err := getAbsolutePath(relativePath)
	if err != nil {
		t.Errorf("getAbsolutePath failed for relative path: %v", err)
	}

	if !filepath.IsAbs(absPath) {
		t.Errorf("Expected absolute path, got: %s", absPath)
	}

	// 既に絶対パスの場合
	if filepath.IsAbs(relativePath) {
		return // このテストは相対パスでのみ有効
	}

	// 絶対パスのテスト
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	absPath2, err := getAbsolutePath(currentDir)
	if err != nil {
		t.Errorf("getAbsolutePath failed for absolute path: %v", err)
	}

	if absPath2 != currentDir {
		t.Errorf("Expected same path for absolute input, got: %s vs %s", absPath2, currentDir)
	}
}

func TestRemoveNodeModules(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// node_modulesディレクトリを作成
	nodeModulesDir := filepath.Join(tempDir, "node_modules")
	if err := os.MkdirAll(nodeModulesDir, 0755); err != nil {
		t.Fatalf("Failed to create node_modules directory: %v", err)
	}

	// テスト用ファイルを作成
	testFile := filepath.Join(nodeModulesDir, "test-package", "index.js")
	if err := os.MkdirAll(filepath.Dir(testFile), 0755); err != nil {
		t.Fatalf("Failed to create test package directory: %v", err)
	}
	if err := os.WriteFile(testFile, []byte("console.log('test');"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// node_modulesが存在することを確認
	if _, err := os.Stat(nodeModulesDir); os.IsNotExist(err) {
		t.Fatalf("node_modules directory should exist before test")
	}

	// removeNodeModules実行
	if err := removeNodeModules(tempDir); err != nil {
		t.Errorf("removeNodeModules failed: %v", err)
	}

	// node_modulesが削除されていることを確認
	if _, err := os.Stat(nodeModulesDir); !os.IsNotExist(err) {
		t.Errorf("node_modules directory should be removed")
	}

	// node_modulesが存在しない場合のテスト
	if err := removeNodeModules(tempDir); err != nil {
		t.Errorf("removeNodeModules should not fail when node_modules doesn't exist: %v", err)
	}
}

// ベンチマークテスト
func BenchmarkFindNpmProjects(b *testing.B) {
	// テスト用の一時ディレクトリを作成
	tempDir := b.TempDir()

	// 複数のプロジェクトを作成
	for i := 0; i < 10; i++ {
		projectDir := filepath.Join(tempDir, "project"+string(rune(i+'0')))
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			b.Fatalf("Failed to create test directory: %v", err)
		}

		packageJson := filepath.Join(projectDir, "package.json")
		if err := os.WriteFile(packageJson, []byte(`{"name": "test"}`), 0644); err != nil {
			b.Fatalf("Failed to create package.json: %v", err)
		}
	}

	// ベンチマーク実行
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := findNpmProjects(tempDir)
		if err != nil {
			b.Fatalf("findNpmProjects failed: %v", err)
		}
	}
}
