# Changelog

NPM Security Scannerの変更履歴

## [1.3.0] - 2024-09-19

### 🎉 メジャーアップデート

#### **新機能**
- **📊 包括的レポートシステム**: 実行ログを構造化収集し、ターミナル/HTML/JSON形式で出力
- **🌐 クリック可能リンク**: file://URL形式でHTMLレポートをワンクリック表示
- **🤖 GitHub Actions CI/CD**: macOS・Windows向け自動バイナリリリース
- **🔍 改善されたSafe Chain検出**: マルチシェル環境対応とより正確なセットアップ状態確認
- **🛡️ 詳細な脆弱性解析**: npm audit出力の詳細解析と脆弱性情報の構造化

#### **機能追加**
- HTML形式の美しいレポート生成（レスポンシブ・グラデーション対応）
- JSON形式の機械読み取り可能レポート
- プロジェクト別の詳細スキャン結果記録
- 実行時間・成功率・エラー詳細の完全追跡
- 脆弱性の重要度・パッケージ・修正状況の記録
- 自動リリース用のチェックサム生成

#### **修正**
- Safe Chain検出の多重シェル確認（zsh/bash/sh対応）
- npm audit exit status 1の正常処理（脆弱性発見＝成功）
- レポート生成の自動実行（ユーザビリティ向上）
- 脆弱性パース精度の大幅改善
- switch文によるパフォーマンス向上

#### **技術改善**
- Clean Code原則に基づく完全リファクタリング
- 定数化によるマジックナンバー撲滅
- 構造体メモリレイアウト最適化
- セキュア権限設定（0o600）
- 関数分割による複雑度低減
- 全linter警告の解消

## [1.2.0] - 2024-09-19

### コードリファクタリング
- **BREAKING**: コード構造の抜本的見直し
  - `install.go`を`safechain.go`に統合・廃止
  - 機能別ファイル分割: `scanner.go`、`safechain.go`、`utils.go`
  - 未使用関数の完全除去: `runSafeChainTest`、`isCommandNotFoundError`

### 修正
- 未使用importの削除とコンパイル警告の解決
- 関数の重複排除と責任分離の改善
- メモリ効率とコードの可読性向上

### 技術的改善
- コード品質: 637行 → 効率的な機能分割
- ファイル構成: 明確な責任分離 (CLI/Scanner/SafeChain/Utils)
- テスト完全動作確認済み

## [1.1.0] - 2024-09-19

### 変更
- **BREAKING**: Aikido Safe ChainからSafe Chainに変更
  - パッケージ名を `@aikido/safe-chain` から `safe-chain-test` に変更
  - コマンドを `aikido scan` から `npm audit` (Safe Chainでラップ) に変更

### 追加
- Safe Chainセットアップ確認機能
- より詳細なインストールガイダンス
- 改善されたデモモード
- 実際の脆弱性検出機能（npm audit使用）
- ターミナル再起動案内

### 修正
- 未使用関数 `ensureDirectoryExists` を削除
- node_modulesディレクトリの確実な除外
- エラーハンドリングの改善

### 技術的変更
- Safe Chainのsetup/teardownコマンドに対応
- npm auditとnpm audit fixを組み合わせたスキャン
- より適切なセットアップ状態の検出

## [1.0.0] - 2024-09-19

### 初回リリース
- NPMプロジェクトの再帰的検索機能
- 一括セキュリティスキャン機能
- node_modules削除とクリーンインストール
- 対話的プロジェクト確認
- デモモード対応
- 包括的なテストスイート
- Docker対応
- Makefile基盤の開発環境
