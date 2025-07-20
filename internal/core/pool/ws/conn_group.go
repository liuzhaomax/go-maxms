package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type IConnGroup interface {
	Store(string, *websocket.Conn)
	Load(string) (*websocket.Conn, bool)
	Delete(string)
	Range(func(string, *websocket.Conn) bool)
	Broadcast(any) error
}

type ConnGroup struct {
	mux       *sync.RWMutex
	createdAt time.Time
	updatedAt time.Time     // 广播时更新，如果n个小时没有结算，则直接判定为闲置房间，直接移除，见HandleInactiveConnGroups
	ID        uint          // roomID
	Conns     *connsSyncMap // userId: conn
}

type connsSyncMap struct {
	m *sync.Map
}

func NewConnGroup(ID uint) *ConnGroup {
	return &ConnGroup{
		mux:       new(sync.RWMutex),
		createdAt: time.Now(),
		updatedAt: time.Now(),
		ID:        ID,
		Conns: &connsSyncMap{
			m: new(sync.Map),
		},
	}
}

func (cg *ConnGroup) Store(key string, value *websocket.Conn) {
	cg.Conns.m.Store(key, value)
}

func (cg *ConnGroup) Load(key string) (value *websocket.Conn, ok bool) {
	v, ok := cg.Conns.m.Load(key)
	return v.(*websocket.Conn), ok
}

func (cg *ConnGroup) Delete(key string) {
	cg.Conns.m.Delete(key)
}

func (cg *ConnGroup) Range(f func(key string, value *websocket.Conn) bool) {
	cg.Conns.m.Range(func(k, v any) bool {
		return f(k.(string), v.(*websocket.Conn))
	})
}

func (cg *ConnGroup) Broadcast(info any) error {
	cg.mux.Lock()
	cg.updatedAt = time.Now() // 更新活动时间
	cg.mux.Unlock()

	var wg sync.WaitGroup
	var broadcastErr error
	var errMux sync.Mutex

	cg.Range(func(key string, conn *websocket.Conn) bool {
		wg.Add(1)
		go func(conn *websocket.Conn, key string) {
			defer wg.Done()

			var msg []byte
			switch v := info.(type) {
			case []byte:
				msg = v
			case string:
				msg = []byte(v)
			default:
				// 可以添加更复杂的序列化逻辑
				msg = []byte(fmt.Sprintf("%v", v))
			}

			// 发送消息
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				cg.Delete(key) // 发送失败则移除连接
				err = conn.Close()
				if err != nil {
					broadcastErr = err
				}

				errMux.Lock()
				if broadcastErr == nil {
					broadcastErr = err
				}
				errMux.Unlock()
			}
		}(conn, key)
		return true
	})

	wg.Wait()
	return broadcastErr
}
