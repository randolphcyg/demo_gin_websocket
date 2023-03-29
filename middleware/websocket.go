package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/randolphcyg/demo_gin_websocket/global"
)

var WsConns sync.Map
var WsUpgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	// 在程序结束时关闭所有连接
	// 注意：只有在程序崩溃或被杀死时才会执行 defer 函数
	defer func() {
		WsConns.Range(func(key, value interface{}) bool {
			conn, success := value.(*websocket.Conn)
			if !success {
				log.Error("invalid websocket connection")
				return true
			}
			conn.Close()
			return true
		})
	}()
}

func MsgHandler() {
	log.Info("MsgHandler start")
	// 创建一个定时器用于服务端心跳
	pingTicker := time.NewTicker(time.Second * 60)

	for {
		select {
		// 从消息通道接收消息，然后推送给前端
		case msg := <-global.NotifyMsg:
			data, err := json.Marshal(msg)
			if err != nil {
				log.Error(err)
				continue
			}
			WsConns.Range(func(key, value interface{}) bool {
				conn, success := value.(*websocket.Conn)
				if !success {
					log.Error("invalid websocket connection")
					return true
				}

				err = conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					log.Error(err)
				}
				return true
			})
		case <-pingTicker.C:
			WsConns.Range(func(key, value interface{}) bool {
				conn, success := value.(*websocket.Conn)
				if !success {
					log.Error("invalid websocket connection")
					return true
				}
				// 服务端心跳:每60秒ping一次客户端，查看其是否在线
				conn.SetWriteDeadline(time.Now().Add(time.Second * 20))
				err := conn.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					log.Println("send ping err:", err)
					conn.Close()
					WsConns.Delete(key)
					return true
				}
				return true
			})

		default:
			time.Sleep(time.Millisecond * 100) // 防止忙等待
		}
	}
}

func CloseWebSocketConn(key interface{}, value interface{}) bool {
	conn, success := value.(*websocket.Conn)
	if !success {
		log.Error("invalid websocket connection")
		return false
	}
	log.Debugf("用户 %s 断开连接", conn.RemoteAddr().String())
	WsConns.Delete(key)
	conn.Close()
	return true
}
