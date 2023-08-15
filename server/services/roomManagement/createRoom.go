package roomManagement

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	chronAuthentication "github.com/TylerAldrich814/Chronicles/services/authentication"
	utilites "github.com/TylerAldrich814/Chronicles/services/utilities"
	"github.com/google/uuid"
)

const BUCKET_NAME string = "chronicle-rooms"
const COLLECTION string = "chronicle-rooms"

type AuthToken chronAuthentication.AuthToken

// "/create-room" HTTP PUT Request.
//    Requires specific form data as input, Verifys the User.
//    Then Creats the Metadata for a Room, stored in FireStare
//    TODO: Create Databases for Room Audio and U->R & U->U Chat logs
func createRoom(
  w http.ResponseWriter,
  r *http.Request,
){
  fb := chronAuthentication.FirebaseAuth{}
  fb.Init().GetClient()

  token, err := fb.GetAuthToken(r.Context(), r)
  if err != nil {
    http.Error(w, err.Error(), http.StatusForbidden)
    return
  }

  if err := r.ParseForm(); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  userName := r.FormValue("userName")
  roomName := r.FormValue("roomName")
  peers := make(map[string]AuthToken) // TODO: Create Peer Struct for nested data.
  roomId, err := uuid.NewUUID()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  // Safety check => Verifies if the provided Room Name was already taken or not.
  if exist, err := fb.DocExists(COLLECTION, roomName); err != nil {
    http.Error(w, "Error while checking if Room exists yet", http.StatusInternalServerError)
    return
  } else if exist {
    http.Error(w, err.Error(), http.StatusForbidden)
    return
  }

  room := map[string]interface{}{
    "roomId":     roomId.String(),
    "ownername":  userName,
    "ownerToken": token,
    "peers":      peers,
  }

  if err := fb.AddFirestoreDoc(COLLECTION, roomName, room); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  // Good Job 👍( TODO. Create DataStrucures directory, pull 'Respinse' into a legit Struct)
  response := struct{
    Status     string            `json:"status"`
    StatusCode int32             `json:"statusCode"`
    Body       map[string]string `json:"body"`
  }{
    Status: "Success",
    StatusCode: 200,
    Body: map[string]string{
      "RoomName": roomName,
      "RoomID": roomId.String(),
    },
  }

  jsonResp,err := json.Marshal(response)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  w.Write(jsonResp)
}

// I was going to use Google Cloud Storage at first, but decided to go with
// Firestore instead. I'll keep this code for now just in case..
func createRoomDirectory(
  ctx       context.Context,
  client    *storage.Client,
  roomId    string,
) error {

  bucket := client.Bucket(BUCKET_NAME)
  prefix := fmt.Sprintf("rooms/%s/", roomId)
  object := bucket.Object(prefix)
  wc := object.NewWriter(ctx)
  emptyBucketObject := []byte("")

  if _, err := wc.Write(emptyBucketObject); err != nil {
    return err
  }
  if err := wc.Close(); err != nil {
    return err
  }
  return nil
}

func createRoomMetadata(
  ctx      context.Context,
  client   *storage.Client,
  roomName string,
  roomId   string,
  creator  string,
) error {
  bucket := client.Bucket(BUCKET_NAME)
  prefix := fmt.Sprintf("rooms/%s/metadata", roomId)
  object := bucket.Object(prefix)
  wc := object.NewWriter(ctx)
  defer wc.Close()

  if _, err := wc.Write([]byte("")); err != nil {
    log.Printf(" --> Failed to Create metadata folder\n  -> Error: %v", err.Error())
    return err
  }

  metaDataObj := bucket.Object(fmt.Sprintf("%s/ownership.json", prefix))
  metadataWriter := metaDataObj.NewWriter(ctx)
  defer metadataWriter.Close()

  ownershipMetadata := struct {
    OwnerID  string `json:"ownerID"`
    RoomName string `json:"roomName"`
  }{
    OwnerID: creator,
    RoomName: roomName,
  }
  metadataEncoder := json.NewEncoder(metadataWriter)
  if err := metadataEncoder.Encode(ownershipMetadata); err != nil {
    log.Printf(" --> Failed to Encode metadata\n  -> Error: %v", err.Error())
    return err
  }
  return nil
}
