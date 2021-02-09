package http_server

import (
	"fmt"
	"runtime"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
)

func (app *App) ReduceLimiter(c *gin.Context) {
	key := c.Param("key")
	maxTokens, err := strconv.Atoi(c.Param("max_tokens"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("Param `max_tokens` must be integer"),
		})
		return
	}
	refillTime, err := strconv.Atoi(c.Param("refill_time"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("Param `refill_time` must be integer"),
		})
		return
	}
	refillAmount, err := strconv.Atoi(c.Param("refill_amount"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("Param `refill_amount` must be integer"),
		})
		return
	}
	tokens, err := strconv.Atoi(c.Param("tokens"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("Param `tokens` must be integer"),
		})
		return
	}

	val, err := app.RateLimiter.Reduce(
		key, int64(maxTokens), int64(refillTime), int64(refillAmount), int64(tokens),
	)

	if err == nil {
		c.JSON(200, gin.H{
			"key":           key,
			"tokens_left":   val,
			"max_tokens":    maxTokens,
			"refill_time":   refillTime,
			"refill_amount": refillAmount,
			"tokens":        tokens,
		})
	} else {
		c.JSON(400, gin.H{
			"error":         err,
			"key":           key,
			"tokens_left":   val,
			"max_tokens":    maxTokens,
			"refill_time":   refillTime,
			"refill_amount": refillAmount,
			"tokens":        tokens,
		})
	}

}

func (app *App) Stat(c *gin.Context) {
	rusage := new(syscall.Rusage)
	syscall.Getrusage(0, rusage)
	userCPU := rusage.Utime.Sec*1e9 + int64(rusage.Utime.Usec)
	maxRSS := float64(rusage.Maxrss)
	numberOfKeys := app.RateLimiter.Len()

	c.JSON(
		200,
		gin.H{
			"num_goroutines": runtime.NumGoroutine(),
			"cpu":            userCPU,
			"max_rss":        maxRSS,
			"number_of_keys": numberOfKeys,
		},
	)
}

func (app *App) Remove(c *gin.Context) {
	key := c.Param("key")

	app.RateLimiter.Remove(key)

	c.JSON(204, gin.H{})
}
