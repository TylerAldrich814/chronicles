package main

import (
	// "fmt"
	// "io/ioutil"
	"log"
	"net/http"

	// "github.com/TylerAldrich814/Chronicles/roomManagement"
	"github.com/TylerAldrich814/Chronicles/clientManagement"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/TylerAldrich814/Chronicles/services/userManagement"
)
var upgrader = websocket.Upgrader{
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request){
  conn, err := upgrader.Upgrade(w,r, nil)
  if err != nil {
    log.Printf(
      " --> WSHandler failed to upgrade HTTP Connection\n  -> Error: %v",
      err.Error(),
    )
  }
  room := r.URL.Query().Get("room")
  var client = &clientManagement.Client{}
  client.
    AddConnection(conn).
    AddSender(make(chan []byte)).
    AddRoom(room)

  go client.Write()
  go client.Read()
}

func signupHandler(w http.ResponseWriter, r *http.Request){

}

func main(){
  r := mux.NewRouter()
  userManagement.UserManagementRoutes(r)
}
