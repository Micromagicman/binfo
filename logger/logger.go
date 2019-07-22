package logger

import "fmt"

func Info(message string) {
	log("INFO", message)
}

func Warning(message string) {
	log("WARNING", message)
}

func Error(message string) {
	log("ERROR", message)
}

func log(level string, message string) {
	fmt.Println("["+level+"] "+message)
}
