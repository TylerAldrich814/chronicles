package chronAuthentication

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"

	// "firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type AuthToken *auth.Token

const PROJECTID string = "ta-chitchat-backend"

type FirebaseAuth struct {
  ctx    context.Context
  app    *firebase.App
  client *auth.Client
}


// Returns firbase.app, is Fatal on credentials and client creation error.
func( fb *FirebaseAuth )Init() *FirebaseAuth{
  credentials, err := FirebaseAuthCredentials()
  if err != nil {
    log.Fatalf(
      " --> FATAL: Failed to create Firebase Client Option\n  -> Error: %v\n",
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
    log.Fatalf(
      " --> FATAL: Failed to create Firebase App Authentication\n  -> Error: %v",
      err.Error(),
    )
  }
  fb.client = client
  return fb
}

func( fb *FirebaseAuth )GetAuthToken(
  ctx context.Context,
  r   *http.Request,
)(  AuthToken,error  ){
  authHeader := r.Header.Get("Authroization")
  idToken := strings.TrimPrefix(authHeader, "Bearer ")
  return fb.client.VerifyIDToken(ctx, idToken)
}

// Validates user on Firebase Auth. Returns user UUID
func( fb *FirebaseAuth)ValidateUserEmail(email string)( string,error ){
  record, err := fb.client.GetUserByEmail(fb.ctx, email)
  if err != nil {
    log.Printf(" --> ERROR: User doesn't exist\n  -> Error: %v", err.Error())
  }
  return record.UID, err
}

func( fb *FirebaseAuth )CreateUser(
  email,
  passwd string,
)(*auth.UserRecord, error ){
  user := auth.UserToCreate{}
  user.Email(email).Password(passwd)

  return fb.client.CreateUser(fb.ctx, &user)
}

func( fb *FirebaseAuth )GetUser(
  email string,
)( *auth.UserRecord,error ){
  return fb.client.GetUserByEmail(fb.ctx, email)
}

func( fb *FirebaseAuth )DeleteUser(email string) error {
  return fb.client.DeleteUser(fb.ctx, email)
}

// Adds a Document to Firestore.
// collection -> Firestore Collection the document will be added to
// title -> Document Title
// documetn -> map[string]interface{}, Documents data
func( fb *FirebaseAuth )AddFirestoreDoc(
  collection string,
  document string,
  fields map[string]interface{},
) error {
  client, err := fb.app.Firestore(fb.ctx)
  if err != nil {
    return err
  }
  defer client.Close()

  _, err = client.Collection(collection).Doc(document).Set(fb.ctx, fields)
  if err != nil {
    return err
  }

  return nil
}

func( fb FirebaseAuth )GetFirestreDoc(
  collection string,
  document   string,
)( map[string]interface{}, error ){
  client, err := fb.app.Firestore(fb.ctx)
  if err != nil {
    return nil, err
  }
  defer client.Close()

  result, err := client.Collection(collection).Doc(document).Get(fb.ctx)
  if err != nil {
    return nil, err
  }
  return result.Data(), nil
}

func( fb *FirebaseAuth )DeleteFirestoreDoc(
  collection string,
  document   string,
) error {
  client, err := fb.app.Firestore(fb.ctx)
  if err != nil {
    return err
  }
  defer client.Close()

  _, err = client.Collection(collection).Doc(document).Delete(fb.ctx)
  if err != nil {
    return err
  }
  return nil
}
func( fb *FirebaseAuth )DocExists(collection, doc string)( bool, error){
  client, err := fb.app.Firestore(fb.ctx)
  if err != nil {
    return false, err
  }

  obj := client.Collection(collection).Doc(doc)

  if doc, err := obj.Get(context.Background()); err != nil {
    return false, nil
  }else if err != nil && doc.Exists(){
    return true, nil
  } else {
    return false, nil
  }
}

func( fb *FirebaseAuth )UpdateFirestoreField(
  collection string,
  document   string,
  fieldUpdate map[string]interface{},
) error {
  client, err := fb.app.Firestore(fb.ctx)
  if err != nil {
    return err
  }
  defer client.Close()

  _, err = client.Collection(collection).Doc(document).Set(fb.ctx, fieldUpdate)
  if err != nil {
    return err
  }
  return nil
}

// Takes in a reference to an Interface, updates interface with fields if found.
// if not, returns an error.
func( fb *FirebaseAuth )GetFirestoreField(
  collection string,
  document   string,
  field      *interface{},
) error {
  client, err := fb.app.Firestore(fb.ctx)
  if err != nil {
    return err
  }
  defer client.Close()

  doc, err := client.Collection(collection).Doc(document).Get(fb.ctx)
  if err != nil {
    return err
  }

  err = doc.DataTo(field)

  return err
}
func( fb *FirebaseAuth)GetFirestoreSubField(
  collection string,
  document   string,
  subfield   string,
)( []interface{}, error ){
  client, err := fb.app.Firestore(fb.ctx)
  if err != nil {
    return nil, err
  }
  defer client.Close()

  ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
  defer cancel()

  doc, err := client.
    Collection(collection).
    Doc(document).
    Get(ctx)
  if err != nil {
    return nil, err
  }
  item, ok := doc.Data()[subfield].([]interface{})
  if !ok {
    return nil, fmt.Errorf("Failed to read subfield '%s'", subfield)
  }
   return item, nil
}

func( fb *FirebaseAuth )DeleteFirstoreField(
  collection string,
  document   string,
  field      string,
) error {
  client, err := fb.app.Firestore(fb.ctx)
  if err != nil {
    return err
  }
  defer client.Close()

  _, err = client.
    Collection(collection).
    Doc(document).
    Update(fb.ctx, []firestore.Update{
      {
        Path: field,
        Value: firestore.Delete,
      },
    })
  if err != nil {
    return err
  }
  return nil
}
