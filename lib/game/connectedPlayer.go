package game

import (
	"github.com/gorilla/websocket"
	"log"
)

type ConnectedPlayer struct {
	*Player
	ws   *websocket.Conn
	room *Room
}

func (cp *ConnectedPlayer) SendMessage(msg string) {
	go func() {
		log.Println("send message: ", msg)
		err := cp.ws.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			cp.room.Leave <- cp
			cp.ws.Close()
		}
	}()
}

// Receive msg from ws in goroutine
func (cp *ConnectedPlayer) receive() {
	defer func() {
		log.Println("start Player leave: ", cp, cp.room.Leave)
		cp.room.Leave <- cp
		cp.ws.Close()
	}()

	for {
		log.Println("run receive")
		_, action, err := cp.ws.ReadMessage()
		if err != nil {
			log.Println("Read message: ", err)
			break
		}
		parsedAction := string(action)
		cp.SetAction(parsedAction)
		cp.room.AddAction <- cp
	}
}

func NewConnectedPlayer(ws *websocket.Conn, player *Player, room *Room) *ConnectedPlayer {
	cp := &ConnectedPlayer{player, ws, room}
	log.Println("NewConnectedPlayer created")
	go cp.receive()
	log.Println("NewConnectedPlayer receive start")

	cp.room.Join <- cp

	return cp
}
