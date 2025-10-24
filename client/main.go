package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type TaskRequest struct {
	FilePath string `json:"file_path"`
	OutPath  string `json:"out_path"`
	Setting  string `json:"setting"`
}

type TaskResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Logs    string `json:"logs"`
}

func main() {
	// コマンドライン引数の定義
	serverURL := flag.String("url", "", "サーバーのURL (例: http://localhost:8080)")
	filePath := flag.String("f", "", "入力ファイルパス")
	outPath := flag.String("o", "", "出力ディレクトリパス")
	setting := flag.String("s", "", "エンコード設定プロファイル")
	logDir := flag.String("logdir", "./logs", "ログディレクトリ")

	flag.Parse()

	// 必須パラメータのチェック
	if *serverURL == "" {
		log.Fatal("エラー: -url パラメータは必須です")
	}
	if *filePath == "" {
		log.Fatal("エラー: -f パラメータは必須です")
	}
	if *outPath == "" {
		log.Fatal("エラー: -o パラメータは必須です")
	}
	if *setting == "" {
		log.Fatal("エラー: -s パラメータは必須です")
	}

	// ログファイルの初期化
	os.MkdirAll(*logDir, 0755)
	logFileName := filepath.Join(*logDir, fmt.Sprintf("client_%s.log", time.Now().Format("20060102_150405")))
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("ログファイルを開けませんでした: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.Printf("クライアントを開始しました - ログファイル: %s\n", logFileName)

	// プレースホルダーの置換
	expandedFilePath := expandPlaceholders(*filePath)
	expandedOutPath := expandPlaceholders(*outPath)
	expandedSetting := expandPlaceholders(*setting)

	log.Printf("リクエストパラメータ:\n")
	log.Printf("  サーバー: %s\n", *serverURL)
	log.Printf("  入力ファイル: %s\n", expandedFilePath)
	log.Printf("  出力ディレクトリ: %s\n", expandedOutPath)
	log.Printf("  設定プロファイル: %s\n", expandedSetting)

	// リクエストの作成
	req := TaskRequest{
		FilePath: expandedFilePath,
		OutPath:  expandedOutPath,
		Setting:  expandedSetting,
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("リクエストのJSONエンコードエラー: %v", err)
	}

	// サーバーにリクエストを送信
	log.Printf("サーバーにリクエストを送信中...\n")
	resp, err := http.Post(*serverURL+"/execute", "application/json", bytes.NewBuffer(reqJSON))
	if err != nil {
		log.Fatalf("サーバーへの接続エラー: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスの読み取り
	var taskResp TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		log.Fatalf("レスポンスのデコードエラー: %v", err)
	}

	// サーバーからのログを出力
	log.Printf("\n=== サーバーからのレスポンス ===\n")
	log.Printf("成功: %v\n", taskResp.Success)
	log.Printf("メッセージ: %s\n", taskResp.Message)
	if taskResp.Logs != "" {
		log.Printf("\n=== AmatsukazeAddTask ログ ===\n%s\n", taskResp.Logs)
	}

	if !taskResp.Success {
		log.Fatal("タスクの実行に失敗しました")
	}

	log.Printf("タスクが正常に完了しました\n")
}

// expandPlaceholders は {ENV_VAR} 形式のプレースホルダーを環境変数で置換します
func expandPlaceholders(input string) string {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		// {} を取り除いて環境変数名を取得
		envVar := match[1 : len(match)-1]
		value := os.Getenv(envVar)
		if value == "" {
			log.Printf("警告: 環境変数 '%s' が見つかりません、元の値を使用します\n", envVar)
			return match
		}
		return value
	})
}
