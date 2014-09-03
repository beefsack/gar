package gar

import (
	"os"
	"os/signal"
)

func onExit(f func()) {
	go func() {
		s := make(chan os.Signal)
		signal.Notify(s, os.Interrupt, os.Kill)
		<-s
		f()
	}()
}
