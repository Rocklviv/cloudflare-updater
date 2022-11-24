package logging

import (
	"log"
	"os"
)

type logger struct {
	InfoLogger     *log.Logger
	DebugLogger    *log.Logger
	ErrorLogger    *log.Logger
	CriticalLogger *log.Logger
}

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string)
	Critical(msg string)
}

func NewLogger() Logger {
	infoHandler := log.New(os.Stdout, "INFO:", log.Ldate|log.Ltime|log.Lshortfile)
	debugHandler := log.New(os.Stdout, "DEBUG:", log.Ldate|log.Ltime|log.Lshortfile)
	errorHandler := log.New(os.Stderr, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile)
	criticalHandler := log.New(os.Stderr, "CRITICAL:", log.Ldate|log.Ltime|log.Lshortfile)
	return &logger{
		InfoLogger:     infoHandler,
		DebugLogger:    debugHandler,
		ErrorLogger:    errorHandler,
		CriticalLogger: criticalHandler,
	}
}

func (l *logger) Info(msg string) {
	l.InfoLogger.Println(msg)
}

func (l *logger) Debug(msg string) {
	l.DebugLogger.Println(msg)
}

func (l *logger) Error(msg string) {
	l.ErrorLogger.Println(msg)
}

func (l *logger) Critical(msg string) {
	l.CriticalLogger.Fatalln(msg)
}
