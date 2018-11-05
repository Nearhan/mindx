# mindx code challenge

### quick start

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


### Not Quickstart

If you want to compile everything you'll need a fair amount of things.

1. Install GO https://golang.org/doc/install
2. Install Protobuff and GRPC https://github.com/protocolbuffers/protobuf/releases/tag/v3.6.1
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

7. Run it all
   ```
   docker-compose -f docker-compose-local.yml up
   ```




### InfluxDB

I decided to go with influxdb as the backing store.
Why?
Influxdb is a time series database, has a udp protocol and can store timestamps as IDS!
It also allows you to query for datapoints with a sql like query lanague.

### Inserter Interface

I used an interface for how I store data points so that I swap this out whenever I watned.
Currently It has two methods, and I implemented an influxdb udp variant

```
// Inserter handles how we store and insert data points
type Inserter interface {
	Insert(*proto.DataPoint) error
	InsertBatch(*proto.DataPoint) error
}
````


### BenchMark tests

I was curious to try various IPC modes to test which one was faster.
Each test uses a different IPC protocol and sends 10k messages for processing.
They all clocked in around the same time so I decided to go with GRPC  + Protobuf since its modern, portable and 
allows for a very flexible API (you can keep adding fields and you don't need to do a schema change)
```
~/Code/golang/src/github.com/Nearhan/mindx master
❯ go test -v
=== RUN   TestImpl
--- PASS: TestImpl (0.11s)
    benchmarks_test.go:62: GRPC took 36.696675ms
    benchmarks_test.go:62: UNIX took 38.101747ms
    benchmarks_test.go:62: PIPES took 36.494892ms
ok      github.com/Nearhan/mindx        0.124s

~/Code/golang/src/github.com/Nearhan/mindx master
❯
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
