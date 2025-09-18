# NPM Security Scanner

NPMパッケージのマルウェア感染対策のためのセキュリティスキャナーツールです。指定されたディレクトリ配下のすべてのNPMプロジェクトを再帰的に検索し、Safe Chainを使用して一括でセキュリティスキャンを実行します。

## 機能

- Safe Chainのグローバルインストール確認と対話的インストール
- ディレクトリ再帰検索によるpackage.json検出
- プロジェクトリストの表示と確認
- node_modules削除とnpm install再実行
- Safe Chainによる一括セキュリティスキャン

## 技術選択

セキュリティツールとしての独立性と安全性を重視し、**Go言語**で実装します。

- 単一バイナリとして配布可能
- 標準ライブラリが充実
- メモリ安全性とセキュリティ面での信頼性
- クロスプラットフォーム対応

## 開発要件

- Clean Code、SOLID原則、OWASP Secure Coding準拠
- 適切なエラーハンドリングとログ出力
- テストコード完備
- セキュリティ最優先の設計

## インストールと使用方法

### ビルド方法

```bash
# 依存関係のインストール
go mod download

# ビルド
make build
# または
go build -o bin/npm-security-scanner .
```

### 使用方法

```bash
# カレントディレクトリをスキャン
./bin/npm-security-scanner

# 特定のディレクトリをスキャン
./bin/npm-security-scanner /path/to/projects

# ヘルプを表示
./bin/npm-security-scanner --help
```

### 実行例

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
```

### Makeコマンド

```bash
# 利用可能なコマンドを表示
make help

# テスト実行
make test

# リンター実行
make lint

# フォーマット
make format

# セキュリティチェック
make security

# すべてのプラットフォーム向けビルド
make build-all
```

## Safe Chain について

このツールは **Safe Chain** と連携してセキュリティスキャンを実行します。

### インストール方法

```bash
# NPMでグローバルインストール
npm install -g safe-chain-test

# セットアップ実行
safe-chain setup

# ターミナル再起動後、確認
safe-chain --version
```

### デモモード

Safe Chainがインストールされていない場合、ツールはデモモードで動作し、実際のスキャンの代わりにモックスキャンを実行します。

## セキュリティ機能

- ✅ **再帰的検索**: 指定ディレクトリ配下のすべてのNPMプロジェクトを検出
- ✅ **node_modules除外**: 依存関係内のpackage.jsonは除外
- ✅ **対話的確認**: スキャン実行前にプロジェクト一覧を表示
- ✅ **クリーンスキャン**: node_modules削除→npm install→スキャン
- ✅ **エラーハンドリング**: 適切なエラーメッセージとログ出力
- ✅ **デモモード**: Safe Chain未インストール時の安全な動作
