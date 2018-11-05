# mindx code challenge

Here is my implementation of the code challenge.

I used following and go into further detail on to why I made those decisions below:

1. Go as the programing langauge
2. GRPC + Protobuf for IPC layer 
3. INFLUXDB as backing store
4. Docker for easy deployment and testing


### Quick Start

1. make sure you have docker on your machine
2. make sure you have docker-compose on your machine 
    ```
    brew install docker-compose
    ```
3. install the influxdb cli 
    ```
    brew install influxdb
    ```
4. Up all the services

    ```
    docker-compose up
    ```
5. spin up the influxdb cli to make sure you get data points. Follow the commands below

```
❯ influx
> use mindx
> select count(*) from "data_point"
name: data-point
time count_confidence count_id count_norm_pos_x count_norm_pos_y count_pupil_diameter count_timestamp
---- ---------------- -------- ---------------- ---------------- -------------------- ---------------
0    454997           454997   454997           459997           459997               459997
```


### Not Quick Start

If you want to compile everything you'll need a fair amount of things.

1. Install GO https://golang.org/doc/install
2. Install Protobuf and GRPC https://github.com/protocolbuffers/protobuf/releases/tag/v3.6.1
3. Install dep
    ```
    brew install dep
    ```
4. Install go dependencies
    ```
    dep ensure -v
    ```
5. Install go grpc plugin  
    ```
    go get -u google.golang.org/grpc
    ```
6. Build everything!
    ```
    ./build.sh
    ```
    this compiles the protobuf files, compiles the binaries and creates docker images 
7. Test stuff!
    ```
    go test -v 
    ```
8. Run it all
   ```
   docker-compose -f docker-compose-local.yml up
   ```


## Considerations

### InfluxDB

Influxdb is a time series database, has a udp protocol and can store timestamps as IDS!
It allows for batch inserting and is distributed from the start. 
It also allows for you to query for data points with a sql like language.

### Inserter Interface

I used an interface for how I store data points so that I swap this out whenever I wanted.
For instance if I wanted to switch to using a postgres database I would just have a postgres implementation
of this interface.

Currently It has two methods, for this challenge I created a influxdb udp inserter

```
// Inserter handles how we store and insert data points
type Inserter interface {
	Insert(*proto.DataPoint) error
	InsertBatch(*proto.DataPoint) error
}
```


### IPC && BenchMark tests

I was curious to try various IPC methods to test which one was faster.
Each test uses a different IPC protocol and sends 10k messages for processing.
They all clocked in around the same time so I decided to go with GRPC since its modern, portable and allows for a very flexible API (Easy Bi Directional Streaming).


For all methods I was using protobuf as the encoding layer.

You can also run the benchmark tests as well 
```
❯ go test -v
=== RUN   TestImpl
--- PASS: TestImpl (0.11s)
    benchmarks_test.go:62: GRPC took 36.440859ms
    benchmarks_test.go:62: UNIX SOCKETS took 35.151095ms
    benchmarks_test.go:62: UNIX PIPES took 35.69541ms
=== RUN   TestConvertDpInfluxPoint
--- PASS: TestConvertDpInfluxPoint (0.00s)
=== RUN   TestCreateNewBatchPoints
--- PASS: TestCreateNewBatchPoints (0.00s)
PASS
ok      github.com/Nearhan/mindx        0.120s
```

### Integration Tests

If you wish to run the integration test you'll need to run the other docker-compose file

```
docker-compose -f docker-compose-integration-test.yml up
```

Then in another shell run

```
go test -v -tags=integration
```

This runs an end 2 end tests; fully sending over 1000 messages
Checks influxdb to see if the correct # of messages appear