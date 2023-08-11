package randomGen

import (
	// "fmt"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	// "github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// func init(){
//   functions.HTTP("RandomGet", randomGen)
// }
//

type Response struct {
  Number int `json:"number"`
}

func RandomGen(w http.ResponseWriter, r *http.Request){
  rand.Seed(time.Now().UnixNano())

  maxParam := r.URL.Query().Get("max")
  maxValue := 20

  if maxParam != "" {
    var err error
    maxValue, err = strconv.Atoi(maxParam)
    if err != nil || maxValue < 2 {
      http.Error(w, "Invalid Max Value, please provide an integer greater than 1", http.StatusBadRequest)
      return
    }
  }
  randomInt := rand.Intn(maxValue) + 1

  response := Response{
    Number:randomInt,
  }

  w.Header().Set("Content-Type", "application/json")
  if err := json.NewEncoder(w).Encode(response); err != nil {
    http.Error(w, "Error Generating Random Number:", http.StatusInternalServerError)
  }
}

