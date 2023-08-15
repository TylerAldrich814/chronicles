package userManagement

import (
	"encoding/json"
	"log"
	"net/http"

	chronAuthentication "github.com/TylerAldrich814/Chronicles/services/authentication"
)

type User struct {
  Email string `json:"email"`
  Uid   string `json:"uid"`
}

func handleUserGet(w http.ResponseWriter, r *http.Request){
  fb := chronAuthentication.FirebaseAuth{}
  fb.Init().GetClient()

  if err := r.ParseForm(); err != nil {
    log.Printf(
      " --> ERROR: Failed to get User Email\n  -> Error: %v",
      err.Error(),
    )
    http.Error(w, "Failed to Get User Email", http.StatusInternalServerError)
    return
  }
  email := r.FormValue("email")

  user, err := fb.GetUser(email)
  if err != nil {
    log.Printf(
      " --> ERROR: Failed to Get User\n  -> Error: %v",
      err.Error(),
    )
    http.Error(w, "Failed to Get User", http.StatusInternalServerError)
    return
  }

  foundUser := User{Email: user.Email, Uid: user.UID }

  userjson, err := json.Marshal(foundUser)
  if err != nil {
    log.Printf(
      " --> ERROR: Failed to Marshal Found User\n  -> Error: %v",
      err.Error(),
    )
    http.Error(w, "Failed to Marshal found user", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  _, writeErr := w.Write(userjson)
  if writeErr != nil {
    log.Printf(
      " --> ERROR: Failed to write JSON Response\n  -> Error: %v",
      err.Error(),
    )
  }
}
