package utils

import (
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func InitLogger() {
	logFile := &lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    50, // megabytes
		MaxBackups: 1,
		MaxAge:     1, // days
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	InfoLogger = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogToFile(message string) {
	InfoLogger.Println(message)
}

func LogError(message string) {
	ErrorLogger.Println(message)
}

func LogToConsole(message string) {
	log.Println(message)
}
