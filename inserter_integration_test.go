// +build integration

package mindx

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/Nearhan/mindx/proto"
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// Runs a generator and processor which communicate over grpc
// and inserts two batches into influxdb
func TestImplEnd2End(t *testing.T) {

	// make done channel
	d1 := make(chan bool)
	expectedBatch := 1000
	test_db := "test_mindx"

	// make influx client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})

	// drop data base when its all over
	defer func() {
		_, err = queryDB(c, fmt.Sprintf("DROP DATABASE %s", test_db))
		if err != nil {
			log.Fatal(err)
		}
	}()

	// setup dp generator
	go func() {
		conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure(), grpc.WithTimeout(time.Second*5))
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

		count := 1000
		log.Println("enter Data Generation Loop")
		for count > 0 {
			dp := RandomizeDP()
			err := streamer.Send(&dp)
			if err != nil {
				fmt.Println(err)
			}
			count--
		}
		streamer.CloseAndRecv()
		d1 <- true
	}()

	// setup processor
	go func() {
		// create influx db inserter
		inserter, err := NewInfluxUDPInserter(test_db, "localhost:8089", expectedBatch)
		if err != nil {
			log.Fatal(err)
		}

		// create Processor
		processor := NewProcessor(inserter)

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 5001))
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

	}()

	<-d1
	// wait some arbitrary amount of time
	// for process to finish inserting into db
	time.Sleep(2 * time.Second)

	// query for the count
	// should equal expected batch #
	q := fmt.Sprintf("SELECT count(*) FROM %s;", "data_point")
	res, err := queryDB(c, q)
	assert.NoError(t, err)
	count := res[0].Series[0].Values[0][1]
	i, err := count.(json.Number).Int64()
	assert.NoError(t, err)
	assert.Equal(t, expectedBatch, int(i))
}

// queryDB convenience function to query the database
func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: "test_mindx",
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
