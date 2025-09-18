# NPM Security Scanner - 使用ガイド

## TL;DR

```bash
# ビルド
make build

# 実行（カレントディレクトリ）
./bin/npm-security-scanner

# 実行（指定ディレクトリ）
./bin/npm-security-scanner /path/to/projects
```

## 詳細な使用方法

### 1. 事前準備

#### Safe Chainのインストール（推奨）

```bash
# NPMでグローバルインストール
npm install -g safe-chain-test

# セットアップ実行
safe-chain setup

# ターミナル再起動

# インストール確認
safe-chain --version
```

**注意**: Safe Chainがインストールされていない場合、デモモードで動作します。

#### Goの依存関係インストール

```bash
cd npm-security-scanner
go mod download
```

### 2. ビルド

```bash
# Makefileを使用（推奨）
make build

# 直接ビルド
go build -o bin/npm-security-scanner .

# クロスプラットフォームビルド
make build-all
```

### 3. 基本的な使用方法

#### カレントディレクトリのスキャン

```bash
./bin/npm-security-scanner
```

#### 特定ディレクトリのスキャン

```bash
./bin/npm-security-scanner /Users/username/projects
./bin/npm-security-scanner ./my-projects
```

#### ヘルプ表示

```bash
./bin/npm-security-scanner --help
```

### 4. 実行フロー

1. **Safe Chain確認**
   - インストール状況をチェック
   - 未インストールの場合はデモモードで継続

2. **プロジェクト検索**
   - 指定ディレクトリを再帰的に検索
   - `package.json`を持つディレクトリを検出
   - `node_modules`内は除外

3. **確認プロンプト**
   - 検出されたプロジェクト一覧を表示
   - ユーザーに実行確認

4. **スキャン実行**
   - 各プロジェクトで`node_modules`を削除
   - `npm install`で依存関係を再インストール
   - Safe Chainでセキュリティスキャン

### 5. 実行例

```bash
$ ./bin/npm-security-scanner examples/
🔍 NPM Security Scanner v1.0.0
Target directory: examples/

🔧 Checking Safe Chain installation...
⚠️  Safe Chain is not installed globally
🔧 Running in demo mode without Safe Chain
📋 To install Safe Chain later:
   1. Run: npm install -g safe-chain-test
   2. Run: safe-chain setup
   3. Restart your terminal

🔍 Searching for NPM projects in examples/...
  📁 Found: examples/demo-project
  📁 Found: examples/nested-project/backend
  📁 Found: examples/nested-project/frontend
🎯 Found 3 NPM project(s)

📋 NPM Projects to be scanned:
  1. examples/demo-project
  2. examples/nested-project/backend
  3. examples/nested-project/frontend

Do you want to proceed with the security scan? [y/N]: y

🚀 Starting security scan for 3 project(s)...

📦 [1/3] Processing: examples/demo-project
  🗑️  Removing node_modules in examples/demo-project...
  ✅ node_modules removed from examples/demo-project
  📦 Running npm install in examples/demo-project...
  ✅ npm install completed in examples/demo-project
  🔍 Running Safe Chain scan in examples/demo-project...
  ⚠️  Safe Chain not found, running demo scan for examples/demo-project
  📊 Demo scan results for examples/demo-project:
  ✅ Demo scan completed - no vulnerabilities detected in examples/demo-project
  💡 Install Safe Chain for real vulnerability scanning
✅ [1/3] Completed: examples/demo-project

... (他のプロジェクトも同様に処理)

✅ All projects scanned successfully!
```

### 6. 開発・デバッグ用コマンド

```bash
# テスト実行
make test

# カバレッジ付きテスト
make test-coverage

# ベンチマークテスト
make bench

# リンター実行
make lint

# コードフォーマット
make format

# セキュリティチェック
make security

# 事前コミットチェック
make pre-commit
```

### 7. エラー対処

#### 「permission denied」エラー

```bash
# 実行権限を付与
chmod +x bin/npm-security-scanner
```

#### 「command not found」エラー

```bash
# パスを確認
ls -la bin/npm-security-scanner

# 絶対パスで実行
/full/path/to/npm-security-scanner/bin/npm-security-scanner
```

#### NPMプロジェクトが見つからない

- `package.json`ファイルが存在することを確認
- `node_modules`内のファイルは除外されます
- 検索対象ディレクトリのパスが正しいことを確認

### 8. 本番環境での使用

#### システムインストール

```bash
# システム全体にインストール
make install

# インストール確認
npm-security-scanner --version
```

#### Dockerでの使用

```bash
# Dockerイメージビルド
make docker-build

# Docker実行
docker run --rm -v $(pwd):/workspace npm-security-scanner /workspace
```

### 9. セキュリティ考慮事項

- **信頼できる環境での実行**: マルウェア検出ツールのため、信頼できる環境で実行してください
- **バックアップ**: 重要なプロジェクトは事前にバックアップを取ることを推奨
- **権限**: 必要最小限の権限で実行してください
- **ログ監視**: 実行ログを適切に監視・保存してください

### 10. トラブルシューティング

| 問題 | 原因 | 解決方法 |
|------|------|----------|
| スキャンが途中で止まる | npm installエラー | プロジェクトの`package.json`を確認 |
| 大量のプロジェクトが検出される | `node_modules`が除外されていない | 最新バージョンに更新 |
| Safe Chainエラー | インストールまたは設定問題 | インストール手順を再確認 |

### 11. サポート

- GitHub Issues: プロジェクトのIssueページ
- セキュリティ報告: `SECURITY.md`を参照
- 機能要望: Issue または Pull Request
