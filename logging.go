package main

import (
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stderr)
}

func LogLn(v any) {
	log.Println(v)
}

func Logf(format string, v ...any) {
	log.Printf(format, v...)
}

func DebugLogln(v any) {
	if (settings.Logging.Debug) {
		log.Println(v)
	}
}

func DebugLogf(format string, v ...any) {
	if (settings.Logging.Debug) {
		log.Printf(format, v...)
	}
}