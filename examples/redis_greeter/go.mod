module redis_greeter

go 1.19

require (
	github.com/tdv/go-care v1.0.3
	google.golang.org/grpc v1.50.1
	proto v0.0.0
	rediscache v0.0.0
)

replace github.com/tdv/go-care v1.0.3 => ../..

replace proto v0.0.0 => ./proto

replace rediscache v0.0.0 => ./rediscache

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
