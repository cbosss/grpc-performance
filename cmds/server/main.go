package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"google.golang.org/grpc/reflection"

	"github.com/cbosss/grpc-performance/proto"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedEchoerServer
}

func (s *server) Echo(ctx context.Context, req *proto.EchoRequest) (*proto.EchoResponse, error) {
	return &proto.EchoResponse{Msg: req.Msg}, nil
}

func main() {
	addr := flag.String("addr", ":5555", "server address")
	paddr := flag.String("paddr", "", "pprof address")
	buff := flag.Int("buff", 32, "read / write buffer sizes (KB)")
	streams := flag.Int("streams", 0, "max amount of stream")
	mu := flag.Bool("mu", false, "enable mutex profiling")

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

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.MaxConcurrentStreams(uint32(*streams)),
		grpc.ReadBufferSize(*buff*1024),
		grpc.WriteBufferSize(*buff*1024),
	)

	ctx, _ := setupSignalContext()
	go func() {
		<-ctx.Done()
		s.GracefulStop()
	}()

	proto.RegisterEchoerServer(s, &server{})
	reflection.Register(s)

	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func setupSignalContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Printf("Received signal '%s', shutting down\n", sig)
		defer cancel()
	}()

	return ctx, cancel
}
