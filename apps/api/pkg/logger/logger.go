package logger

import (
	"log"
	"os"
)

// Placeholder for structured logging setup
// Future: Integrate with proper logging library (e.g., zap, zerolog)

func Init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
