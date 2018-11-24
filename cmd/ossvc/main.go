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
// export OSSVC_ADDR=443
// export OSSVC_KEYPATH=server.key
// export OSSVC_CERTPATH=server.crt
// export OSSVC_TIMEOUTSEC=15
// export OSSVC_KEYLENGTH=1024
// export OSSVC_HASH=sha256
// export OSSVC_IDHEX=1A3F1
// export OSSVC_KHEX=AFFEE
// export OSSVC_QHEX=BEEFF
// export OSSVC_KHEX=AFFEE
type Configuration struct {
	Addr     string `default:":443"`
	KeyPath  string `default:"server.key"`
	CertPath string `default:"server.crt"`

	TimeoutSec int    `default:"15"`
	KeyLength  int    `default:"1024"`
	Hash       string `default:"sha256"`
	IDHex      string
	KHex       string
	Q0Hex      string
}

func main() {

	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	// === configuration ===

	var c Configuration
	err := envconfig.Process("ossvc", &c)
	if err != nil {
		logger.Log("err", fmt.Sprintf("%+v", errors.Wrap(err, "failed to process environment variables")))
		os.Exit(1)
	}

	httpAddr := flag.String("ossvc.addr", c.Addr, "http listen address")
	keyPath := flag.String("ossvc.key.path", c.KeyPath, "server key path")
	certPath := flag.String("ossvc.cert.path", c.CertPath, "server cert path")
	timeoutSec := flag.Int("ossvc.timeout.sec", c.TimeoutSec, "http timeout in seconds")
	keyLength := flag.Int("ossvc.key.length", c.KeyLength, "key length")
	hashName := flag.String("ossvc.hash", c.Hash, "hash function")
	idhex := flag.String("ossvc.id.hex", c.IDHex, "secret k in hex")
	khex := flag.String("ossvc.k.hex", c.KHex, "secret k in hex")
	q0hex := flag.String("ossvc.q0.hex", c.Q0Hex, "secret Q0 in hex")
	flag.Parse()

	hashFn := getHashBy(*hashName)
	id := getOrGenerate(*idhex, *keyLength)
	k := getOrGenerate(*khex, *keyLength)
	q0 := getOrGenerate(*q0hex, *keyLength)

	logger.Log("service", "starting", "state", "configured")

	// === service layer ===

	logger.Log("service", "starting")
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
	mux.Handle("/v1/logout", t.MakeLogoutHandler())

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
		ReadTimeout:  time.Duration(*timeoutSec) * time.Second,
		WriteTimeout: time.Duration(*timeoutSec) * time.Second,
	}
	go func(server *http.Server) {
		logger.Log("service", "started", "listening", httpAddr)
		err := server.ListenAndServeTLS(*certPath, *keyPath)
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
