package librespot

import (
	"fmt"
	"time"

	"github.com/dubeyKartikay/lazyspotify/core/deamon"
)

type Librespot struct {
	Deamon deamon.DeamonManager
	Server *LibrespotApiServer
	Client *LibrespotApiClient
	Ready  chan error
}

func InitLibrespot(panicOnDaemonFailure bool) (*Librespot, error) {
	deamonManager, err := deamon.NewDeamonManager([]string{"/Users/user/personal/go-librespot/daemon"})
	if err != nil {
		return nil, err
	}

	librespotApiServer := NewLibrespotApiServer("127.0.0.1", "4040")
	librespotApiClient := NewLibrespotApiClient(librespotApiServer)
	l := &Librespot{Deamon: deamonManager, Server: librespotApiServer, Client: librespotApiClient, Ready: make(chan error, 1)}
	go notifyWhenReady(l)
	return l, nil
}

func notifyWhenReady(l *Librespot) {
	for range 900 {
		healthRes, err := l.Client.GetHealth()

		if err == nil && healthRes.PlaybackReady {
			l.Ready <- nil
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	l.Ready <- fmt.Errorf("daemon did not become ready before timeout")
}
