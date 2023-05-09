package logger

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type Logger struct {
	DebugEnabled bool
}

func (l *Logger) DisableDebug() {
	l.DebugEnabled = false
}

func (l *Logger) EnableDebug() {
	l.DebugEnabled = true
}

func NewLogger() *Logger {
	return &Logger{
		DebugEnabled: false,
	}
}

func (l *Logger) Info(message string) {
	fmt.Printf("%s %s %s\n", color.RedString(time.Now().Format("2006-01-02 15:04:05")), color.HiBlueString("[ INFO ]"), color.HiBlueString(message))
}

func (l *Logger) Warning(message string) {
	fmt.Printf("%s %s %s\n", color.RedString(time.Now().Format("2006-01-02 15:04:05")), color.HiRedString("[ WARNING ]"), color.HiRedString(message))
}

func (l *Logger) Error(message string) {
	fmt.Printf("%s %s %s\n", color.RedString(time.Now().Format("2006-01-02 15:04:05")), color.RedString("[ ERROR ]"), color.RedString(message))
}

func (l *Logger) Debug(message string) {
	if l.DebugEnabled {
		fmt.Printf("%s %s %s\n", color.RedString(time.Now().Format("2006-01-02 15:04:05")), color.GreenString("[ DEBUG ]"), color.HiGreenString(message))
	}
}
