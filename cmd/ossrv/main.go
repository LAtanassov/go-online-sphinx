package main

import (
	"crypto/rand"
	"crypto/sha256"
	"flag"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/LAtanassov/go-online-sphinx/pkg/service"

	"github.com/go-kit/kit/log"
)

func main() {

	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var users service.UserRepository
	users = service.NewUserRepository()

	bits := big.NewInt(8)
	max := new(big.Int)
	max.Exp(big.NewInt(2), bits, nil)

	k, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.Log(err)
		os.Exit(1)
	}

	q0, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.Log(err)
		os.Exit(1)
	}

	// TODO: load from env. or file later on
	cfg := service.NewConfiguration(big.NewInt(1), k, q0, bits, sha256.New)

	var svc service.Service
	svc = service.New(users, cfg)

	mux := http.NewServeMux()

	mux.Handle("/v1/register", service.MakeRegisterHandler(svc))
	mux.Handle("/v1/login/expk", service.MakeExpKHandler(svc))
	mux.Handle("/v1/login/challenge", service.MakeChallengeHandler(svc))

	mux.Handle("/v1/metadata", service.MakeMetadataHandler(svc))
	mux.Handle("/v1/add", service.MakeAddHandler(svc))
	mux.Handle("/v1/get", service.MakeGetHandler(svc))

	http.Handle("/", service.MakeAccessControl(mux))
	http.Handle("/_status/liveness", service.MakeLivenessHandler())
	http.Handle("/_status/readiness", service.MakeReadinessHandler())

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}
