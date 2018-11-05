package main

import (
	"context"
	"log"
	"time"

	"github.com/Nearhan/mindx"
	"github.com/Nearhan/mindx/proto"
	"github.com/caarlos0/env"
	"google.golang.org/grpc"
)

type config struct {
	SendAddr string `env:"SEND_ADDR" envDefault:"processor:5001"`
}

func main() {

	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("unable to parse config", err)
	}

	conn, err := grpc.Dial(cfg.SendAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second*5))

	if err != nil {
		log.Fatal(err)
	}

	log.Println("connection to processor established")

	// create proto client
	client := proto.NewAgentClient(conn)
	streamer, err := client.Process(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// spin off go routine to send data
	go func() {
		log.Println("entering data generation Loop")
		for {
			dp := mindx.RandomizeDP()
			err := streamer.Send(&dp)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	done := mindx.HandleSignal()
	<-done
	log.Println("exiting generator")
}
