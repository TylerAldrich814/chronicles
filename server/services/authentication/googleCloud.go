package chronAuthentication

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GoogleCloudAuth struct {
  ctx     context.Context
  client  *storage.Client
  buckets map[string]string
}

func( gc *GoogleCloudAuth)Init() *GoogleCloudAuth {
  ctx := context.Background()
  credentials, err := GCSAuthCredentials()
  if err != nil {
    log.Fatalf(
      " --> FATAL: Failed to create Google Cloud Client\n  -> Error: %v\n",
      err.Error(),
    )
  }
  client, err := storage.NewClient(ctx, option.WithCredentialsJSON(credentials))
  if err != nil {
    log.Fatalf(
      " --> FATAL: Failed to create new Google Cloud Client\n  -> Error: %v\n",
      err.Error(),
    )
  }
  gc.ctx = ctx
  gc.client = client
  gc.buckets = make(map[string]string)
  return gc
}

func( gc *GoogleCloudAuth)Close() {
  gc.client.Close()
}

func( gc *GoogleCloudAuth)WriteFileToBucket(
  bucketName string,
  parent     string,
  filename   string,
  data       interface{},
) error {
  prefix := fmt.Sprintf("%s/.keep", parent)
  var jsonData  = []byte("")

  if data != nil {
    prefix = fmt.Sprintf("%s/%s", parent, filename)

    var err error
    if jsonData, err = json.Marshal(data); err != nil {
      log.Printf(" --> ERROR: AddFileToBucket - Data Marshal failed\n  -> Error: %v", err.Error())
      return err
    }
  }

  buf := bytes.NewBuffer(jsonData)
  ctx, cancel := context.WithTimeout(gc.ctx, time.Second*50)
  defer cancel()

  wc := gc.client.
    Bucket(bucketName).
    Object(prefix).
    NewWriter(ctx)

  if _, err := io.Copy(wc, buf); err != nil {
    log.Fatalf(
      " --> FATAL: WriteFileToBucket: io.Copy Failed \n  -> Error: %v\n",
      err.Error(),
    )
  }

  if err := wc.Close(); err != nil {
    log.Printf(" --> WriteFileToBucket Write.Close Failed\n  -> Error: %v", err.Error())
    return err
  }

  return nil
}
