# Greeter
The example demonstrates features quite close to the canonical 'Hello World' example.  

Using various flag-set you can try various cases like with/without memoization on the client- or server-side or at the both.    

In order to make an impression close to long time response computation having added the fake delay (time.Sleep) on the server-side. Having been playing with flags you'll notice how everything works, compare and evaluate the results.  

# Build
```bash
git clone https://github.com/tdv/go-care
cd go-care/examples/greeter/
make
```

# Run
## Getting info about usage
```bash
cd out
./server -help

Usage of ./server:
  -help
    	Print usage instructions and exit.
  -memoization
    	Use response memoization. (default true)
  -port int
    	The server port. (default 55555)
  -with-reflection
    	Enable reflection for the service.
    	
    	
./client -help

Usage of ./client:
  -help
    	Print usage instructions and exit.
  -host string
    	The server host. (default "localhost")
  -memoization
    	Use response memoization on the client side.
  -name string
    	The name for greeting. (default "Client")
  -port int
    	The server port. (default 55555)
  -repeat uint
    	Number of the request repetitions. (default 1)
```
Now you're familiar with the options. Let's go to try  

## Let's try out
### With the server-side memoization
Run the server with the memoization     
```bash
./server -memoization=1
Started on the port 55555. Press Ctrl+C to quit.
2022/10/25 22:44:30 Handling the request. 'SayHello' has been called for 'User'
```
Run the client without the memoization and with the request repetition for three times.  
```bash
./client -name="User" -repeat=3
2022/10/25 22:44:32 Response from the server is 'Hello: User'
2022/10/25 22:44:32 Response from the server is 'Hello: User'
2022/10/25 22:44:32 Response from the server is 'Hello: User'
```
Look at the output.
- On the server-side the real handler has been called once in spite of on the client-side three requests have been sent.
- All requests reached the server, but only one handled really, the others obtained from the cache and the handler didn't call.
- Notice that on the client-side the time difference among the responses is not significant. But you could notice that the very first request had been processing longer than the others processed. We had added fake long computation within a handler via the delay (time.Sleep) and that worked for the first request.

### With the client-side memoization
Run the server without the memoization
```bash
./server -memoization=0
Started on the port 55555. Press Ctrl+C to quit.
2022/10/25 23:04:18 Handling the request. 'SayHello' has been called for 'User'
2022/10/25 23:04:20 Handling the request. 'SayHello' has been called for 'User'
2022/10/25 23:04:22 Handling the request. 'SayHello' has been called for 'User'
```
Run the client without the memoization too and with the request repetition for three times.
```bash
./client -name="User" -repeat=3
2022/10/25 23:04:20 Response from the server is 'Hello: User'
2022/10/25 23:04:22 Response from the server is 'Hello: User'
2022/10/25 23:04:24 Response from the server is 'Hello: User'
```
Look at the results. There is a significant (2 second) delay at the both sides.  

Let's try out memoization only on the client-side.

Run the server without the memoization again
```bash
./server -memoization=0
Started on the port 55555. Press Ctrl+C to quit.
2022/10/25 23:11:04 Handling the request. 'SayHello' has been called for 'User'
```
Run the client with the memoization
```bash
./client -memoization=1 -name="User" -repeat=3
2022/10/25 23:11:06 Response from the server is 'Hello: User'
2022/10/25 23:11:06 Response from the server is 'Hello: User'
2022/10/25 23:11:06 Response from the server is 'Hello: User'
```
In this case only the first request had reached the server, the others remained on the client-side and obtained from the client's cache.  
There is no delay amongs the responses, only the first had been waited too long (in ratio with others).

### At the both sides
```bash
./server -memoization=1
Started on the port 55555. Press Ctrl+C to quit.
2022/10/25 23:17:27 Handling the request. 'SayHello' has been called for 'User'


./client -memoization=1 -name="User" -repeat=3
2022/10/25 23:17:29 Response from the server is 'Hello: User'
2022/10/25 23:17:29 Response from the server is 'Hello: User'
2022/10/25 23:17:29 Response from the server is 'Hello: User'
```
In this case  nothing interesting. Everything is predicted.

# Note
All cases have been considered with similar params. If we use different params, there will be the similar results like above, but different for each param-set because a key of the response is the hash of the request (of all params)  

## In additional
The server's key '-with-reflection' might be used for discovering all services and methods via grpc_cli at least like below
```bash
…out$grpc_cli ls localhost:55555
api.GreeterService
grpc.reflection.v1alpha.ServerReflection

…out$grpc_cli ls localhost:55555 api.GreeterService
SayHello
```

*Good luck discovering the features of the go-care package.*