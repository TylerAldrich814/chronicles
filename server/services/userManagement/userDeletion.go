package userManagement

import (
	"encoding/json"
	"log"
	"net/http"

	chronAuthentication "github.com/TylerAldrich814/Chronicles/services/authentication"
)

// /remove-user/
//   --> First uses proivided email address from client to get
//       user data from firebase. With the procided user UID
//       we then call upon Firebase to delete this user.
func handleUserDeletion(w http.ResponseWriter, r *http.Request){
  fb := chronAuthentication.FirebaseAuth{}
  fb.Init().GetClient()

  userData := User{}
  if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
    log.Printf(
      " --> ERROR: Failed to Decode Request Body\n  -> Error: %v",
      err.Error(),
    )
    http.Error(w, "Failed to Decode Request Body", http.StatusInternalServerError)
    return
  }
  email := userData.Email

  user, err := fb.GetUser(email)
  if err != nil {
    log.Printf(
      " --> ERROR: Failed to Get User\n  -> Error: %v",
      err.Error(),
    )
    http.Error(w, "Failed to Get User", http.StatusInternalServerError)
    return
  }
  uid := user.UID

  if err := fb.DeleteUser(uid); err != nil {
    log.Printf(
      "Failed to delete user \"%v\"\n Error: %v",
      email, err.Error(),
    )
    http.Error(w, "Failed to delete user", http.StatusInternalServerError)
  }
  log.Printf("Successfully Delteed user \"%v\"", email)
  w.WriteHeader(http.StatusOK)
}
