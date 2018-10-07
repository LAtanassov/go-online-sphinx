package osphinx

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

// ExpK wraps service.ExpK and writes log msg
func (s *LoggingService) ExpK(r, q *big.Int) (b0 *big.Int, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "ExpK",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.ExpK(r, q)
}
