package service

import (
	"math/big"
	"time"

	"github.com/go-kit/kit/log"
)

// LoggingService wraps around a Service interface
type LoggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) *LoggingService {
	return &LoggingService{logger, s}
}

// Register wraps service.Register and writes log msg
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

// ExpK wraps service.ExpK and writes log msg
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

// Verify wraps service.Verify and writes log msg
func (s *LoggingService) Verify(ski, g, q *big.Int) (r *big.Int, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Verify",
			"took", time.Since(begin),

			"ski", ski.Text(16),
			"g", g.Text(16),

			"r", r.Text(16),
			"err", err,
		)
	}(time.Now())

	return s.Service.Verify(ski, g, q)
}
