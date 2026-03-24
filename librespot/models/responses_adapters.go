package models

import "encoding/json"

func DecodeHealthResponse(data []byte) (HealthResponse, error) {
	var healthResponse HealthResponse
	err := json.Unmarshal(data, &healthResponse)
	return healthResponse, err
}

func DecodeVolumeResponse(data []byte) (VolumeResponse, error) {
	var volumeResponse VolumeResponse
	err := json.Unmarshal(data, &volumeResponse)
	return volumeResponse, err
}
