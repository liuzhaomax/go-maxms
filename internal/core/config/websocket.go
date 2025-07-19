package config

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
)

type webSocketConfig struct {
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
			return r.Header.Get("Origin") == fmt.Sprintf("%s://%s", cfg.Server.Ws.Protocol, cfg.App.Domain) ||
				r.Header.Get("Origin") == fmt.Sprintf("%s://%s:%s", cfg.Server.Ws.Protocol, cfg.Server.Ws.Host, cfg.Server.Ws.Port)
		},
		Error: func(w http.ResponseWriter, r *http.Request, status int, err error) {
			LogFailure(ext.ProtocolUpgradeFailed, "WebSocket upgrade failed", err)
			http.Error(w, "WebSocket upgrade failed", status)
		},
	}
}
