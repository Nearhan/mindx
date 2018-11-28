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

	f, err := os.OpenFile("/tmp/namedPipe", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	defer func() {
		if f != nil {
			f.Close()
		}
	}()

	if err != nil {
		log.Println(err)
	}

	count := 0
	for {
		buf := make([]byte, 107)
		_, err := f.Read(buf)
		if err == io.EOF {
			f.Close()
			log.Println("messages consumed", count)
			return
		}
		count++
	}

}
