package service

import (
	"context"
	"math/big"
	"os"

	"github.com/go-kit/kit/endpoint"
)

func makeExpKEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(expKRequest)

		ski, sID, sNonce, bd, q0, kv, err := s.ExpK(req.cID, req.cNonce, req.b, req.q)

		// TODO: should store within session
		os.Setenv("SKi", ski.Text(16))

		return expKResponse{sID: sID, sNonce: sNonce, bd: bd, q0: q0, kv: kv, Err: err}, nil
	}
}

type expKRequest struct {
	cID    *big.Int
	cNonce *big.Int
	b      *big.Int
	q      *big.Int
}

type expKResponse struct {
	sID    *big.Int
	sNonce *big.Int
	bd     *big.Int
	q0     *big.Int
	kv     *big.Int
	Err    error
}

func makeRegisterEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(registerRequest)
		err := s.Register(req.cID)
		return registerResponse{Err: err}, nil
	}
}

type registerRequest struct {
	cID *big.Int
}

type registerResponse struct {
	Err error `json:"error,omitempty"`
}

func makeChallengeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(challengeRequest)

		// TODO: should be stored within session
		ski := new(big.Int)
		ski.SetString(os.Getenv("SKi"), 16)

		// verify MAC of request
		r, err := s.Challenge(ski, req.g, req.q)
		return challengeResponse{r: r, Err: err}, nil
	}
}

type challengeRequest struct {
	g *big.Int
	q *big.Int
}

type challengeResponse struct {
	r   *big.Int
	Err error `json:"error,omitempty"`
}

func makeMetadataEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(metadataRequest)
		// verify MAC of request
		domains, err := s.GetMetadata()
		return metadataResponse{domains: domains, Err: err}, nil
	}
}

type metadataRequest struct {
	cID *big.Int
	mac []byte
}

type metadataResponse struct {
	domains []string
	Err     error `json:"error,omitempty"`
}

func makeAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addRequest)
		// verify MAC of request
		err := s.Add(req.domain)
		return addResponse{Err: err}, nil
	}
}

type addRequest struct {
	domain string
}

type addResponse struct {
	Err error `json:"error,omitempty"`
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRequest)
		// verify MAC of request
		bj, qj, err := s.Get(req.domain, req.bmk)
		return getResponse{
			bj:  bj,
			qj:  qj,
			Err: err}, nil
	}
}

type getRequest struct {
	domain string
	bmk    *big.Int
}

type getResponse struct {
	bj  *big.Int
	qj  *big.Int
	Err error
}

func (r expKResponse) error() error { return r.Err }
