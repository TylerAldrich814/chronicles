package roomManagement

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Conneciton represents a Connected Client
type Connection struct {
   Conn *websocket.Conn
   Room string
}

// RoomManager Handles the Clients and Connections
type RoomManager struct {
  sync.RWMutex
  Rooms map[string][]*Connection
}
func( rm *RoomManager )JoinRoom(conn *websocket.Conn, roomName string){
  rm.Lock()
  defer rm.Unlock()

  room, exists := rm.Rooms[roomName]
  if !exists {
    room = []*Connection{}
  }

  room = append(room, &Connection{Conn: conn, Room: roomName})
  rm.Rooms[roomName] = room
}

func(rm *RoomManager) Broadcast(
  conn *websocket.Conn,
  roomName string,
  message []byte,
){
  rm.RLock()
  defer rm.Unlock()

  room, exists := rm.Rooms[roomName]
  if !exists {
    log.Printf(" --> Warning: You are not connected to the Room \"%v\"", roomName)
    return
  }

  for _, connection := range room {
    if connection.Conn != conn {
      if err := connection.Conn.WriteMessage(
        websocket.BinaryMessage,
        message,
      ); err != nil {
        log.Printf(
          "  --> Failed to Communitcate with all other connections in the Room"+
          " \"%v\"\n  -> Error: %v",
          connection.Room, err.Error(),
        )
      }
    }
  }
}

func( rm *RoomManager )LeaveRoom(conn *websocket.Conn, roomName string){
  rm.Lock()
  defer rm.Unlock()

  room, exists := rm.Rooms[roomName]
  if !exists {
    log.Printf(" --> Warning: You are not connected to the Room \"%v\"", roomName)
    return
  }

  for i, connection := range room {
    if connection.Conn == conn {
      room = append(room[:i], room[i+1:]...)
      break
    }
  }
  rm.Rooms[roomName] = room
}

