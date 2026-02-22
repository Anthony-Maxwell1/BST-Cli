package ws

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

const wsURL = "ws://localhost:5000"

func SendPacket(packet map[string]any) {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		fmt.Println("WS error:", err)
		return
	}
	defer conn.Close()

	if err := conn.WriteJSON(packet); err != nil {
		fmt.Println("Write error:", err)
		return
	}

	var response any
	if err := conn.ReadJSON(&response); err != nil {
		fmt.Println("Read error:", err)
		return
	}

	out, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(out))
}
