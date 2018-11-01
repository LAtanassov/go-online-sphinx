package service

import (
	"math/big"
	"time"

	"github.com/go-kit/kit/metrics"
)

// NewInstrumentingMiddleware returns an instance of the instrumented middleware.
func NewInstrumentingMiddleware(counter metrics.Counter, latency metrics.Histogram) Middleware {
	return func(next Service) Service {
		return &instrumentingService{
			requestCount:   counter,
			requestLatency: latency,
			Service:        next,
		}
	}
}

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func (s *instrumentingService) Register(cID *big.Int) (err error) {

	defer func(begin time.Time) {
		s.requestCount.With("method", "Register").Add(1)
		s.requestLatency.With("method", "Register").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Register(cID)
}

func (s *instrumentingService) ExpK(cID, cNonce, b, q *big.Int) (ski, sID, sNonce, bd, q0, kv *big.Int, err error) {

	defer func(begin time.Time) {
		s.requestCount.With("method", "ExpK").Add(1)
		s.requestLatency.With("method", "ExpK").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.ExpK(cID, cNonce, b, q)
}
func (s *instrumentingService) Challenge(ski, g, q *big.Int) (r *big.Int, err error) {

	defer func(begin time.Time) {
		s.requestCount.With("method", "Challenge").Add(1)
		s.requestLatency.With("method", "Challenge").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Challenge(ski, g, q)
}

func (s *instrumentingService) VerifyMAC(mac []byte, ski *big.Int, data ...[]byte) (err error) {

	defer func(begin time.Time) {
		s.requestCount.With("method", "VerifyMAC").Add(1)
		s.requestLatency.With("method", "VerifyMAC").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.VerifyMAC(mac, ski, data...)
}

func (s *instrumentingService) GetMetadata(cID *big.Int) (domains []string, err error) {

	defer func(begin time.Time) {
		s.requestCount.With("method", "GetMetadata").Add(1)
		s.requestLatency.With("method", "GetMetadata").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetMetadata(cID)
}
func (s *instrumentingService) Add(cID *big.Int, domain string) (err error) {

	defer func(begin time.Time) {
		s.requestCount.With("method", "Add").Add(1)
		s.requestLatency.With("method", "Add").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Add(cID, domain)
}
func (s *instrumentingService) Get(cID *big.Int, domain string, bmk *big.Int, q *big.Int) (bj, qj *big.Int, err error) {

	defer func(begin time.Time) {
		s.requestCount.With("method", "Get").Add(1)
		s.requestLatency.With("method", "Get").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Get(cID, domain, bmk, q)
}
