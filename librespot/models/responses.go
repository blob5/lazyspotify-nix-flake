package models

type HealthResponse struct {
	PlaybackReady bool `json:"playback_ready"`
}

type VolumeResponse struct {
	Value int `json:"value"`
	Max   int `json:"max"`
}
