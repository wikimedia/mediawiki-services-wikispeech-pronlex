package dbapi

import (
	"golang.org/x/net/websocket"
	"log"
)

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

// WebSockLogger is a logger for printing messages to a web socket
type WebSockLogger struct {
	websock *websocket.Conn
}

func NewWebSockLogger(websock *websocket.Conn) WebSockLogger {
	return WebSockLogger{websock: websock}
}

func (wsl WebSockLogger) Write(msg string) {
	websocket.Message.Send(wsl.websock, msg)
}
