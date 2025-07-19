package ws

import (
    "context"
    "fmt"
    "github.com/gorilla/websocket"
    "github.com/liuzhaomax/go-maxms/internal/core"
    "sync"
    "sync/atomic"
    "time"
)

var once sync.Once
var wsPool *WsPool

func init() {
    once.Do(func() {
        wsPool = &WsPool{
            mux:    new(sync.RWMutex),
            groups: []*ConnGroup{},
        }
    })
}

func InitWsPool() *WsPool {
    return wsPool
}

type IWsPool interface {
    GetGroupByName(string) *ConnGroup
    Push(*ConnGroup)
    Remove(string) bool
    Range(func(*ConnGroup) bool)
    Filter(func(*ConnGroup) bool) []*ConnGroup
    HandleInactiveConnGroups(time.Duration)
}

type WsPool struct {
    mux    *sync.RWMutex
    groups []*ConnGroup // [{ Name: roomName, Conns: [{userId: conn}] }]
}

func (wp *WsPool) GetGroupByName(name string) *ConnGroup {
    wp.mux.RLock()
    defer wp.mux.RUnlock()

    for _, group := range wp.groups {
        if group.Name == name {
            return group
        }
    }
    return nil
}

func (wp *WsPool) Push(connGroup *ConnGroup) {
    wp.mux.Lock()
    defer wp.mux.Unlock()
    wp.groups = append(wp.groups, connGroup)
}

func (wp *WsPool) Remove(name string) bool {
    wp.mux.Lock()
    defer wp.mux.Unlock()

    for i, group := range wp.groups {
        if group.Name == name {
            ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
            defer cancel()

            // 异步关闭所有连接
            go func(ctx context.Context, g *ConnGroup) {
                var (
                    wg      *sync.WaitGroup
                    success int32
                    failed  int32
                )

                g.Conns.m.Range(func(key, value interface{}) bool {
                    select {
                    case <-ctx.Done():
                        return false // 超时停止处理新连接
                    default:
                        wg.Add(1)
                        go func(k interface{}, conn *websocket.Conn) {
                            defer wg.Done()

                            // 优雅关闭流程
                            if err := CloseConn(ctx, conn); err != nil {
                                core.LogFailure(core.CloseException, fmt.Sprintf("关闭连接 %v 失败", k), err)
                                atomic.AddInt32(&(failed), 1)
                                return
                            }

                            // 从连接池移除
                            g.Conns.m.Delete(k)
                            atomic.AddInt32(&success, 1)
                        }(key, value.(*websocket.Conn))
                        return true
                    }
                })

                wg.Wait()
                core.LogSuccess(fmt.Sprintf("房间 %s 关闭完成: 成功%d个, 失败%d个",
                    g.Name, success, failed))
            }(ctx, group)

            // 从组列表中移除
            wp.groups[i] = wp.groups[len(wp.groups)-1]
            wp.groups = wp.groups[:len(wp.groups)-1]
            return true
        }
    }
    return false
}

// CloseConn 优雅关闭单个连接
func CloseConn(ctx context.Context, conn *websocket.Conn) error {
    // 1. 发送关闭帧
    err := conn.WriteControl(
        websocket.CloseMessage,
        websocket.FormatCloseMessage(websocket.CloseNormalClosure, "conn_group_close"),
        time.Now().Add(3*time.Second),
    )
    if err != nil {
        return fmt.Errorf("发送关闭帧失败: %w", err)
    }

    // 2. 等待关闭确认或超时
    select {
    case <-time.After(1 * time.Second): // 给客户端响应时间
    case <-ctx.Done():
        return ctx.Err()
    }

    // 3. 最终关闭
    err = conn.Close()
    if err != nil {
        return fmt.Errorf("最终关闭失败: %w", err)
    }

    return nil
}

func (wp *WsPool) Range(fn func(group *ConnGroup) bool) {
    wp.mux.RLock()
    defer wp.mux.RUnlock()

    for _, group := range wp.groups {
        if !fn(group) {
            break
        }
    }
}

func (wp *WsPool) Filter(filterFn func(group *ConnGroup) bool) []*ConnGroup {
    wp.mux.RLock()
    defer wp.mux.RUnlock()

    var result []*ConnGroup
    for _, group := range wp.groups {
        if filterFn(group) {
            result = append(result, group)
        }
    }
    return result
}

// HandleInactiveConnGroups 启动检查协程(5小时不活动视为闲置) -> go HandleInactiveConnGroups(5*time.Hour)
func (wp *WsPool) HandleInactiveConnGroups(timeout time.Duration) {
    for {
        // 每小时检查一次
        time.Sleep(1 * time.Hour)

        // 找出所有超时的房间
        inactiveGroups := wp.Filter(func(group *ConnGroup) bool {
            return time.Since(group.updatedAt) > timeout
        })

        // 移除这些房间
        for _, group := range inactiveGroups {
            wp.Remove(group.Name)
            core.LogSuccess(fmt.Sprintf(
                "移除闲置房间: %s (最后活动时间: %s, 已闲置: %s)",
                group.Name,
                group.updatedAt.Format("2006-01-02T15:04:05"),
                time.Since(group.updatedAt).Truncate(time.Second), // 去除纳秒精度
            ))
        }
    }
}
