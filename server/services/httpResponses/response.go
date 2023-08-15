package httpResponses

import "encoding/json"

type Response struct {
  Status     string  `json:"status"`
  StatusCode int32   `json:"statusCode"`
  Body       map[string]interface{}  `json:"body"`
}

func(resp *Response)AddStatus(status string) *Response {
  resp.Status = status
  return resp
}

func(resp *Response)AddStatusCode(sc int32) *Response {
  resp.StatusCode = sc
  return resp
}

func(resp *Response)AddBody(body map[string]interface{}) *Response {
  resp.Body = body
  return resp
}

func(resp *Response)Build()( []byte, error ){
  return json.Marshal(resp)
}
