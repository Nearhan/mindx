package main

import (
	"io"
	"log"
	"os"

	"github.com/Nearhan/mindx"
)

func main() {

	done := mindx.HandleSignal()

	defer func() {
		if err := os.Remove("/tmp/namedPipe"); err != nil {
			log.Println(err)
		}
	}()

	go func(d chan bool) {

		defer func() { d <- true }()

		f, err := os.OpenFile("/tmp/namedPipe", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
		defer func() {
			if f != nil {
				f.Close()
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
