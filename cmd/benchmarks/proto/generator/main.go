package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Nearhan/mindx/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
)

var send = flag.String("send", "localhost:5001", "where to send datapoints to")

func main() {

	flag.Parse()
	conn, err := grpc.Dial(*send, grpc.WithInsecure(), grpc.WithTimeout(time.Second*5))
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

	count := 10000
	log.Println("enter Data Generation Loop")
	for count > 0 {
		dp := GetDP()
		err := streamer.Send(&dp)
		if err != nil {
			fmt.Println(err)
		}
		count--
	}
	x, err := streamer.CloseAndRecv()
	fmt.Println(x, err)
	log.Println("done")
}

// RandomizeDP randomizes a Data Point
func GetDP() proto.DataPoint {
	t := time.Time{}
	id := int32(0)
	normx := float32(.1)
	normy := float32(.1)
	confidence := float32(.1)
	pupil := int32(1)

	dp := proto.DataPoint{
		Timestamp: &timestamp.Timestamp{
			Seconds: int64(t.Second()),
			Nanos:   int32(t.Nanosecond()),
		},
		Id:           id,
		NormPosX:     normx,
		NormPosY:     normy,
		Confidence:   confidence,
		PupilDiameter: pupil,
	}
	return dp
}
