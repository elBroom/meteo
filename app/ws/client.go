package ws

import (
	"time"

	"github.com/elBroom/meteo/app/schema"
	"github.com/fasthttp-contrib/websocket"
	"github.com/mailru/easyjson"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 1 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan *schema.Indication
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, conn *websocket.Conn) {
	client := &Client{hub: hub, conn: conn, send: make(chan *schema.Indication, 256)}
	client.hub.register <- client

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			b, _ := easyjson.Marshal(message)
			client.conn.WriteMessage(websocket.TextMessage, b)
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
