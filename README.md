[![build workflow](https://github.com/tdv/go-care/actions/workflows/main.yml/badge.svg?event=push)](https://github.com/tdv/go-care/actions)

# go-care
A library for gRPC response memoization in Go, aimed to improve common performance of a project having made real response computation lower cost.  

go-care is an abbreviation of the 'Cached Response in Go'.  

# Introduction
Having faced an issue of lack of the memoization in grpc this library has been written as a solution. Maybe there was a solution, but I had not been searching for the one enough well. Tricky question... Nonetheless, let's say a bit about the go-care package.    

The package aims to improve performance via caching response for respective request params. Having included additional information from incoming data there is a possibility to make the most unique key for the cached response.  

Essential gist of the approach is to make an interceptor wherein compute the key --a hash of the incoming data-- and memorize/restore the responses, having been making the lower computation cost at the expense of the decreased number of the real computation.      

The package can be used with built-in LRU cache or with an external cache like Redis, Memcached, or something else. In order to use an external cache there is only the need to implement a caching interface.

# Features
- server- and client-side response memoization
- flexibility and room space for customization
- loosely coupled architecture
- thread safety
- easy to use :)

**Note**
- The built-in implementation of the in-memory cache doesn't support eviction by the TTL. It has developed only for demo and small MVP. In production you need to use go-care with Redis, Memcached, or other cache. That might be done by implementing the 'Cache' interface and providing the one via 'Options'.
- Currently, if a request contains sequences --like slices, arrays, maps, etc-- the response caching might be sensitive for the items’ order inside the sequences and if the orders are not the same for the same requests, some requests will reach the handlers, omitting the cache. It will be fixed in a little while.  

# Usage
1. Add the package to the project
2. On the server-side it might be articulated like below in pseudocode
```go
package main

import (
  ...
  "github.com/tdv/go-care"
  ...
)

// Your server implementation
// type server struct {
//   api.UnimplementedYourServerServiceServer
// }

// ...

func main() {
  // Other your code
  ...

  // Creating the options
  opts := care.NewOptions()

  // Adding methods for the memoization. 
  // You need to define the methods pool 
  // for response memoization.
  opts.Methods.Add("/api.GreeterService/SayHello", time.Second*60)

  // Other customization might be done via the opts variable.
  // See the examples.

  // Creating the server-side interceptor.
  unary := care.NewServerUnaryInterceptor(opts)
  // Providing / applying the interceptor
  grpcsrv := grpc.NewServer(unary)
  // Other your code

  ...
}
```
3. On the client-side the similar way
 ```go
package main

import (
  ...
  "github.com/tdv/go-care"
  ...
)

func main() {
  ...

  opts := care.NewDefaultOptions()
  opts.Methods.Add("/api.GreeterService/SayHello", time.Second*60)

  unary := care.NewClientUnaryInterceptor(opts)

  conn, err := grpc.Dial(
    fmt.Sprintf("%s:%d", *host, *port),
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    unary,
  )

  if err != nil {
    log.Fatalf...
  }

  defer conn.Close()

	...
}
```
More details you'll find among the [examples](https://github.com/tdv/go-care/tree/main/examples/greeter)
 
# Examples
The [examples](https://github.com/tdv/go-care/tree/main/examples) demonstrate go-care package's features. Having run server and client with different param-set you can try the ones.
- ['Greeter'](https://github.com/tdv/go-care/tree/main/examples/greeter) is close to the canonical 'Hello World' example, demonstrating all basic features.    
- ['Redis Greeter'](https://github.com/tdv/go-care/tree/main/examples/redis_greeter) is the same, but with Redis like an external cache is.   

# Compiler and OS
The package had developed and tested in Go 1.19 within Ubuntu 20.04. Hopefully, many other OS and compiler versions will be appropriate quite well.
