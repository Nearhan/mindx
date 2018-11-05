package main

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/Nearhan/mindx"
)

func main() {

	done := mindx.HandleSignal()

	defer func() {
		if err := os.Remove("/tmp/example.sock"); err != nil {
			log.Println(err)
		}
	}()

	go func(d chan bool) {

		defer func() { d <- true }()

		l, err := net.Listen("unix", "/tmp/example.sock")
		defer func() {
			if l != nil {
				l.Close()
			}
		}()

		if err != nil {
			log.Println(err)
		}

		fd, err := l.Accept()
		if err != nil {
			log.Println(err)
		}

		count := 0
		for {
			buf := make([]byte, 107)
			_, err := fd.Read(buf)
			if err == io.EOF {
				fd.Close()
				log.Println("messages consumed", count)
				return
			}
			count++
		}

	}(done)
	<-done
}
