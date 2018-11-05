package mindx

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// HandleSignal handles signals
func HandleSignal() chan bool {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("exiting")
		done <- true
	}()

	return done
}
