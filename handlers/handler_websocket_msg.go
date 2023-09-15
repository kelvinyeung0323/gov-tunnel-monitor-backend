package handlers

import (
	"encoding/json"
	"log"

	"github.com/ky/gov-tunnel-monitor-backend/base/websocket"
	"github.com/ky/gov-tunnel-monitor-backend/handlers/monitor"
)

func HandlerWebsocketRecievedMessage(d *websocket.Client, msg *websocket.Message) error {
	if msg.MsgType == websocket.WS_BROCAST || msg.MsgType == websocket.WS_P2P {
		//这两种情况暂不处理，只处理CMD类型
		return nil
	}
	b, err := json.Marshal(msg.Body)
	if err != nil {
		log.Printf("error:命令操作错误 %v", err)
		return err
	}
	oper := &monitor.Operation{}
	err = json.Unmarshal(b, oper)
	if err != nil {
		log.Printf("error:命令操作错误 %v", err)
		return err
	}
	monitor.Hub.OperationHandler(oper)
	//处理完信息后不作反馈
	return nil
}
