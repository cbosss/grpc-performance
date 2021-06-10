package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"

	"github.com/jhump/protoreflect/desc"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	"github.com/cbosss/grpc-performance/proto"
	pb "google.golang.org/protobuf/proto"
)

func main() {
	addr := flag.String("addr", ":5555", "server address")
	paddr := flag.String("paddr", "", "pprof address")
	mu := flag.Bool("mu", false, "enable mutex profiling")
	rps := flag.Uint("rps", 5000, "requests per second")
	concurrency := flag.Uint("concurrency", 500, "number of goroutines reading/writing")
	connections := flag.Uint("connections", 1, "number of http2 connections to grpc server")

	flag.Parse()

	if *paddr != "" {
		go func() {
			if *mu {
				runtime.SetMutexProfileFraction(5)
			}
			log.Printf("pprof on %s\n", *paddr)
			log.Println(http.ListenAndServe(*paddr, nil))
		}()
	}

	log.Printf("Running...\n")
	log.Printf("rps=%d concurrency=%d connections=%d\n", *rps, *concurrency, *connections)

	b, _ := pb.Marshal(&proto.EchoRequest{Msg: "7a1af845-ad1e-4411-b07d-93b089ab16d8"})

	report, err := runner.Run(
		"echo.Echoer.Echo",
		*addr,
		runner.WithBinaryDataFunc(func(mtd *desc.MethodDescriptor, callData *runner.CallData) []byte {
			return b
		}),
		runner.WithInsecure(true),
		runner.WithRPS(*rps),
		runner.WithRunDuration(30*time.Second),
		runner.WithConcurrency(*concurrency),
		runner.WithConnections(*connections),
	)
	log.Printf("Running... done!\n")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	printer := printer.ReportPrinter{
		Out:    os.Stdout,
		Report: report,
	}

	printer.Print("summary")
}
