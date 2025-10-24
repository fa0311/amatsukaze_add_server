FROM golang:1.21-alpine AS builder

WORKDIR /build

# Go モジュールのコピーとダウンロード
COPY go.mod go.sum* ./
RUN go mod download || true

# ソースコードのコピー
COPY main.go .

# ビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# 実行用の軽量イメージ
FROM ubuntu:24.04

# 必要なパッケージをインストール
RUN apt-get update && apt-get install -y \
    ca-certificates \
    wget \
    curl \
    xz-utils \
    && rm -rf /var/lib/apt/lists/*

# AmatsukazeAddTaskをインストール
RUN UBUNTU_VERSION=24.04 \
    && mkdir -p /tmp/Amatsukaze \
    && curl -s https://api.github.com/repos/rigaya/Amatsukaze/releases/latest \
    | grep "browser_download_url.*tar.xz" | grep "Ubuntu${UBUNTU_VERSION}" | cut -d : -f 2,3 | tr -d \" \
    | wget -i - -O - | tar -xJ -C /tmp/Amatsukaze \
    && install /tmp/Amatsukaze/exe_files/AmatsukazeAddTask /usr/local/bin/ \
    && rm -rf /tmp/Amatsukaze

WORKDIR /app

# ビルドしたバイナリをコピー
COPY --from=builder /build/server .

# ログディレクトリを作成
RUN mkdir -p /app/logs

# ポート公開
EXPOSE 8080

# サーバーの実行
CMD ["./server"]
