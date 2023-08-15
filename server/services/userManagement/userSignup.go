package userManagement

import (
	"log"
	"net/http"

	chronAuthentication "github.com/TylerAldrich814/Chronicles/services/authentication"
	"github.com/TylerAldrich814/Chronicles/services/httpResponses"
)
const BUCKET_NAME string = "chronicles_users"
const COLLECTION string = "users"

// '/signup' PUT HTTP Endpoint. Handles signing a user up using
// both Firebase Auth and stores User metadata with FireStore
func handleUserSignup(w http.ResponseWriter, r *http.Request){
  fb := chronAuthentication.FirebaseAuth{}
  fb.Init().GetClient()
  resp := httpResponses.Response{}

  if err := r.ParseForm(); err != nil {
    http.Error(w, "Failed to Parse Form Data", http.StatusBadRequest)
  }

  email    := r.FormValue("email")
  username := r.FormValue("userName")
  passw    := r.FormValue("password")

  if len(email) == 0 || len(username) == 0 || len(passw) == 0 {
    http.Error(w, "Email, UserName or Password were missing", http.StatusNotAcceptable)
    return
  }

  if err := createUserMetadata(email, username); err != nil {
    http.Error(w, "Failed to add user metadata to Google Cloud", http.StatusInternalServerError)
  }

  _, err := fb.CreateUser(email, passw)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  response, err := resp.AddStatus("Successful").AddStatusCode(200).Build()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  _, writeErr := w.Write(response)
  if writeErr != nil {
    log.Printf(
      " --> ERROR: Failed to write JSON Response\n  -> Error: %v",
      err.Error(),
    )
  }
  log.Printf("Successfully Created User %v\n", email)
}

func createUserMetadata(email, userName string) error {
  fb := chronAuthentication.FirebaseAuth{}
  fb.Init().GetClient()

  fields := map[string]interface{}{
    "Email":       email,
    "OwnedRooms":  make(map[string]string),
    "JoinedRooms": make(map[string]string),
  }

  if err := fb.AddFirestoreDoc(COLLECTION, userName, fields); err != nil {
    log.Printf(" --> ERROR: Failed to Save user to Firestore\n  -> Error: %v\n", err.Error())
    return err
  }

  return nil
}

func CHECKOWNERSHIP(){
  // GET users/*UID/profile.json
  // GET rooms/name/metadata/ownership.json
  // if user.UID == room.owner.uid { return tru }
}
