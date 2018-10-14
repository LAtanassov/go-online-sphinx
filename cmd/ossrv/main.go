package main

import (
	"crypto/rand"
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

	var repo service.Repository
	repo = service.NewInMemoryRepository()

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

	var svc service.Service
	svc = service.New("sID", k, q0, bits, repo)
	svc = service.NewLoggingService(logger, svc)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/v1/register", service.MakeRegisterHandler(svc, httpLogger))
	mux.Handle("/v1/login/expk", service.MakeExpKHandler(svc, httpLogger))
	mux.Handle("/v1/login/verify", service.MakeVerifyHandler(svc, httpLogger))

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
