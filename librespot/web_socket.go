package librespot

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/dubeyKartikay/lazyspotify/core/logger"
	"github.com/dubeyKartikay/lazyspotify/librespot/models"
)

type eventSocket struct {
	serverURL string
	events    chan models.PlayerEvent
	done      chan struct{}
}

func newEventSocket(serverURL string) *eventSocket {
	return &eventSocket{
		serverURL: serverURL,
		events:    make(chan models.PlayerEvent, 32),
		done:      make(chan struct{}),
	}
}

func (s *eventSocket) Events() <-chan models.PlayerEvent {
	return s.events
}

func (s *eventSocket) Start() {
	go s.run()
}

func (s *eventSocket) Close() {
	select {
	case <-s.done:
		return
	default:
		close(s.done)
	}
}

func (s *eventSocket) run() {
	defer close(s.events)

	backoff := 200 * time.Millisecond
	for {
		select {
		case <-s.done:
			return
		default:
		}

		err := s.connectAndRead()
		if err != nil {
			logger.Log.Warn().Err(err).Msg("event websocket disconnected")
			if closeStatus := websocket.CloseStatus(err); closeStatus == websocket.StatusNormalClosure || closeStatus == websocket.StatusGoingAway {
				return
			}
		}

		select {
		case <-s.done:
			return
		case <-time.After(backoff):
			if backoff < 3*time.Second {
				backoff *= 2
				if backoff > 3*time.Second {
					backoff = 3 * time.Second
				}
			}
		}
	}
}

func (s *eventSocket) connectAndRead() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	wsURL := toWebSocketURL(s.serverURL) + "/events"
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close(websocket.StatusNormalClosure, "")
	logger.Log.Info().Str("url", wsURL).Msg("connected to event websocket")

	for {
		select {
		case <-s.done:
			return nil
		default:
		}

		var envelope models.EventEnvelope
		if err := wsjson.Read(context.Background(), conn, &envelope); err != nil {
			return err
		}

		eventBytes, err := json.Marshal(envelope)
		if err != nil {
			logger.Log.Warn().Err(err).Msg("failed to marshal websocket event envelope")
			continue
		}

		event, err := models.DecodePlayerEvent(eventBytes)
		if err != nil {
			logger.Log.Debug().Err(err).Str("type", string(envelope.Type)).Msg("ignoring unsupported websocket event")
			continue
		}

		select {
		case s.events <- event:
		case <-s.done:
			return nil
		}
	}
}

func toWebSocketURL(serverURL string) string {
	if strings.HasPrefix(serverURL, "https://") {
		return "wss://" + strings.TrimPrefix(serverURL, "https://")
	}
	if strings.HasPrefix(serverURL, "http://") {
		return "ws://" + strings.TrimPrefix(serverURL, "http://")
	}
	return fmt.Sprintf("ws://%s", serverURL)
}
