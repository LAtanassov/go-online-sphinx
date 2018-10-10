package server

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

// Register wraps service.Register and instruments it.
func (s *InstrumentingService) Register(id string) (err error) {

	defer func(begin time.Time) {
		s.requestCount.With("method", "Register").Add(1)
		s.requestLatency.With("method", "Register").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Register(id)
}

// ExpK wraps service.ExpK and instruments it.
func (s *InstrumentingService) ExpK(uID string, r, q *big.Int) (sID string, sNonce, bd, q0, kv *big.Int, err error) {

	defer func(begin time.Time) {
		s.requestCount.With("method", "ExpK").Add(1)
		s.requestLatency.With("method", "ExpK").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.ExpK(uID, r, q)
}
