package http_server

import (
	"context"
	"net/http"
	"runtime"
	"syscall"

	"github.com/danielgtaylor/huma/v2"
	"github.com/kgantsov/limiter/pkg/limiter"
)

type (
	Handler struct {
		RateLimiter *limiter.RateLimiter
		PathMap     map[string]string
	}
)

func (h *Handler) Reduce(ctx context.Context, input *ReduceRequest) (*ReduceResponse, error) {
	key := input.Key

	val, err := h.RateLimiter.Reduce(
		key, input.MaxTokens, input.RefillTime, input.RefillAmount, input.Tokens,
	)

	if err != nil {
		return nil, huma.Error409Conflict("Failed to reduce tokens", err)
	}

	res := &ReduceResponse{
		Status: http.StatusOK,
		Body: ReduceResponseBody{
			Status: "OK",
			Key:    input.Key,
			Tokens: val,
		},
	}
	return res, nil
}

func (h *Handler) Remove(ctx context.Context, input *RemoveRequest) (*RemoveResponse, error) {
	key := input.Key

	h.RateLimiter.Remove(key)

	res := &RemoveResponse{Status: http.StatusOK}
	res.Body.Status = "OK"
	return res, nil
}

func (h *Handler) Stats(ctx context.Context, input *StatsRequest) (*StatsResponse, error) {

	rusage := new(syscall.Rusage)
	syscall.Getrusage(0, rusage)
	userCPU := rusage.Utime.Sec*1e9 + int64(rusage.Utime.Usec)
	maxRSS := float64(rusage.Maxrss)
	numberOfKeys := h.RateLimiter.Len()

	res := &StatsResponse{
		Status: http.StatusOK,
		Body: StatsResponseBody{
			NumGoroutines: runtime.NumGoroutine(),
			CPU:           userCPU,
			MaxRSS:        maxRSS,
			NumberOfKeys:  numberOfKeys,
		},
	}

	return res, nil
}
