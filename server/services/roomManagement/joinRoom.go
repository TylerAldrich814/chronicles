package roomManagement

import (
	"net/http"
	"strings"

	chronAuthentication "github.com/TylerAldrich814/Chronicles/services/authentication"
)

func joinRoom(
  w http.ResponseWriter,
  r *http.Request,
){
  fb := chronAuthentication.FirebaseAuth{}
  fb.Init().GetClient()

  token, err := fb.GetAuthToken(r.Context(), r)
  if err != nil {
    http.Error(w, err.Error(), http.StatusUnauthorized)
    return
  }
  userName := r.FormValue("userName")
  roomName := r.FormValue("roomName")

  if exist, err := fb.DocExists(COLLECTION, roomName); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  } else if !exist {
    http.Error(w, err.Error(), http.StatusNoContent)
    return
  }
  users, err := fb.GetFirestoreSubField(COLLECTION, roomName, "Users")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  var Status string
  for _, user := range users {
      if userId, ok := user.(string); ok && userId == token.UID {
        Status = "User already exists in room."

      }
    }
  }
}
