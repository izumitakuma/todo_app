package utils

import (
	"io"
	"log"
	"os"
)

// LoggingSettings 設定されたログファイルと標準出力にログを出力する設定を行う
func LoggingSettings(logfile string) {
	// ログファイルを開く
	logFile, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	// 標準出力とログファイルに同時にログを出力する設定
	multiLogFile := io.MultiWriter(os.Stdout, logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(multiLogFile)
}
