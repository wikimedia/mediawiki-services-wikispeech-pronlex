package dbapi

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/websocket"
)

// Logger is an interface for logging progress and other messages
type Logger interface {
	Progress(string)
	Write(string)
	LogInterval() int
}

// StderrLogger is a logger for printing messages to standard error. Implements the dbapi.Logger interface.
type StderrLogger struct {
	LogIntervalVar int
}

// Progress logs progress info
func (l StderrLogger) Progress(s string) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("\r%s\r", s))
}

// Write logs a message string
func (l StderrLogger) Write(s string) {
	log.Println(s)
}

// LogInterval speficies logging interval (to be used by the calling process)
func (l StderrLogger) LogInterval() int {
	return l.LogIntervalVar
}

// StdoutLogger is a logger for printing messages to standard out. Implements the dbapi.Logger interface.
type StdoutLogger struct {
	LogIntervalVar int
}

// Progress logs progress info
func (l StdoutLogger) Progress(s string) {
	//fmt.Fprintf(os.Stdout, fmt.Sprintf("\r%s\r", s))
	//fmt.Print(".")
	log.Printf(fmt.Sprintf("%s\n", s))
}

// Write logs a message string
func (l StdoutLogger) Write(s string) {
	fmt.Println(s)
}

// LogInterval speficies logging interval (to be used by the calling process)
func (l StdoutLogger) LogInterval() int {
	return l.LogIntervalVar
}

// WebSockLogger is a logger for printing messages to a web socket. Implements the dbapi.Logger interface.
type WebSockLogger struct {
	websock        *websocket.Conn
	LogIntervalVar int
}

// NewWebSockLogger creates a new websock logger using the input connection
func NewWebSockLogger(websock *websocket.Conn) WebSockLogger {
	return WebSockLogger{websock: websock}
}

// Progress logs progress info
func (l WebSockLogger) Progress(msg string) {
	websocket.Message.Send(l.websock, msg)
}

// Write logs a message string
func (l WebSockLogger) Write(msg string) {
	websocket.Message.Send(l.websock, msg)
}

// LogInterval speficies logging interval (to be used by the calling process)
func (l WebSockLogger) LogInterval() int {
	return l.LogIntervalVar
}

// SilentLogger is a muted logger, used for testing to skip too much confusing test output
type SilentLogger struct {
}

// Write logs a message string
func (l SilentLogger) Write(s string) {

}

// Progress logs progress info
func (l SilentLogger) Progress(s string) {

}

// LogInterval speficies logging interval (to be used by the calling process)
func (l SilentLogger) LogInterval() int {
	return -1
}
