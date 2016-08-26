package dbapi

import "log"

// Logger is an interface for logging progress and other messages
type Logger interface {
	Write(string)
}

// StderrLogger is a logger for printing messages to standard error
type StderrLogger struct {
}

func (l StderrLogger) Write(s string) {
	log.Printf(s)
}
