package mindx

import (
	"os/exec"
	"testing"
	"time"
)

func runImpl(gen, proc string) time.Duration {

	ch1 := make(chan bool)
	ch2 := make(chan bool)

	stat := time.Now()

	cmd1 := exec.Command("./cmd/benchmarks/proto/processor/processor")
	cmd2 := exec.Command("./cmd/benchmarks/proto/generator/generator")

	go func(c chan bool) {
		cmd1.Run()
		c <- true

	}(ch1)

	go func(c chan bool) {
		cmd2.Run()
		c <- true

	}(ch2)

	<-ch2
	<-ch1

	return time.Since(stat)
}

func TestImpl(t *testing.T) {

	testCases := []struct {
		Name  string
		gen   string
		pross string
	}{
		{
			Name:  "GRPC",
			gen:   "./cmd/benchmarks/proto/generator/generator",
			pross: "./cmd/benchmarks/proto/processor/processor",
		}, {
			Name:  "UNIX SOCKETS",
			gen:   "./cmd/benchmarks/unix/generator/generator",
			pross: "./cmd/benchmarks/unix/processor/processor",
		},
		{
			Name:  "UNIX PIPES",
			gen:   "./cmd/benchmarks/pipes/generator/generator",
			pross: "./cmd/benchmarks/pipes/processor/processor",
		},
	}

	for _, tc := range testCases {
		d := runImpl(tc.gen, tc.pross)
		t.Logf("%s took %s", tc.Name, d)
	}

}
