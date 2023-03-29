package controller

import (
	"sync"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/randolphcyg/demo_gin_websocket/global"
	"github.com/randolphcyg/demo_gin_websocket/middleware"
)

var once sync.Once

func WsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := middleware.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		// 把与客户端的链接添加到客户端链接池中
		middleware.WsConns.Store(ws.RemoteAddr().String(), ws)
		defer middleware.CloseWebSocketConn(ws.RemoteAddr().String(), ws)

		log.Debugf("用户 %s 连接成功", ws.RemoteAddr().String())

		once.Do(middleware.MsgHandler)
		var wsMsg global.WsMsg
		for {
			err := ws.ReadJSON(&wsMsg)
			if err != nil {
				return
			}
			log.Debug("websocket recv:", wsMsg)
		}
	}
}
