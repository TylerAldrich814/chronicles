package chroniclesRecording

import (
	"log"
	"net/http"

	// "cloud.google.com/go/pubsub"
	"github.com/TylerAldrich814/Chronicles/roomManagement"
	"github.com/gorilla/websocket"
)
type RoomManager = roomManagement.RoomManager
type Connection  = roomManagement.Connection

var upgrader = websocket.Upgrader{
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,
}


func WSHandler(w http.ResponseWriter, r *http.Request){
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Printf(
      " --> Failed to upgrade HTTP connection to Websocket\n  -> Error: %v",
      err.Error(),
    )
  }
  defer conn.Close()

  roomName := r.URL.Query().Get("room")
  if roomName == "" {
    log.Printf(" --> Failed, Room Name was not provided!")
    return
  }

  manager.JoinRoom(conn, roomName)

  for {
    messageType, p, err := conn.ReadMessage()
    if err != nil {
      log.Printf(
        " --> Failed to read Connection Message\n  -> Error: %v",
        err.Error(),
      )
    }
    if err := conn.WriteMessage(messageType, p); err != nil {
      log.Printf(
        " --> Failed to Write to Connection\n  -> Error: %v",
        err.Error(),
      )
    }

    // Hnad the Audio data here.
  }
}

func handleAudioChunk(conn *websocket.Conn, roomName string, audioChunk []byte){

}


var manager = RoomManager{Rooms: make(map[string][]*Connection)}

func WSmain(){
  http.HandleFunc("/", WSHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}
