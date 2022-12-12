package wsproxy

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Socket    *websocket.Conn
	SessionID uint64
}

type Gateway struct {
	/*TODO*/
}

func (gw *Gateway) Handler(w http.ResponseWriter, r *http.Request) {
	/*TODO*/
}

func ConnectHandler(w http.ResponseWriter, r *http.Request) {

	// upgrader := websocket.Upgrader{
	// 	/*TODO*/
	// }

	// wsConn, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	/*TODO*/
	// 	return
	// }

	// conn := Connection{
	// 	Socket:    wsConn,
	// 	SessionID: 0, /*will be filled after authentication*/
	// }
}
