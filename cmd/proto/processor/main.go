package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Nearhan/mindx"
	"github.com/Nearhan/mindx/proto"
	"github.com/caarlos0/env"
	"google.golang.org/grpc"
)

// config type
type config struct {
	BatchSize int    `env:"BATCH_SIZE" envDefault:"500"`
	Port      int    `env:"PORT" envDefault:"5001"`
	DbName    string `env:"DB_NAME" envDefault:"mindx"`
	DbConn    string `env:"DB_CONN_STR" envDefault:"db:8089"`
}

func main() {

	// parse config
	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("unable to parse config", err)
	}

	// create influx db inserter
	inserter, err := mindx.NewInfluxUDPInserter(cfg.DbName, cfg.DbConn, cfg.BatchSize)
	if err != nil {
		log.Fatal(err)
	}

	// create Processor
	processor := mindx.NewProcessor(inserter)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterAgentServer(grpcServer, processor)
	log.Println("starting processor")
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatal(err)
	}

	// waiting for process to end
	done := mindx.HandleSignal()
	<-done
	log.Println("exiting processor")
}
