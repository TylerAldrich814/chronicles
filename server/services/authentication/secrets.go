package chronAuthentication

import (
	"context"
	"log"
	"cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)
const _FIREBASE_SECRET string = "projects/ta-chitchat-backend/secrets/Chronicles_Firebase_Authentication/versions/latest"
const _GCS_SECRET      string = "projects/ta-chitchat-backend/secrets/chronicle-gcs-credentials/versions/latest"

func accessSecretVersion(secretName string)( []byte, error ){
  ctx := context.Background()
  client, err := secretmanager.NewClient(ctx)
  if err != nil {
    log.Printf(
      " --> Failed to create new SecretManager Client\n  -> Error: %v",
      err.Error(),
    )
    return nil, err
  }
  defer client.Close()

  req := &secretmanagerpb.AccessSecretVersionRequest{
    Name: secretName,
  }
  result, err := client.AccessSecretVersion(ctx, req)
  if err != nil {
    log.Printf(
      " --> Failed to Access Secret\n  -> Error: %v",
      err.Error(),
    )
    return nil, err
  }

  return result.Payload.Data, err
}

func GCSAuthCredentials()( []byte, error ){
  credentials, err := accessSecretVersion(_GCS_SECRET)
  if err != nil {
    log.Printf(" --> Failed to obtain GCS Secret Credentials\n --> Error %v\n", err.Error())
    return nil, err
  }
  return credentials, nil
}

func FirebaseAuthCredentials()( []byte, error ){
  credentials, err := accessSecretVersion(_FIREBASE_SECRET)
  if err != nil {
    return nil, err
  }

  return credentials, nil
}
