package roomManagement

import (
	"log"
	"testing"
)

func TestCreateRoom(t *testing.T){
  creator  := "Tyler"
  roomName := "TestRoom"
  roomId   := "a3bf123e9c1"

  if err := createRoom(roomName, roomId, creator); err != nil {
    log.Printf(" --> Failed to create Room %v", roomName)
  }
}
