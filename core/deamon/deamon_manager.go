package deamon

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
)

type DeamonManager struct {
	deamonProcess           DeamonProcess
	cmd                     []string
	restartOnFailure        bool
	restartCount            int
	deamonErrorChannel      chan error
	RestartFailErrorChannel chan error
	stopping                atomic.Bool
}

func NewDeamonManager(cmd []string) (DeamonManager, error) {
	deamon, err := NewDeamonProcess(context.Background(), cmd)
	if err != nil {
		return DeamonManager{}, err
	}
	return DeamonManager{
		deamonProcess:           deamon,
		cmd:                     cmd,
		restartOnFailure:        true,
		deamonErrorChannel:      make(chan error, 1),
		RestartFailErrorChannel: make(chan error, 1),
	}, nil
}

func (d *DeamonManager) StartDeamon() error {
	fmt.Println("Starting Deamon")
	d.stopping.Store(false)
	err := d.deamonProcess.StartDeamon()
	if err != nil {
		return err
	}
	go d.deamonProcess.MonitorDeamon(d.deamonErrorChannel)
	go d.listenForErrors()
	return nil
}

func (d *DeamonManager) RestartDeamon() error {
	d.StopDeamon()
	deamon, err := NewDeamonProcess(context.Background(), d.cmd)
	d.deamonProcess = deamon
	if err != nil {
		return err
	}
	return d.StartDeamon()
}

func (d *DeamonManager) StopDeamon() {
	if d.deamonProcess.cmd.Process == nil {
		return
	}
	d.stopping.Store(true)
	err := d.deamonProcess.cmd.Process.Signal(os.Interrupt)
	if err != nil {
		d.forceKill()
	}
	//  time.Sleep(5 * time.Second)
	// fmt.Printf("%+v",d.deamonProcess.cmd.Process)
	// fmt.Printf("%+v",d.deamonProcess.cmd.ProcessState)
}

func (d *DeamonManager) listenForErrors() {
	err := <-d.deamonErrorChannel
	if d.stopping.Load() {
		return
	}
	fmt.Println("deamon error:", err)
	fmt.Printf("%+v", d)
	if !d.restartOnFailure {
		d.StopDeamon()
		d.reportRestartFailure(err)
		return
	}
	if d.restartCount >= 3 {
		d.reportRestartFailure(fmt.Errorf("max daemon retry breached: %w", err))
		return
	}
	d.restartCount++
	err = d.RestartDeamon()
	if err != nil {
		d.reportRestartFailure(fmt.Errorf("failed to restart daemon: %w", err))
		return
	}
}

func (d *DeamonManager) reportRestartFailure(err error) {
	select {
	case d.RestartFailErrorChannel <- err:
	default:
	}
}
func (d *DeamonManager) forceKill() {
	fmt.Println("Force killing process...")
	if err := d.deamonProcess.cmd.Process.Kill(); err != nil {
		fmt.Printf("Failed to kill process: %v\n", err)
	}
	if d.deamonProcess.cancel != nil {
		d.deamonProcess.cancel()
	}
}
