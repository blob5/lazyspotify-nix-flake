package models

import "encoding/json"

type HealthResponse struct {
	PlaybackReady bool `json:"playback_ready"`
}

type VolumeResponse struct {
	Value int `json:"value"`
	Max   int `json:"max"`
}

type EventType string

const (
	EventTypeMetadata EventType = "metadata"
	EventTypePlaying  EventType = "playing"
	EventTypePaused   EventType = "paused"
	EventTypeStopped  EventType = "stopped"
	EventTypeSeek     EventType = "seek"
	EventTypeVolume   EventType = "volume"
)

type EventEnvelope struct {
	Type EventType       `json:"type"`
	Data json.RawMessage `json:"data"`
}

type MetadataEventData struct {
	ContextURI    string   `json:"context_uri"`
	URI           string   `json:"uri"`
	Name          string   `json:"name"`
	ArtistNames   []string `json:"artist_names"`
	AlbumName     string   `json:"album_name"`
	AlbumCoverURL *string  `json:"album_cover_url"`
	Position      int      `json:"position"`
	Duration      int      `json:"duration"`
}

type PlayingEventData struct {
	ContextURI string `json:"context_uri"`
	URI        string `json:"uri"`
	Resume     bool   `json:"resume"`
	PlayOrigin string `json:"play_origin"`
}

type PausedEventData struct {
	ContextURI string `json:"context_uri"`
	URI        string `json:"uri"`
	PlayOrigin string `json:"play_origin"`
}

type StoppedEventData struct {
	PlayOrigin string `json:"play_origin"`
}

type SeekEventData struct {
	ContextURI string `json:"context_uri"`
	URI        string `json:"uri"`
	Position   int    `json:"position"`
	Duration   int    `json:"duration"`
	PlayOrigin string `json:"play_origin"`
}

type VolumeEventData struct {
	Value int `json:"value"`
	Max   int `json:"max"`
}

type PlayerEvent struct {
	Type     EventType
	Metadata *MetadataEventData
	Playing  *PlayingEventData
	Paused   *PausedEventData
	Stopped  *StoppedEventData
	Seek     *SeekEventData
	Volume   *VolumeEventData
}
