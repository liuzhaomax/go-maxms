package core

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type WebSocket struct {
	ReadBufferSize    int      `mapstructure:"read_buffer_size"`
	WriteBufferSize   int      `mapstructure:"write_buffer_size"`
	HandshakeTimeout  int      `mapstructure:"handshake_timeout"`
	Subprotocols      []string `mapstructure:"subprotocols"`
	EnableCompression bool     `mapstructure:"enable_compression"`
}

func InitWebSocket() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:    cfg.Lib.WebSocket.ReadBufferSize * 1024,
		WriteBufferSize:   cfg.Lib.WebSocket.WriteBufferSize * 1024,
		HandshakeTimeout:  time.Duration(cfg.Lib.WebSocket.HandshakeTimeout) * time.Second,
		Subprotocols:      cfg.Lib.WebSocket.Subprotocols,
		EnableCompression: cfg.Lib.WebSocket.EnableCompression,
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get("Origin") == cfg.App.Domain
		},
		Error: func(w http.ResponseWriter, r *http.Request, status int, err error) {
			LogFailure(ProtocolUpgradeFailed, "WebSocket upgrade failed", err)
			http.Error(w, "WebSocket upgrade failed", status)
		},
	}
}
