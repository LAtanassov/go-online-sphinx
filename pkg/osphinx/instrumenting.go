package osphinx

import (
	"math/big"
	"time"

	"github.com/go-kit/kit/metrics"
)

// InstrumentingService wraps a Service and intruments method calls.
type InstrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) *InstrumentingService {
	return &InstrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

// ExpK wraps service.ExpK and instruments it.
func (s *InstrumentingService) ExpK(r, q *big.Int) (b0 *big.Int, err error) {

	defer func(begin time.Time) {
		s.requestCount.With("method", "ExpK").Add(1)
		s.requestLatency.With("method", "ExpK").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.ExpK(r, q)
}
