package httpResponses


type User struct {
  Email       string
  UserName    string
  JoinedRooms map[string]string
  OwnedRooms  map[string]string
}
