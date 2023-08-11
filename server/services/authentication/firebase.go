package chronAuthentication

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
  "firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type FirebaseAuth struct {
  ctx    context.Context
  app    *firebase.App
  client *auth.Client
}

const PROJECTID string = "ta-chitchat-backend"

// Returns firbase.app, is Fatal on credentials and client creation error.
func( fb *FirebaseAuth )Init() *FirebaseAuth{
  credentials, err := FirebaseAuthCredentials()
  if err != nil {
    log.Fatalf(
      " --> FATAL: Failed to create Firebase Client Option\n  -> Error: %v",
      err.Error(),
    )
  }
  ctx := context.Background()
  fb.ctx = ctx
  opt := option.WithCredentialsJSON([]byte(credentials))
  config := &firebase.Config{ProjectID: PROJECTID}
  app, err := firebase.NewApp(fb.ctx, config, opt)
  if err != nil {
    log.Fatalf(
      " --> FATAL: Failed to create Firebase Context\n  -> Error: %v",
      err.Error(),
    )
  }
  fb.app = app
  return fb
}

func( fb *FirebaseAuth )GetClient() *FirebaseAuth{
  client, err := fb.app.Auth(fb.ctx)
  if err != nil {
    log.Fatal(
      " --> FATAL: Failed to create Firebase App Authentication\n  -> Error: %v",
      err.Error(),
    )
  }
  fb.client = client
  return fb
}
