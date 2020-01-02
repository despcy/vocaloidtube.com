package controller

import (
	"encoding/json"
	"log"

	"github.com/kataras/iris/websocket"
	viewrender "github.com/yangchenxi/VOCALOIDTube/view"
)

var WebsocketChan chan []byte

var NativeMessageHandler = func(nsConn *websocket.NSConn, msg websocket.Message) error {
	log.Printf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())

	//nsConn.Conn.Server().Broadcast(nsConn, msg)
	return nil
}

var ConnectHandler = func(c *websocket.Conn) error {
	log.Printf("[%s] Connected to server!", c.ID())
	//do the check and broadcast

	return nil
}

func SendWebSocData(Vcount string, VTime string, ip string, Loc string, VPath string, Scheck string) {
	packet := viewrender.StatusPageWSData{
		VisitCount: Vcount,
		TrafficData: viewrender.Traffic{
			VisitTime:     VTime,
			IpAddr:        ip,
			Location:      Loc,
			Path:          VPath,
			SecurityCheck: Scheck,
		},
	}
	JsonByte, err := json.Marshal(packet)

	if err != nil {

		log.Println(err)
	} else {

		WebsocketChan <- JsonByte
	}
}

var DisconnectHandler = func(c *websocket.Conn) {
	log.Printf("[%s] Disonnected to server!", c.ID())

}
