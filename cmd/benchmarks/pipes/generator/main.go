package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/Nearhan/mindx/proto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func main() {
	pipeFile := "/tmp/namedPipe"
	os.Remove(pipeFile)

	err := syscall.Mkfifo(pipeFile, 0666)
	if err != nil {
		log.Fatal("Make named pipe file error:", err)
	}
	file, err := os.OpenFile(pipeFile, os.O_CREATE, 0600)

	defer file.Close()
	if err != nil {
		panic(err)
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
		_, err = file.Write(b.Bytes())
		if err != nil {
			log.Println(err)
		}
		count--
	}

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
