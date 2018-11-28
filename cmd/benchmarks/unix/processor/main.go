package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {

	defer func() {
		if err := os.Remove("/tmp/example.sock"); err != nil {
			log.Println(err)
		}
	}()

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
}
