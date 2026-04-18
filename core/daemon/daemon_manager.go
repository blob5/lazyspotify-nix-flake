package daemon

import (
	"context"
	"fmt"
	"github.com/dubeyKartikay/lazyspotify/core/logger"
	"os"
)

type DaemonManager struct {
	daemonProcess           DaemonProcess
	cmd                     []string
	restartOnFailure        bool
	restartCount            int
	daemonErrorChannel      chan error
	RestartFailErrorChannel chan error
}

func NewDaemonManager(cmd []string) (DaemonManager, error) {
	daemon, err := NewDaemonProcess(context.Background(), cmd)
	if err != nil {
		return DaemonManager{}, err
	}
	return DaemonManager{
		daemonProcess:           daemon,
		cmd:                     cmd,
		restartOnFailure:        true,
		daemonErrorChannel:      make(chan error, 1),
		RestartFailErrorChannel: make(chan error, 1),
	}, nil
}

func (d *DaemonManager) StartDaemon() error {
	logger.Log.Info().Msg("starting daemon")
	d.restartOnFailure = true
	err := d.daemonProcess.StartDaemon()
	if err != nil {
		return err
	}
	go d.daemonProcess.MonitorDaemon(d.daemonErrorChannel)
	go d.listenForErrors()
	return nil
}

func (d *DaemonManager) RestartDaemon() error {
	d.StopDaemon()
	daemon, err := NewDaemonProcess(context.Background(), d.cmd)
	d.daemonProcess = daemon
	if err != nil {
		return err
	}
	return d.StartDaemon()
}

func (d *DaemonManager) StopDaemon() {
	d.restartOnFailure = false
	if d.daemonProcess.cmd.Process == nil {
		return
	}
	err := d.daemonProcess.cmd.Process.Signal(os.Interrupt)
	if err != nil {
		d.forceKill()
	}
}

func (d *DaemonManager) listenForErrors() {
	err := <-d.daemonErrorChannel
	logger.Log.Error().Err(err).Msgf("daemon error: %+v", d)
	if !d.restartOnFailure {
		return
	}
	if d.restartCount >= 3 {
		d.reportRestartFailure(fmt.Errorf("max daemon retry breached: %w", err))
		return
	}
	d.restartCount++
	err = d.RestartDaemon()
	if err != nil {
		d.reportRestartFailure(fmt.Errorf("failed to restart daemon: %w", err))
		return
	}
}

func (d *DaemonManager) reportRestartFailure(err error) {
	select {
	case d.RestartFailErrorChannel <- err:
	default:
	}
}

func (d *DaemonManager) forceKill() {
	logger.Log.Warn().Msg("force killing process")
	if err := d.daemonProcess.cmd.Process.Kill(); err != nil {
		logger.Log.Error().Err(err).Msg("failed to kill process")
	}
	if d.daemonProcess.cancel != nil {
		d.daemonProcess.cancel()
	}
}
