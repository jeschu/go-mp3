package lang

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func WaitForSignal(signals ...os.Signal) {
	quit := make(chan os.Signal)
	signal.Notify(quit, signals...)
	s := <-quit
	log.Printf("got signal %v", s)
}

func WaitForIntOrTerm() { WaitForSignal(syscall.SIGINT, syscall.SIGTERM) }
