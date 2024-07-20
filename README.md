# limiter

![Build status](http://thekoss.ml:8000/api/badges/kgantsov/limiter/status.svg) 

Limiter is a rate limiter application that has http and redis interfaces


### Strart limiter service taht listens HTTP requests
```bash
./limiter --http_port=3000 --redis_port=46379
```

### Start with enabled prometheus
```bash
./limiter --http_port=3000 --redis_port=46379 --prometheus true
```

### Usage

#### Reduce limiter that have 5 max tokens that refils every 10 seconds by 5 tokens (reduces 1 token per request)
```bash
curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://127.0.0.1:3000/API/v1/limiter/test/5/10/5/1/
```

#### typical respose
```bash
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 02 Oct 2018 19:11:08 GMT
Content-Length: 91

{"key":"test","max_tokens":5,"refill_amount":5,"refill_time":10,"tokens":1,"tokens_left":4}
```

### Benchmarks

```bash
go test -cpu 1,2,4,8 -run XXX -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/kgantsov/limiter/pkg/limiter
BenchmarkReduce_100_1000_1_1000_1            	 8154126	       150.8 ns/op	      13 B/op	       1 allocs/op
BenchmarkReduce_100_1000_1_1000_1-2          	12831616	        82.32 ns/op	      13 B/op	       1 allocs/op
BenchmarkReduce_100_1000_1_1000_1-4          	27081812	        46.19 ns/op	      13 B/op	       1 allocs/op
BenchmarkReduce_100_1000_1_1000_1-8          	26452123	        46.04 ns/op	      13 B/op	       1 allocs/op
BenchmarkReduce_10000000_1000_10_1000_10     	 2576054	       461.0 ns/op	     173 B/op	       2 allocs/op
BenchmarkReduce_10000000_1000_10_1000_10-2   	 4805696	       275.5 ns/op	     183 B/op	       2 allocs/op
BenchmarkReduce_10000000_1000_10_1000_10-4   	 8261760	       138.7 ns/op	     114 B/op	       2 allocs/op
BenchmarkReduce_10000000_1000_10_1000_10-8   	 8805577	       128.7 ns/op	     109 B/op	       2 allocs/op
PASS
ok  	github.com/kgantsov/limiter/pkg/limiter	15.079s
```
