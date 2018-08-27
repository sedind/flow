package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/fatih/color"
)

const lFormat = "=== %s ==="

// Logger object
type Logger struct {
	log *log.Logger
}

// New creates new logger
func New(name string, enableColors bool) *Logger {
	color.NoColor = !enableColors
	if runtime.GOOS == "windows" {
		color.NoColor = true
	}
	return &Logger{
		log: log.New(os.Stdout, fmt.Sprintf("%s: ", name), log.LstdFlags),
	}
}

// Success logs success message
func (l *Logger) Success(msg interface{}, args ...interface{}) {
	l.log.Print(color.GreenString(fmt.Sprintf(lFormat, msg), args...))
}

// Error logs error message
func (l *Logger) Error(msg interface{}, args ...interface{}) {
	l.log.Print(color.RedString(fmt.Sprintf(lFormat, msg), args...))
}

// Print logs message
func (l *Logger) Print(msg interface{}, args ...interface{}) {
	l.log.Printf(fmt.Sprintf(lFormat, msg), args...)
}
