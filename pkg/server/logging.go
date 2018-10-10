package server

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
func (s *LoggingService) Register(id string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Register",
			"took", time.Since(begin),
			"id", id,
			"err", err,
		)
	}(time.Now())

	return s.Service.Register(id)
}

// ExpK wraps service.ExpK and writes log msg
func (s *LoggingService) ExpK(uID string, r, q *big.Int) (sID string, sNonce, bd, q0, kv *big.Int, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "ExpK",
			"uID", uID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.Service.ExpK(uID, r, q)
}
