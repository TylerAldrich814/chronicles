package httpResponses

import chronAuthentication "github.com/TylerAldrich814/Chronicles/services/authentication"

type AuthToken = chronAuthentication.AuthToken

type Room struct {
  Uid        string    `json:"uid"`
  RoomName   string    `json:"roomName"`
  Owner      string    `json:"owner"`
  OwnerToken AuthToken `json:"authToken"`
  users      map[string]AuthToken `json:"users"`
}

