package main

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/LAtanassov/go-online-sphinx/pkg/service"
	"github.com/pkg/errors"

	kitlog "github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/kelseyhightower/envconfig"
)

// Configuration represents a set of environment variables loaded at startup e.g.:
// #!/bin/sh
// export OSSRV_ADDR=8080
// export OSSRV_KEYLENGTH=1024
// export OSSRV_HASH=sha256
// export OSSRV_IDHEX=1A3F1
// export OSSRV_KHEX=AFFEE
// export OSSRV_QHEX=BEEFF
type Configuration struct {
	Addr      string `default:":8080"`
	KeyLength int    `default:"1024"`
	Hash      string `default:"sha256"`
	IDHex     string
	KHex      string
	Q0Hex     string
}

func main() {

	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	logger.Log("service", "starting")

	// === configuration ===

	var c Configuration
	err := envconfig.Process("ossrv", &c)
	if err != nil {
		logger.Log("err", fmt.Sprintf("%+v", errors.Wrap(err, "failed to process environment variables")))
		os.Exit(1)
	}

	httpAddr := flag.String("ossrv.addr", c.Addr, "http listen address")
	keyLength := flag.Int("ossrv.key.length", c.KeyLength, "key length")
	hashName := flag.String("ossrv.hash", c.Hash, "hash function")
	idhex := flag.String("ossrv.id.hex", c.IDHex, "secret k in hex")
	khex := flag.String("ossrv.k.hex", c.KHex, "secret k in hex")
	q0hex := flag.String("ossrv.q0.hex", c.Q0Hex, "secret Q0 in hex")
	flag.Parse()

	hashFn := getHashBy(*hashName)
	id := getOrGenerate(*idhex, *keyLength)
	k := getOrGenerate(*khex, *keyLength)
	q0 := getOrGenerate(*q0hex, *keyLength)

	logger.Log("service", "starting", "state", "configured")

	// === service layer ===

	users := service.NewUserRepository()

	fieldKeys := []string{"method"}
	cfg := service.NewConfiguration(id, k, q0, big.NewInt(int64(*keyLength)), hashFn)

	var svc service.Service
	svc = service.New(users, cfg)
	svc = service.NewLoggingMiddleware(kitlog.With(logger, "component", "online_sphinx"))(svc)
	svc = service.NewInstrumentingMiddleware(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "online_sphinx",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "online_sphinx",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys))(svc)

	// === transport layer ===

	t := service.NewHTTPTransport(svc, kitlog.With(logger, "component", "transport"))

	mux := http.NewServeMux()

	mux.Handle("/v1/register", t.MakeRegisterHandler())
	mux.Handle("/v1/login/expk", t.MakeExpKHandler())
	mux.Handle("/v1/login/challenge", t.MakeChallengeHandler())

	mux.Handle("/v1/metadata", t.MakeMetadataHandler())
	mux.Handle("/v1/add", t.MakeAddHandler())
	mux.Handle("/v1/get", t.MakeGetHandler())

	handler := http.NewServeMux()
	handler.Handle("/", service.MakeAccessControl(mux))
	handler.Handle("/metrics", promhttp.Handler())
	handler.Handle("/_status/liveness", t.MakeLivenessHandler())
	handler.Handle("/_status/readiness", t.MakeReadinessHandler())

	// === startup ===

	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	server := &http.Server{
		Addr:         *httpAddr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func(server *http.Server) {
		logger.Log("service", "started", "listening", httpAddr)
		err := server.ListenAndServe()
		if err != nil {
			logger.Log("err", fmt.Sprintf("%+v", errors.Wrap(err, "failed to start http service")))
			os.Exit(1)
		}
	}(server)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals

	// === shutdown ===

	logger.Log("service", "shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log("err", fmt.Sprintf("%+v", errors.Wrap(err, "failed to shutdown http service")))
		os.Exit(1)
	}

	logger.Log("service", "stopped")
}

// === utils ===

func getHashBy(name string) func() hash.Hash {
	switch name {
	case "sha256":
		return sha256.New
	case "sha512":
		return sha512.New
	default:
		return sha256.New
	}
}

func getOrGenerate(hex string, keyLength int) *big.Int {
	return big.NewInt(1)
}
