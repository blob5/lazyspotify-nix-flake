package player

import (
	"context"
	"fmt"
	"time"

	"github.com/dubeyKartikay/lazyspotify/librespot"
)

func PlayTrack(ctx context.Context, uri string) error {
	l, err := librespot.InitLibrespot(true)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if err := l.Deamon.StartDeamon(); err != nil {
		return err
	}
	defer l.Deamon.StopDeamon()

	select {
	case err := <-l.Ready:
		if err != nil {
			return err
		}
	case err := <-l.Deamon.RestartFailErrorChannel:
		return fmt.Errorf("daemon is no longer available: %w", err)
	case <-ctx.Done():
		return ctx.Err()
	}
	fmt.Println(uri)
	res := l.Client.Play(ctx, uri, "", false)
	if res >= 400 {
		return fmt.Errorf("failed to play track: daemon returned status %d", res)
	}

	select {
	case err := <-l.Deamon.RestartFailErrorChannel:
		return fmt.Errorf("daemon is no longer available: %w", err)
	case <-time.After(30 * time.Second):
	case <-ctx.Done():
		return ctx.Err()
	}
	fmt.Println(res)
	return nil

}
