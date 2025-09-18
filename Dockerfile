# Multi-stage build for npm-security-scanner
FROM golang:1.21-alpine AS builder

# セキュリティのためのユーザー作成
RUN adduser -D -s /bin/sh -u 1001 scanner

# 作業ディレクトリ設定
WORKDIR /app

# Go modulesファイルをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド（静的リンク）
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static" -s -w' -o npm-security-scanner .

# 最終イメージ（軽量）
FROM alpine:3.18

# セキュリティアップデート
RUN apk --no-cache add ca-certificates tzdata && \
    apk --no-cache update && \
    apk --no-cache upgrade

# 非rootユーザーを作成
RUN adduser -D -s /bin/sh -u 1001 scanner

# Node.js and npmをインストール（Aikido Safe Chainの実行に必要）
RUN apk add --no-cache nodejs npm

# 作業ディレクトリ設定
WORKDIR /app

# バイナリをコピー
COPY --from=builder /app/npm-security-scanner /usr/local/bin/npm-security-scanner

# 実行権限を設定
RUN chmod +x /usr/local/bin/npm-security-scanner

# ユーザーを切り替え
USER scanner

# ヘルスチェック
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD npm-security-scanner --version || exit 1

# エントリーポイント
ENTRYPOINT ["npm-security-scanner"]

# デフォルトコマンド
CMD ["--help"]
