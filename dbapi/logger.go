package dbapi

import (
	"log"

	"golang.org/x/net/websocket"
)

// Logger is an interface for logging progress and other messages
type Logger interface {
	Write(string)
	LogInterval() int
}

// StderrLogger is a logger for printing messages to standard error. Implements the dbapi.Logger interface.
type StderrLogger struct {
}

func (l StderrLogger) Write(s string) {
	log.Println(s)
}
func (l StderrLogger) LogInterval() int {
	return 10000
}

// WebSockLogger is a logger for printing messages to a web socket. Implements the dbapi.Logger interface.
type WebSockLogger struct {
	websock *websocket.Conn
}

func NewWebSockLogger(websock *websocket.Conn) WebSockLogger {
	return WebSockLogger{websock: websock}
}

func (l WebSockLogger) Write(msg string) {
	//log.Println(msg)
	websocket.Message.Send(l.websock, msg)
}

func (l WebSockLogger) LogInterval() int {
	return 100
}

// SilentLogger is a muted logger, used for testing to skip too much confusing test output
type SilentLogger struct {
}

func (l SilentLogger) Write(s string) {

}
func (l SilentLogger) LogInterval() int {
	return 10000
}
