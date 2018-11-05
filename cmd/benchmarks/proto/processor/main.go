package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/Nearhan/mindx/proto"
	"google.golang.org/grpc"
)

var _ AgentServer = (*Processor)(nil)

type Processor struct {
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) Process(stream Agent_ProcessServer) error {
	count := 0
	for {
		_, err := stream.Recv()
		count++
		if err == io.EOF {
			fmt.Println(err)
			fmt.Println(count)
			p.server.Stop()
			return nil
		}
	}

	return nil
}

func main() {

	grpcServer := grpc.NewServer()
	processor := NewProcessor()

	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatal(err)
	}

	proto.RegisterAgentServer(grpcServer, processor)
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatal(err)
	}

}
