package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/randolphcyg/demo_gin_websocket/controller"
	"github.com/randolphcyg/demo_gin_websocket/global"
)

func main() {
	r := gin.Default()
	// 前端调用此接口建立websocket长连接
	r.GET("/websocket", controller.WsHandler())

	// send data to front end
	r.POST("/send", send)

	err := r.Run("localhost:8090")
	if err != nil {
		panic(err)
	}
}

type DemoData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Total int `json:"total"`
		List  []struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
			Addr string `json:"addr"`
		} `json:"list"`
	} `json:"data"`
}

func send(*gin.Context) {
	data := DemoData{
		Code: 1,
		Msg:  "success",
		Data: struct {
			Total int `json:"total"`
			List  []struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
				Addr string `json:"addr"`
			} `json:"list"`
		}{
			Total: 2,
			List: []struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
				Addr string `json:"addr"`
			}{
				{Name: "xiao ming",
					Age:  13,
					Addr: "shanghia",
				},
				{Name: "xiao hong",
					Age:  12,
					Addr: "shanghia",
				},
			},
		},
	}
	log.Info("send data by websocket")
	global.Notify(global.TaskSend, data)
}
