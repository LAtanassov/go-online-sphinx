package service

import (
	"math/big"
	"time"

	"github.com/go-kit/kit/log"
)

// LoggingService wraps around a service layer and logs
// method calls, latency and parameters and errrors.
//
// ATTENTION: some parameter are secrets and should not be logged.
// This is a prototype and should evaluate the protocol.
type LoggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) *LoggingService {
	return &LoggingService{logger, s}
}

// Register wraps service.Register and logs method calls.
func (s *LoggingService) Register(cID *big.Int) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Register",
			"took", time.Since(begin),
			"cID", cID,
			"err", err,
		)
	}(time.Now())

	return s.Service.Register(cID)
}

// ExpK wraps service.ExpK and logs method calls.
func (s *LoggingService) ExpK(cID, cNonce, b, q *big.Int) (ski, sID, sNonce, bd, q0, kv *big.Int, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "ExpK",
			"took", time.Since(begin),
			"cID", cID.Text(16),
			"b", b.Text(16),
			"q", q.Text(16),
			"ski", sID.Text(16),
			"sID", sID.Text(16),
			"sNonce", sNonce.Text(16),
			"bd", bd.Text(16),
			"q0", q0.Text(16),
			"kv", kv.Text(16),

			"err", err,
		)
	}(time.Now())

	return s.Service.ExpK(cID, cNonce, b, q)
}

// Challenge wraps service.Challenge and logs method calls.
func (s *LoggingService) Challenge(ski, g, q *big.Int) (r *big.Int, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Challenge",
			"took", time.Since(begin),
			"ski", ski.Text(16),
			"g", g.Text(16),
			"r", r.Text(16),
			"err", err,
		)
	}(time.Now())

	return s.Service.Challenge(ski, g, q)
}

// VerifyMAC wraps service.VerifyMAC and logs method calls.
func (s *LoggingService) VerifyMAC(mac []byte, key []byte, data ...[]byte) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "VerifyMAC",
			"took", time.Since(begin),
			"mac", mac,
			"key", key,
			"data", len(data),
			"err", err,
		)
	}(time.Now())

	return s.Service.VerifyMAC(mac, key, data...)
}

// GetMetadata wraps service.GetMetadata and logs method calls.
func (s *LoggingService) GetMetadata() (domains []string, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "GetMetadata",
			"took", time.Since(begin),
			"domains", len(domains),
			"err", err,
		)
	}(time.Now())

	return s.Service.GetMetadata()
}

// Add wraps service.Add and logs method calls.
func (s *LoggingService) Add(domain string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Add",
			"took", time.Since(begin),
			"domain", domain,
			"err", err,
		)
	}(time.Now())

	return s.Service.Add(domain)
}

// Get wraps service.Get and logs method calls.
func (s *LoggingService) Get(domain string, bmk *big.Int) (bj, qj *big.Int, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Get",
			"took", time.Since(begin),
			"domain", domain,
			"bmk", bmk.Text(16),
			"err", err,
		)
	}(time.Now())

	return s.Service.Get(domain, bmk)
}
