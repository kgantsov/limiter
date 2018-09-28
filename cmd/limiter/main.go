package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"syscall"

	"github.com/kgantsov/limiter/pkg/limiter"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"
)

type App struct {
	RateLimiter *limiter.RateLimiter
}

type RateLimiterParams struct {
	MaxTokens    int64 `json:"max_tokens"`
	RefillTime   int64 `json:"refill_time"`
	RefillAmount int64 `json:"refill_amount"`
}

func DefineRoutes(app *App, r *gin.Engine) {
	r.GET("/stats/", app.Stat)
	v1 := r.Group("/API/v1")
	{
		v1.GET("/limiter/:key/:max_tokens/:refill_time/:refill_amount/:tokens/", app.ReduceLimiter)
	}
}

func ListenAndServe(app *App, port int, debug bool) {
	log.Infof("Strarting API service on a port: %d", port)

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	if debug {
		r.Use(gin.Logger())
	}

	DefineRoutes(app, r)

	r.Run(fmt.Sprintf(":%d", port))
}

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
	var currentRSS int64

	VmRSSRe := regexp.MustCompile(`VmRSS:\s+(\d+)\s+kB`)
	pid := os.Getpid()

	output, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err == nil {
		result := VmRSSRe.FindStringSubmatch(string(output))
		currentRSS, _ = strconv.ParseInt(result[1], 10, 64)
	}

	rusage := new(syscall.Rusage)
	syscall.Getrusage(0, rusage)
	userCPU := rusage.Utime.Sec*1e9 + int64(rusage.Utime.Usec)
	maxRSS := float64(rusage.Maxrss)

	c.JSON(
		200,
		gin.H{
			"num_goroutines": runtime.NumGoroutine(),
			"CPU":            userCPU,
			"current_rss":    currentRSS,
			"max_rss":        maxRSS,
		},
	)
}

func main() {
	port := flag.Int("port", 9000, "Port")

	flag.Parse()

	app := &App{
		RateLimiter: limiter.NewRateLimiter(),
	}

	ListenAndServe(app, *port, false)
}
