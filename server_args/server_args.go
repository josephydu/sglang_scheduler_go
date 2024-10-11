package server_args

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

type ServerArgs struct {
	Host              string
	Port              int
	LogLevel          string
	LogLevelHttp      string
	LogRequests       bool
	ShowTimeCost      bool
	RandomSeed        int
	LoadBalanceMethod string
}

func ParseArgs() *ServerArgs {
	host := flag.String("host", "127.0.0.1", "The host of the server ")
	port := flag.Int("port", 23000, "The port of the server")
	randomSeed := flag.Int("random-seed", rand.Intn(1000000000), "The random seed")
	logLevel := flag.String("log-level", "info", "The logging level of all loggers.")
	logLevelHttp := flag.String("log-level-http", "", "The logging level of HTTP server. If not set, reuse --log-level by default.")
	logRequests := flag.Bool("log-requests", false, "Log the inputs and outputs of all requests.")
	showTimeCost := flag.Bool("show-time-cost", false, "Show time cost of custom marks.")
	loadBalanceMethod := flag.String("load-balance-method", "round_robin", "The load balancing strategy for data parallelism.")

	flag.Parse()

	return &ServerArgs{
		Host:              *host,
		Port:              *port,
		RandomSeed:        *randomSeed,
		LogLevel:          *logLevel,
		LogLevelHttp:      *logLevelHttp,
		LogRequests:       *logRequests,
		ShowTimeCost:      *showTimeCost,
		LoadBalanceMethod: *loadBalanceMethod,
	}
}

func (s *ServerArgs) URL() string {
	return fmt.Sprintf("http://%s:%d", s.Host, s.Port)
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
