package deamon

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type DeamonProcess struct {
	cmd *exec.Cmd
	cancel context.CancelFunc
}

func NewDeamonProcess(ctx context.Context, args []string)(DeamonProcess,error){
  cmd := exec.CommandContext(ctx, args[0], args[1:]...)
  ctx, cancel := context.WithCancel(ctx)
  return DeamonProcess{cmd: cmd, cancel: cancel}, nil
}

func (d *DeamonProcess) StartDeamon() error {
	logFile, err := os.OpenFile("daemon.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666) // TODO: create deamon log in config DIR
	if err != nil {
		d.cancel()
		return fmt.Errorf("could not open log file: %w", err)
	}
	d.cmd.Stdout = logFile
	d.cmd.Stderr = logFile

	if err := d.cmd.Start(); err != nil {
		d.cancel()
		return fmt.Errorf("failed to start daemon: %w", err)
	}
	fmt.Println("deamon process", d.cmd.Process)
	return nil
}

func (d *DeamonProcess) MonitorDeamon(channel chan error){
	err := d.cmd.Wait()
	if (err != nil){
    channel <- err
	}
}


