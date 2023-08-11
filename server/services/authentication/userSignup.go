package chronAuthentication

import (
	"encoding/json"
	"log"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/google/uuid"
)

func handleUserSignup(w http.ResponseWriter, r *http.Request){
  fb := FirebaseAuth{}
  fb.Init().GetClient()

  email := r.URL.Query().Get("email")
  passw := r.URL.Query().Get("password")
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
  user.
    Email(email).
    Password(passw).
    UID(uid.String())

  record, err := fb.client.CreateUser(fb.ctx, &user)
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
}
