# limiter

![Build status](http://thekoss.ml:8000/api/badges/kgantsov/limiter/status.svg) 

Limiter is a rate limiter application that has http and redis interfaces


### Strart limiter service taht listens HTTP requests

    ./limiter --http_port=3000 --redis_port=46379

### Start with enabled prometheus

    ./limiter --http_port=3000 --redis_port=46379 --prometheus true


### Usage

#### Reduce limiter that have 5 max tokens that refils every 10 seconds by 5 tokens (reduces 1 token per request)
    curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://127.0.0.1:3000/API/v1/limiter/test/5/10/5/1/

#### typical respose

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=utf-8
    Date: Tue, 02 Oct 2018 19:11:08 GMT
    Content-Length: 91

    {"key":"test","max_tokens":5,"refill_amount":5,"refill_time":10,"tokens":1,"tokens_left":4}


### Benchmarks

    go test -cpu 1,2,4,8 -run XXX -bench . -benchmem
    goos: darwin
    goarch: amd64
    pkg: github.com/kgantsov/limiter/pkg/limiter
    BenchmarkReduce_100_1000_1_1000_1            	 2480442	       455 ns/op	      29 B/op	       3 allocs/op
    BenchmarkReduce_100_1000_1_1000_1-2          	 4903908	       265 ns/op	      29 B/op	       3 allocs/op
    BenchmarkReduce_100_1000_1_1000_1-4          	 8643420	       150 ns/op	      29 B/op	       3 allocs/op
    BenchmarkReduce_100_1000_1_1000_1-8          	11123530	        95.4 ns/op	      29 B/op	       3 allocs/op
    BenchmarkReduce_10000000_1000_10_1000_10     	 1210626	       917 ns/op	     193 B/op	       4 allocs/op
    BenchmarkReduce_10000000_1000_10_1000_10-2   	 1939578	       559 ns/op	     227 B/op	       4 allocs/op
    BenchmarkReduce_10000000_1000_10_1000_10-4   	 4674661	       297 ns/op	     196 B/op	       4 allocs/op
    BenchmarkReduce_10000000_1000_10_1000_10-8   	 6048070	       207 ns/op	     163 B/op	       4 allocs/op
    PASS
    ok  	github.com/kgantsov/limiter/pkg/limiter	14.993s
