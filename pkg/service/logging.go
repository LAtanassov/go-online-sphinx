package service

import (
	"math/big"
	"time"

	"github.com/go-kit/kit/log"
)

// NewLoggingMiddleware returns a new instance of a logging middleware.
func NewLoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingService{logger, next}
	}
}

type loggingService struct {
	logger log.Logger
	Service
}

func (s *loggingService) Register(cID *big.Int) (err error) {

	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Register",
			"cID", cID.Text(16),

			"err", err,

			"took", time.Since(begin),
		)
	}(time.Now())

	return s.Service.Register(cID)
}

func (s *loggingService) ExpK(cID, cNonce, b, q *big.Int) (ski, sID, sNonce, bd, q0, kv *big.Int, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "ExpK",
			"cID", cID.Text(16),
			"cNonce", cNonce.Text(16),
			"b", b.Text(16),
			"q", q.Text(16),

			"ski", ski.Text(16),
			"sID", sID.Text(16),
			"sNonce", sNonce.Text(16),
			"bd", bd.Text(16),
			"q0", q0.Text(16),
			"kv", kv.Text(16),
			"err", err,

			"took", time.Since(begin),
		)
	}(time.Now())

	return s.Service.ExpK(cID, cNonce, b, q)
}
func (s *loggingService) Challenge(ski, g, q *big.Int) (r *big.Int, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Challenge",
			"ski", ski.Text(16),
			"g", g.Text(16),
			"q", q.Text(16),

			"r", r.Text(16),
			"err", err,

			"took", time.Since(begin),
		)
	}(time.Now())

	return s.Service.Challenge(ski, g, q)
}

func (s *loggingService) VerifyMAC(mac []byte, cID *big.Int, data ...[]byte) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "VerifyMAC",
			"mac", mac,
			"cID", cID.Text(16),
			"data", data,

			"err", err,

			"took", time.Since(begin),
		)
	}(time.Now())

	return s.Service.VerifyMAC(mac, cID, data...)
}

func (s *loggingService) GetMetadata(cID *big.Int) (domains []string, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "GetMetadata",
			"cID", cID.Text(16),

			"domains", domains,
			"err", err,

			"took", time.Since(begin),
		)
	}(time.Now())

	return s.Service.GetMetadata(cID)
}
func (s *loggingService) Add(cID *big.Int, domain string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Add",
			"cID", cID.Text(16),
			"domain", domain,

			"err", err,

			"took", time.Since(begin),
		)
	}(time.Now())

	return s.Service.Add(cID, domain)
}
func (s *loggingService) Get(cID *big.Int, domain string, bmk *big.Int, q *big.Int) (bj, qj *big.Int, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Get",
			"cID", cID.Text(16),
			"domain", domain,
			"bmk", bmk.Text(16),
			"q", q.Text(16),

			"bj", bj.Text(16),
			"qj", qj.Text(16),
			"err", err,

			"took", time.Since(begin),
		)
	}(time.Now())

	return s.Service.Get(cID, domain, bmk, q)
}
