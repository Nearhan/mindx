package mindx

import (
	"math/rand"
	"time"

	"github.com/Nearhan/mindx/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

// RandomizeDP randomizes a Data Point
func RandomizeDP() proto.DataPoint {
	rand.Seed(time.Now().UTC().UnixNano())
	t := time.Now()
	id := int32(rand.Intn(2))
	normx := rand.Float32()
	normy := rand.Float32()
	confidence := rand.Float32()
	pupil := rand.Int31()

	dp := proto.DataPoint{
		Timestamp: &timestamp.Timestamp{
			Seconds: int64(t.Second()),
			Nanos:   int32(t.Nanosecond()),
		},
		Id:            id,
		NormPosX:      normx,
		NormPosY:      normy,
		Confidence:    confidence,
		PupilDiameter: pupil,
	}
	return dp
}
