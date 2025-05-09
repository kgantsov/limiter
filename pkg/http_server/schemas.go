package http_server

type ReduceRequest struct {
	Key          string `path:"key" maxLength:"1024" example:"user:1" doc:"Key for the lock"`
	MaxTokens    int64  `path:"max_tokens" minimum:"1" example:"1000" doc:"Maximum number of tokens in the bucket"`
	RefillTime   int64  `path:"refill_time" minimum:"1" example:"60" doc:"Time in seconds to refill the bucket"`
	RefillAmount int64  `path:"refill_amount" minimum:"1" example:"100" doc:"Number of tokens to refill the bucket"`
	Tokens       int64  `path:"tokens" minimum:"1" example:"1" doc:"Number of tokens to reduce"`
}

type ReduceResponseBody struct {
	Status string `json:"status" example:"OK" doc:"Status of the reduce operation"`
	Key    string `json:"key" example:"user:1" doc:"Key for the lock"`
	Tokens int64  `json:"tokens" example:"1" doc:"Number of tokens left in the bucket"`
}

type ReduceResponse struct {
	Status int
	Body   ReduceResponseBody
}

type RemoveRequest struct {
	Key          string `path:"key" maxLength:"1024" example:"user:1" doc:"Key for the lock"`
	MaxTokens    int64  `path:"max_tokens" minimum:"1" example:"1000" doc:"Maximum number of tokens in the bucket"`
	RefillTime   int64  `path:"refill_time" minimum:"1" example:"60" doc:"Time in seconds to refill the bucket"`
	RefillAmount int64  `path:"refill_amount" minimum:"1" example:"100" doc:"Number of tokens to refill the bucket"`
}

type RemoveResponse struct {
	Status int
	Body   struct {
		Status string `json:"status" example:"OK" doc:"Status of the remove operation"`
	}
}

type StatsRequest struct {
}

type StatsResponseBody struct {
	NumGoroutines int     `json:"num_goroutines" example:"1" doc:"Number of goroutines"`
	CPU           int64   `json:"cpu" example:"53373" doc:"CPU usage in nanoseconds"`
	MaxRSS        float64 `json:"max_rss" example:"15777792" doc:"Max RSS in bytes"`
	NumberOfKeys  int64   `json:"number_of_keys" example:"543" doc:"Number of keys in the rate limiter"`
}

type StatsResponse struct {
	Status int
	Body   StatsResponseBody
}
