// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// The license can be found in the LICENSE file.
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
	"rediscache"
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
		redishost   = flag.String("redishost", "localhost", "The Redis server host.")
		redisport   = flag.Int("redisport", 6379, "The Redis server port.")
		redisdb     = flag.Int("redisdb", 0, "The Redis DB id.")
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

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if memoization != nil && *memoization {
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

	defer conn.Close()

	client := api.NewGreeterServiceClient(conn)

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
