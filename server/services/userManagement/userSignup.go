package userManagement

import ( "encoding/json"
	"log"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	chronAuthentication "github.com/TylerAldrich814/Chronicles/services/authentication"
)

func handleUserSignup(w http.ResponseWriter, r *http.Request){
  fb := chronAuthentication.FirebaseAuth{}
  fb.Init().GetClient()

  if err := r.ParseForm(); err != nil {
    http.Error(w, "Failed to Parse Form Data", http.StatusBadRequest)
  }

  email := r.FormValue("email")
  passw := r.FormValue("password")

  uid := func() *uuid.UUID{
    reties := 5
    for retry := 1; retry <= reties; retry++ {
      uid, err := uuid.NewUUID()
      if err != nil {
        log.Printf(
          " --> Warn: Failed to create UUID, %v trys left\n  -> Error: %v",
          reties - retry, err.Error(),
        )
      }
      return &uid
    }
    log.Fatal(" --> FATAL: Failed to create UUID")
    return nil
  }()
  if uid == nil {
    http.Error(w, "Failed to create User UUID, Please try again", http.StatusInternalServerError)
    return
  }

  user := auth.UserToCreate{}
  user.Email(email).Password(passw).UID(uid.String())

  record, err := fb.CreateUser(&user)
  if err != nil {
    log.Printf(
      " --> ERROR: Failed to Create user\n  -> Error: %v",
      err.Error(),
    )
  }

  recordJson, err := json.Marshal(record)
  if err != nil {
    log.Printf(
      " --> Error: Failed to Marshal User Record\n  -> Error: %v",
      err.Error(),
    )
    http.Error(w, "Failed to Marshal User Record", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  _, writeErr := w.Write(recordJson)
  if writeErr != nil {
    log.Printf(
      " --> ERROR: Failed to write JSON Response\n  -> Error: %v",
      err.Error(),
    )
  }
  log.Printf("Successfully Created User %v\n", email)
}
