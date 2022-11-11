// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/tdv/go-care"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	api "proto/api/pb"
	"time"
)

func main() {
	var (
		port        = flag.Int("port", 55555, "The server port.")
		host        = flag.String("host", "localhost", "The server host.")
		memoization = flag.Bool("memoization", false, "Use response memoization on the client side.")
		name        = flag.String("name", "Client", "The name for greeting.")
		repeat      = flag.Uint("repeat", 1, "Number of the request repetitions.")
		help        = flag.Bool("help", false, "Print usage instructions and exit.")
	)

	flag.Parse()

	if help != nil && *help {
		flag.Usage()
		os.Exit(0)
	}

	grpcopts := make([]grpc.DialOption, 0, 1)

	grpcopts = append(grpcopts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if memoization != nil && *memoization {
		opts := care.NewOptions()
		opts.Methods.Add("/api.GreeterService/SayHello", time.Second*60)

		unary := care.NewClientUnaryInterceptor(opts)
		grpcopts = append(grpcopts, unary)
	}
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", *host, *port),
		grpcopts...,
	)

	if err != nil {
		log.Fatalf("Failed to dial the server. Error: %v\n", err)
	}

	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	client := api.NewGreeterServiceClient(conn)

	ctx := context.Background()

	var count uint = 1
	if repeat != nil {
		count = *repeat
	}

	req := api.SayHelloRequest{
		Name: *name,
	}

	for ; count > 0; count-- {
		resp, err := client.SayHello(ctx, &req)
		if err != nil {
			log.Fatalf("Failed to call 'SayHello' Error %v\n", err)
		}

		log.Printf("Response from the server is '%s'\n", resp.Greeting)
	}
}
