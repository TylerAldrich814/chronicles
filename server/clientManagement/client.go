package clientManagement

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
  conn *websocket.Conn
  room string
  send chan[]byte
}

func( c *Client)AddConnection(conn *websocket.Conn) *Client {
  c.conn = conn
  return c
}
func( c *Client)AddRoom(room string) *Client {
  c.room = room
  return c
}
func( c *Client )AddSender(send chan[]byte) *Client {
  c.send = send
  return c
}

func( c *Client)Read(){
  defer c.conn.Close()
  for {
    messageType, _, err := c.conn.ReadMessage()
    if err != nil {
      log.Printf(
        " --> Client Failed to Read next message\n  -> Error: %v",
        err.Error(),
      )
      return
    }
    if messageType == websocket.BinaryMessage {
      // Audio message
    } else if messageType == websocket.TextMessage {
      // room Message
    }
    // TODO: Implement our logic for handling and routing our messages
  }
}

func( c *Client)Write() {
  defer c.conn.Close()
  for message := range c.send{
    err := c.conn.WriteMessage(websocket.TextMessage, message)
    if err != nil {
      log.Printf(
        " --> Client failed to Write their next message\n  -> Error: %v",
        err.Error(),
      )
      return
    }
  }
}
