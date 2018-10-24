package pressctrlc

import (
	"os"
	"os/signal"
)

// ToExit blocks until os.Interrupt signal is received. This usually happens when CTRL+C is pressed
// This is so small you should probably copy paste to your project instead of importing
func ToExit() {
	exitCh := make(chan struct{})
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		select {
		case <-signalCh:
		}
		exitCh <- struct{}{}
	}()
	<-exitCh
}
