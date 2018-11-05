package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Nearhan/mindx/proto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func main() {
	c, err := net.Dial("unix", "/tmp/example.sock")

	if err != nil {
		log.Println(err)
		return
	}

	count := 10000
	dp := GetDP()
	b := new(bytes.Buffer)
	m := jsonpb.Marshaler{EnumsAsInts: false, EmitDefaults: true}
	err = m.Marshal(b, &dp)
	fmt.Println(b.Len())
	if err != nil {
		panic(err)
	}

	for count > 0 {
		_, err = c.Write(b.Bytes())
		if err != nil {
			log.Println(err)
		}
		count--
	}
	c.Close()

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
