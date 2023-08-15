package userManagement

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// Sinple test that tests all three functionalites
//  - Firebase User Creation
//  - Firebase User Retrieval
//  - Firebase User Deletion
func TestUserSignupAndDelete(t *testing.T){
  createAndDelete := false
  onlyDelete := false


  email := "SuperFuperNewUser1233@email.com"
  userName := "SuperDuperMyUsername666"
  passw := "SuperSecurePassword123"
  values := url.Values{}

  values.Add("email", email)
  values.Add("password", passw)
  values.Add("userName", userName)

  data := strings.NewReader(values.Encode())
  req, err := http.NewRequest("POST", "/signup", data)
  if err != nil {
    t.Fatal(err)
  }
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

  rr := httptest.NewRecorder()

  if !onlyDelete {
    log.Printf("Creating User %v...\n", email)
    handleUserSignup(rr, req)
    if status := rr.Code; status != http.StatusCreated {
      t.Errorf(
        "Handler returned wrong status code: got %v want %v",
        status, http.StatusOK,
      )
      return
    }
    time.Sleep(4 * time.Second)
  }

  if createAndDelete || onlyDelete {
    log.Printf("Deleting User %v...\n", email)

    userData := struct{
      Email string `json:"email"`
    }{
      Email: email,
    }
    jsonData, err := json.Marshal(userData)
    if err != nil {
      t.Fatalf(" --> Failed to Marshal UserData: %v\n", err.Error())
      return
    }

    req, err = http.NewRequest("DELETE", "/remove-user", bytes.NewReader(jsonData))
    if err != nil {
      t.Fatalf(" --> Failed to create new HTTP Request: %v\n", err.Error())
      return
    }
    req.Header.Add("Content-Type", "application/json")
    rr = httptest.NewRecorder()

    handleUserDeletion(rr, req)
    if status := rr.Code; status != http.StatusOK {
      t.Errorf(
        "Handler returned wrong status code: got %v want %v",
        status, http.StatusOK,
      )
      return
    }
  }
}
