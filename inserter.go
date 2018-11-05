package mindx

import (
	"fmt"
	"log"
	"time"

	"github.com/Nearhan/mindx/proto"
	"github.com/influxdata/influxdb/client/v2"
)

// Inserter handles how we store and insert data points
type Inserter interface {
	Insert(*proto.DataPoint) error
	InsertBatch(*proto.DataPoint) error
}

// interface check
var _ (Inserter) = (*InfluxUDPInserter)(nil)

// InfluxUDPInserter an influx db & udp implementation of Inserter interface
type InfluxUDPInserter struct {
	client.Client // Embedded Influx DB Client
	Points        client.BatchPoints
	BatchSize     int
	DBName        string
}

// NewInfluxUDPInserter constructor for influxdb
func NewInfluxUDPInserter(dbName, dbConnStr string, batchSize int) (*InfluxUDPInserter, error) {
	var i InfluxUDPInserter
	c, err := client.NewUDPClient(client.UDPConfig{Addr: dbConnStr})

	if err != nil {
		return nil, err
	}

	// Create a new point batch
	bps, err := CreateNewBatchPoints(dbName)
	if err != nil {
		return nil, err
	}

	i.Points = bps
	i.Client = c
	i.DBName = dbName
	i.BatchSize = batchSize
	return &i, nil
}

// Insert adds single dp to influx
func (i *InfluxUDPInserter) Insert(dp *proto.DataPoint) error {
	log.Fatal("not implemented yet")
	return nil
}

// InsertBatch adds dp to influx as a batch
func (i *InfluxUDPInserter) InsertBatch(dp *proto.DataPoint) error {
	// convert dp to influx point
	point, err := ConvertDpInfluxPoint(dp)
	if err != nil {
		fmt.Println("erroring out when converting to point")
		return err
	}

	i.Points.AddPoint(point)

	// check if we should batch
	if len(i.Points.Points()) >= i.BatchSize {
		if err := i.Client.Write(i.Points); err != nil {
			log.Println("error when writing to influx", err)
			return err
		}

		// create new bps
		bps, err := CreateNewBatchPoints(i.DBName)
		if err != nil {
			return err
		}

		i.Points = bps
		return nil
	}

	return nil
}

// ConvertDpInfluxPoint takes a grpc datapoint and converts it into InfluxDB type
func ConvertDpInfluxPoint(dp *proto.DataPoint) (*client.Point, error) {

	ts := dp.GetTimestamp()
	fields := map[string]interface{}{
		"id":             dp.GetId(),
		"confidence":     dp.GetConfidence(),
		"norm_pos_x":     dp.GetNormPosX(),
		"norm_pos_y":     dp.GetNormPosY(),
		"pupil_diameter": dp.GetPupilDiameter(),
		"timestamp":      ts.String(),
	}

	// set id of this datapoint to be the timestamp itself
	id := time.Unix(ts.GetSeconds(), int64(ts.GetNanos()))
	return client.NewPoint("data_point", nil, fields, id)
}

// CreateNewBatchPoints creates a new container for influxdb points
func CreateNewBatchPoints(db string) (client.BatchPoints, error) {
	// Create a new point batch
	bps, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: "ns",
	})
	if err != nil {
		return nil, err
	}
	return bps, nil
}
