package global

var (
	NotifyMsg = make(chan WsMsg, 1000)
)

type MsgType int8

const (
	BasicSysInfo MsgType = 100 + iota
	TaskSend
)

type WsMsg struct {
	MsgType MsgType     `json:"msg_type"`
	Data    interface{} `json:"data"`
}

func Notify(msgType MsgType, data interface{}) {
	wsMsg := WsMsg{
		MsgType: msgType,
		Data:    data,
	}

	select {
	case NotifyMsg <- wsMsg:
		return
	default:
		<-NotifyMsg
		NotifyMsg <- wsMsg
	}
}
