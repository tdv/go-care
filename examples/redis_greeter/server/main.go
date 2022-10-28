package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/tdv/go-care"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"proto/api/pb"
	"rediscache"
	"time"
)

type server struct {
	api.UnimplementedGreeterServiceServer
}

func (s *server) SayHello(ctx context.Context, req *api.SayHelloRequest) (*api.SayHelloResponse, error) {
	log.Printf("Handling the request. 'SayHello' has been called for '%s'\n", req.Name)

	if len(req.Name) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Empty name.")
	}

	resp := api.SayHelloResponse{
		Greeting: "Hello: " + req.Name,
	}

	// Adding something like  the long time computation in onder to
	// try out the profit of the memoization usage,
	time.Sleep(time.Second * 2)

	return &resp, status.New(codes.OK, "Ok.").Err()
}

func main() {
	var (
		port           = flag.Int("port", 55555, "The server port.")
		memoization    = flag.Bool("memoization", true, "Use response memoization.")
		withReflection = flag.Bool("with-reflection", false, "Enable reflection for the service.")
		help           = flag.Bool("help", false, "Print usage instructions and exit.")
		redishost      = flag.String("redishost", "localhost", "The Redis server host.")
		redisport      = flag.Int("redisport", 6379, "The Redis server port.")
		redisdb        = flag.Int("redisdb", 0, "The Redis DB id.")
	)

	flag.Parse()

	if help != nil && *help {
		flag.Usage()
		os.Exit(0)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen the interface. Error: %v\n", err)
	}

	var grpcsrv *grpc.Server

	if memoization != nil && *memoization {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		cache, err := rediscache.New(
			ctx,
			time.Millisecond*500,
			*redishost,
			*redisport,
			*redisdb,
		)

		if err != nil {
			log.Fatalf("Failed to create Redis cache instance. Error: %v\n", err)
		}

		opts := care.NewOptions()
		opts.Methods.Add("/api.GreeterService/SayHello", time.Second*60)
		opts.Cache = cache
		unary := care.NewServerUnaryInterceptor(opts)
		grpcsrv = grpc.NewServer(unary)
	} else {
		grpcsrv = grpc.NewServer()
	}

	srv := server{}
	api.RegisterGreeterServiceServer(grpcsrv, &srv)

	if withReflection != nil && *withReflection {
		reflection.Register(grpcsrv)
	}

	quit := make(chan struct{})

	go func() {
		if err = grpcsrv.Serve(listener); err != nil {
			log.Printf("Failed to run the server. Error: %v\n", err)
		}
		quit <- struct{}{}
	}()

	fmt.Printf("Started on the port %d. Press Ctrl+C to quit.\n", *port)

	<-quit
}
