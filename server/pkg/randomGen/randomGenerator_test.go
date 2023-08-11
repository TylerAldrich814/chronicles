package randomGen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func TestRandomNumberWithinBounds(t *testing.T){
  url := "https://us-central1-ta-chitchat-backend.cloudfunctions.net/GenerateRandomNumber"

  t.Run("Test Log", func(t *testing.T){
    for i := 10; i < 100009; i*= 10 {
      var urlReq = url + "?max=" + strconv.Itoa(i)

      resp, err := http.Get(urlReq )
      if err != nil {
        t.Errorf("Failed to GET url request: Error %v\n", err.Error())
      }
      // Read the response body
      body, err := ioutil.ReadAll(resp.Body)
      if err != nil {
        t.Errorf("Failed to Read The Body: Error %v\n", err.Error())
      }

      var response Response
      if err := json.Unmarshal(body, &response); err != nil {
        t.Errorf("Failed to Unmarshal the Response Body: Error %v\n", err.Error())
      }
      number := response.Number
      if number > i || number < 1 {
        t.Errorf("Random Generator failed to produce a number between 1 & %v", i)
      }else {
        t.Logf("Test Passed :: 1 < %v < %v", number, i)
        fmt.Printf("Test Passed :: 1 < %v < %v\n", number, i)
      }
    }
  })

  t.Run("Num Array", func(t *testing.T){
    tests := []struct{
      max int
    }{
      {max: 20},
      {max: 50},
      {max: 100},
    }

    for _, tt := range tests {
      var urlReq = url + "?max=" + strconv.Itoa(tt.max)

      resp, err := http.Get(urlReq )
      if err != nil {
        t.Errorf("Failed to GET url request: Error %v\n", err.Error())
      }
      // Read the response body
      body, err := ioutil.ReadAll(resp.Body)
      if err != nil {
        t.Errorf("Failed to Read The Body: Error %v\n", err.Error())
      }

      var response Response
      if err := json.Unmarshal(body, &response); err != nil {
        t.Errorf("Failed to Unmarshal the Response Body: Error %v\n", err.Error())
      }
      number := response.Number
      if number > tt.max || number < 1 {
        t.Errorf("Random Generator failed to produce a number between 1 & %v", tt.max)
      }else {
        t.Logf("Test Passed :: 1 < %v < %v", number, tt.max)
        fmt.Printf("Test Passed :: 1 < %v < %v\n", number, tt.max)
      }
    }
  })
}
