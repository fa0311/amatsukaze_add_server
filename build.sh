#!/bin/bash

# Amatsukaze クライアント/サーバー ビルドスクリプト

set -e

echo "=== Amatsukaze ビルドスクリプト ==="

# サーバーのビルド（Linux用）
echo ""
echo "📦 サーバーをビルド中..."
cd server
go build -o ../bin/amatsukaze-server main.go
echo "✅ サーバービルド完了: bin/amatsukaze-server"
cd ..

# クライアントのビルド（Windows用）
echo ""
echo "📦 クライアント（Windows用）をビルド中..."
cd client
GOOS=windows GOARCH=amd64 go build -o ../bin/amatsukaze-client.exe main.go
echo "✅ クライアントビルド完了: bin/amatsukaze-client.exe"
cd ..

# クライアントのビルド（Linux用、テスト用）
echo ""
echo "📦 クライアント（Linux用）をビルド中..."
cd client
go build -o ../bin/amatsukaze-client main.go
echo "✅ クライアントビルド完了: bin/amatsukaze-client"
cd ..

echo ""
echo "🎉 すべてのビルドが完了しました！"
echo ""
echo "出力ファイル:"
echo "  - bin/amatsukaze-server (Linux サーバー)"
echo "  - bin/amatsukaze-client.exe (Windows クライアント)"
echo "  - bin/amatsukaze-client (Linux クライアント)"
