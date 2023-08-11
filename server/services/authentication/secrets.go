package chronAuthentication

import (
	"context"
	"log"
	"cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

const SECRETNAME string = "projects/ta-chitchat-backend/secrets/Chronicles_Firebase_Authentication/versions/1"

func accessSecretVersion(secretName string)( string, error ){
  ctx := context.Background()
  client, err := secretmanager.NewClient(ctx)
  if err != nil {
    log.Printf(
      " --> Failed to create new SecretManager Client\n  -> Error: %v",
      err.Error(),
    )
    return "", err
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
    return "", err
  }

  return string(result.Payload.Data), err
}

func FirebaseAuthCredentials()( []byte, error ){
  credentials, err := accessSecretVersion(SECRETNAME)
  if err != nil {
    return nil, err
  }

  return []byte(credentials), nil
}
