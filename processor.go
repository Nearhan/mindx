package mindx

import (
	"log"

	"github.com/Nearhan/mindx/proto"
)

var _ proto.AgentServer = (*Processor)(nil)

// Processor handles data points
type Processor struct {
	Inserter Inserter
}

// NewProcessor create a new Processor
func NewProcessor(i Inserter) *Processor {
	return &Processor{Inserter: i}
}

// Process reads off the grpc stream
func (p *Processor) Process(stream proto.Agent_ProcessServer) error {

	// read off grpc stream
	for {
		dp, err := stream.Recv()
		if err != nil {
			return err
		}

		// batch insert
		if err := p.Inserter.InsertBatch(dp); err != nil {
			log.Println("unable to insert dp due to ", err)
			return err
		}
	}

	return nil
}
