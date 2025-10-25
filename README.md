# Amatsukaze クライアント/サーバー

AmatsukazeAddTask を実行するためのシンプルなクライアント/サーバーシステムです。

## 構成

- **サーバー**: Linux (Docker) 上で実行され、`AmatsukazeAddTask` コマンドを実行します
- **クライアント**: Windows で実行され、サーバーにタスクリクエストを送信します
- **通信**: HTTP/JSON による RESTful API
- **ログ**: サーバーとクライアント両方でログファイルに出力

## 特徴

- ✅ Go 言語で実装されたシンプルで高速な設計
- ✅ 環境変数のプレースホルダーサポート（`{変数名}` 形式）
- ✅ サーバーからクライアントへのログストリーミング
- ✅ Docker 対応で簡単なデプロイ
- ✅ すべてのログをファイルに永続化

## 必要な環境

### サーバー側

- Docker & Docker Compose
- Linux 環境
- `AmatsukazeAddTask` コマンドがインストールされていること

### クライアント側

- Windows OS
- ビルド済みのクライアントバイナリ

## ダウンロード

ビルド済みのクライアントバイナリは[GitHub Releases](../../releases)からダウンロードできます。

- `amatsukaze-client-windows-amd64.exe` - Windows 用クライアント
- `amatsukaze-client-linux-amd64` - Linux 用クライアント

## セットアップ

### 1. ビルド（開発者向け）

```bash
# ビルドスクリプトに実行権限を付与
chmod +x build.sh

# ビルド実行
./build.sh
```

ビルドが完了すると以下のファイルが生成されます:

- `bin/amatsukaze-server` (Linux サーバー用)
- `bin/amatsukaze-client.exe` (Windows クライアント用)
- `bin/amatsukaze-client` (Linux クライアント用、テスト用)

### 2. サーバーの起動 (Docker)

```bash
# Docker Composeでサーバーを起動
docker-compose up -d

# ログを確認
docker-compose logs -f
```

サーバーは `http://localhost:8080` でリッスンします。

### 3. クライアントの使用 (Windows)

基本的な使用方法:

```cmd
amatsukaze-client.exe ^
  -url http://server-ip:8080 ^
  -s "make" ^
  -f "/app/share/TV-Record/example.ts" ^
  -o "/app/share/TV-Encoded"
```

#### プレースホルダーを使用する例:

```cmd
REM 環境変数を設定
set FileName=example.ts

REM プレースホルダーを使用
amatsukaze-client.exe ^
  -url http://server-ip:8080 ^
  -s "make" ^
  -f "/app/share/TV-Record/{FileName}" ^
  -o "/app/share/TV-Encoded"
```

#### コマンドライン引数:

| 引数      | デフォルト値 | 説明                       |
| --------- | ------------ | -------------------------- |
| `-url`    | (必須)       | サーバーの URL             |
| `-s`      | (必須)       | エンコード設定プロファイル |
| `-f`      | (必須)       | 入力ファイルパス           |
| `-o`      | (必須)       | 出力ディレクトリパス       |
| `-logdir` | `./logs`     | ログディレクトリパス       |
| `-nolog`  | `false`      | ログファイルを作成しない   |

## ログファイル

### サーバー側

- 場所: `/app/logs/server_YYYYMMDD_HHMMSS.log`
- Docker volume: `./logs/server/` にマウント
- `NOLOG` 環境変数を `true` に設定すると、ログファイルを作成せず標準出力のみに出力

### クライアント側

- 場所: `./logs/client_YYYYMMDD_HHMMSS.log`
- デフォルトでカレントディレクトリの `logs/` フォルダ
- `-nolog` フラグを使用すると、ログファイルを作成せず標準出力のみに出力

#### ログファイルを無効にする例:

**クライアント:**

```cmd
amatsukaze-client.exe -nolog -url http://server-ip:8080 -s "make" -f "/path/to/file" -o "/path/to/output"
```

**サーバー (docker-compose.yml):**

```yaml
environment:
  - NOLOG=true
```

## プレースホルダー機能

クライアントは `{環境変数名}` 形式のプレースホルダーをサポートしています。

例:

```cmd
set MY_FILE=video.ts
set OUTPUT_DIR=/app/share/output

amatsukaze-client.exe ^
  -f "/app/share/input/{MY_FILE}" ^
  -o "{OUTPUT_DIR}"
```

上記は以下のように展開されます:

```
-f "/app/share/input/video.ts"
-o "/app/share/output"
```

## API 仕様

### エンドポイント: POST /execute

**リクエスト:**

```json
{
  "file_path": "/app/share/TV-Record/example.ts",
  "out_path": "/app/share/TV-Encoded",
  "setting": "make"
}
```

**レスポンス (成功時):**

```json
{
  "success": true,
  "message": "タスクが正常に実行されました",
  "logs": "AmatsukazeAddTask の標準出力"
}
```

**レスポンス (失敗時):**

```json
{
  "success": false,
  "message": "エラーメッセージ",
  "logs": "エラー詳細"
}
```

## Docker Compose 設定

`docker-compose.example.yml` をコピーして `docker-compose.yml` として使用してください:

```bash
cp docker-compose.example.yml docker-compose.yml
```

### 環境変数

以下の環境変数でサーバーの動作を変更できます:

| 環境変数          | デフォルト値 | 説明                            |
| ----------------- | ------------ | ------------------------------- |
| `AMATSUKAZE_IP`   | `127.0.0.1`  | AmatsukazeServer の IP アドレス |
| `AMATSUKAZE_PORT` | `32768`      | AmatsukazeServer のポート番号   |
| `NOLOG`           | `false`      | `true` でログファイルを無効化   |

`docker-compose.yml` で設定例:

```yaml
environment:
  - AMATSUKAZE_IP=127.0.0.1
  - AMATSUKAZE_PORT=32768
```

### ボリュームマウント

- ログディレクトリ `./logs/server` をマウント
- 必要に応じて共有ディレクトリを追加してマウント

## トラブルシューティング

### サーバーに接続できない

1. サーバーが起動しているか確認:

   ```bash
   docker-compose ps
   ```

2. ファイアウォール設定を確認

3. サーバーの IP アドレスとポートが正しいか確認

### AmatsukazeAddTask が見つからない

Dockerfile で `AmatsukazeAddTask` のインストール手順を追加する必要があります。
`server/Dockerfile` を編集してください。

### ログが出力されない

1. ログディレクトリの書き込み権限を確認
2. Docker volume のマウントパスを確認

## ライセンス

このプロジェクトは MIT ライセンスの下で公開されています。

## 開発

### サーバーの開発モード実行

```bash
cd server
go run main.go
```

### クライアントの開発モード実行

```bash
cd client
go run main.go -url http://localhost:8080 -s "make" -f "/path/to/file" -o "/path/to/output"
```
