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

    > $ go test -run XXX -bench . -benchmem
    goos: darwin
    goarch: amd64
    pkg: github.com/kgantsov/limiter/pkg/limiter
    BenchmarkReduce_100_1000_1_1000_1-4        	 2000000	       745 ns/op	       0 B/op	       0 allocs/op
    BenchmarkReduce_10000000_1000_1_1000_1-4   	 1000000	      2130 ns/op	     139 B/op	       0 allocs/op
    PASS
    ok  	github.com/kgantsov/limiter/pkg/limiter	150.661s
