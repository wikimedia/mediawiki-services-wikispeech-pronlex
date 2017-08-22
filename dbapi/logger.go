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

func (l StderrLogger) Progress(s string) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("\r%s\r", s))
}
func (l StderrLogger) Write(s string) {
	log.Println(s)
}
func (l StderrLogger) LogInterval() int {
	return l.LogIntervalVar
}

// StdoutLogger is a logger for printing messages to standard out. Implements the dbapi.Logger interface.
type StdoutLogger struct {
	LogIntervalVar int
}

func (l StdoutLogger) Progress(s string) {
	fmt.Fprintf(os.Stdout, fmt.Sprintf("\r%s\r", s))
}
func (l StdoutLogger) Write(s string) {
	fmt.Println(s)
}
func (l StdoutLogger) LogInterval() int {
	return l.LogIntervalVar
}

// WebSockLogger is a logger for printing messages to a web socket. Implements the dbapi.Logger interface.
type WebSockLogger struct {
	websock        *websocket.Conn
	LogIntervalVar int
}

func NewWebSockLogger(websock *websocket.Conn) WebSockLogger {
	return WebSockLogger{websock: websock}
}

func (l WebSockLogger) Write(msg string) {
	websocket.Message.Send(l.websock, msg)
}

func (l WebSockLogger) Progress(msg string) {
	websocket.Message.Send(l.websock, msg)
}

func (l WebSockLogger) LogInterval() int {
	return l.LogIntervalVar
}

// SilentLogger is a muted logger, used for testing to skip too much confusing test output
type SilentLogger struct {
}

func (l SilentLogger) Write(s string) {

}
func (l SilentLogger) Progress(s string) {

}
func (l SilentLogger) LogInterval() int {
	return -1
}
