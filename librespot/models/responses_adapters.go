package models

import (
	"encoding/json"
	"fmt"
)

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

func DecodeResolveTracksResponse(data []byte) (ResolveTracksResponse, error) {
	var resolveTracksResponse ResolveTracksResponse
	err := json.Unmarshal(data, &resolveTracksResponse)
	return resolveTracksResponse, err
}

func DecodeEventEnvelope(data []byte) (EventEnvelope, error) {
	var envelope EventEnvelope
	err := json.Unmarshal(data, &envelope)
	return envelope, err
}

func DecodePlayerEvent(data []byte) (PlayerEvent, error) {
	envelope, err := DecodeEventEnvelope(data)
	if err != nil {
		return PlayerEvent{}, err
	}

	event := PlayerEvent{Type: envelope.Type}
	switch envelope.Type {
	case EventTypeMetadata:
		var payload MetadataEventData
		if err := json.Unmarshal(envelope.Data, &payload); err != nil {
			return PlayerEvent{}, err
		}
		event.Metadata = &payload
	case EventTypePlaying:
		var payload PlayingEventData
		if err := json.Unmarshal(envelope.Data, &payload); err != nil {
			return PlayerEvent{}, err
		}
		event.Playing = &payload
	case EventTypePaused:
		var payload PausedEventData
		if err := json.Unmarshal(envelope.Data, &payload); err != nil {
			return PlayerEvent{}, err
		}
		event.Paused = &payload
	case EventTypeStopped:
		var payload StoppedEventData
		if err := json.Unmarshal(envelope.Data, &payload); err != nil {
			return PlayerEvent{}, err
		}
		event.Stopped = &payload
	case EventTypeSeek:
		var payload SeekEventData
		if err := json.Unmarshal(envelope.Data, &payload); err != nil {
			return PlayerEvent{}, err
		}
		event.Seek = &payload
	case EventTypeVolume:
		var payload VolumeEventData
		if err := json.Unmarshal(envelope.Data, &payload); err != nil {
			return PlayerEvent{}, err
		}
		event.Volume = &payload
	default:
		return PlayerEvent{}, fmt.Errorf("unsupported event type: %s", envelope.Type)
	}

	return event, nil
}
