package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type TaskRequest struct {
	FilePath string `json:"file_path"`
	OutPath  string `json:"out_path"`
}

const (
	// AmatsukazeAddTaskの固定設定
	amatsukazeIP      = "192.168.70.2"
	amatsukazePort    = "32768"
	amatsukazeService = "make"
)

type TaskResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Logs    string `json:"logs"`
}

var logFile *os.File

func main() {
	// ログファイルの初期化
	logDir := "/app/logs"
	os.MkdirAll(logDir, 0755)
	
	logFileName := filepath.Join(logDir, fmt.Sprintf("server_%s.log", time.Now().Format("20060102_150405")))
	var err error
	logFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("ログファイルを開けませんでした: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.Printf("サーバーを起動しました - ログファイル: %s\n", logFileName)

	http.HandleFunc("/execute", handleExecute)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	bindAddr := os.Getenv("BIND_ADDR")
	if bindAddr == "" {
		bindAddr = "0.0.0.0"
	}
	
	listenAddr := bindAddr + ":" + port
	log.Printf("サーバーが %s でリッスンしています\n", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("サーバーエラー: %v", err)
	}
}

func handleExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドのみサポートされています", http.StatusMethodNotAllowed)
		return
	}

	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("リクエストのデコードエラー: %v\n", err)
		respondError(w, "無効なリクエスト形式", err.Error())
		return
	}

	log.Printf("タスクリクエストを受信: File=%s, Out=%s\n",
		req.FilePath, req.OutPath)

	// AmatsukazeAddTaskコマンドを構築（IP、ポート、サービスは固定）
	ipWithPort := amatsukazeIP + ":" + amatsukazePort
	args := []string{
		"-ip", ipWithPort,
		"-s", amatsukazeService,
		"-f", req.FilePath,
		"-o", req.OutPath,
	}

	log.Printf("コマンド設定: IP=%s, Port=%s, Service=%s\n",
		amatsukazeIP, amatsukazePort, amatsukazeService)

	log.Printf("コマンドを実行: AmatsukazeAddTask %v\n", args)

	// コマンドを実行
	cmd := exec.Command("AmatsukazeAddTask", args...)
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	// コマンドの出力をログに記録
	if outputStr != "" {
		log.Printf("AmatsukazeAddTask 出力:\n%s\n", outputStr)
	}

	if err != nil {
		log.Printf("コマンド実行エラー: %v\n", err)
		respondError(w, "コマンド実行に失敗しました", fmt.Sprintf("%v\n%s", err, outputStr))
		return
	}

	log.Printf("タスクが正常に完了しました\n")

	// 成功レスポンス
	resp := TaskResponse{
		Success: true,
		Message: "タスクが正常に実行されました",
		Logs:    outputStr,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func respondError(w http.ResponseWriter, message string, details string) {
	resp := TaskResponse{
		Success: false,
		Message: message,
		Logs:    details,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(resp)
}
