package models

import "encoding/json"


func DecodeHealthResponse(data []byte) (HealthResponse, error) {
  var healthResponse HealthResponse
  err := json.Unmarshal(data, &healthResponse)
  return healthResponse, err
}
