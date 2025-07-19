package pool

import (
    "github.com/liuzhaomax/go-maxms/internal/core/pool/ws"
    "sync"
)

var once sync.Once
var pool *Pool

func init() {
    once.Do(func() {
        pool = &Pool{}
    })
}

func InitPool() *Pool {
    pool.WsPool = ws.InitWsPool()
    return pool
}

type Pool struct {
    WsPool *ws.WsPool
}

// Pool: {
//     WsPool: {
//         groups: [
//             {
//                 Name: roomName,
//                 Conns: map{
//                         userId: conn
//                     }
//                 },
//             {
//                 Name: roomName,
//                 Conns: map{
//                     userId: conn
//                 }
//             }
//         ]
//     },
//     TcpPool: [
//
//     ]
// }
