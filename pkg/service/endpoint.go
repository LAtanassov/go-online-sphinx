package service

import (
	"context"
	"math/big"
	"os"

	"github.com/go-kit/kit/endpoint"
)

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

func (r expKResponse) error() error { return r.Err }

func makeExpKEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(expKRequest)

		ski, sID, sNonce, bd, q0, kv, err := s.ExpK(req.cID, req.cNonce, req.b, req.q)

		// TODO: should store within session
		os.Setenv("SKi", ski.Text(16))

		return expKResponse{sID: sID, sNonce: sNonce, bd: bd, q0: q0, kv: kv, Err: err}, nil
	}
}

type registerRequest struct {
	cID *big.Int
}

type registerResponse struct {
	Err error `json:"error,omitempty"`
}

func makeRegisterEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(registerRequest)
		err := s.Register(req.cID)
		return registerResponse{Err: err}, nil
	}
}

type verifyRequest struct {
	g *big.Int
	q *big.Int
}

type verifyResponse struct {
	r   *big.Int
	Err error `json:"error,omitempty"`
}

type metadataRequest struct {
	cID *big.Int
	mac []byte
}

type metadataResponse struct {
	domains []Domain
	Err     error `json:"error,omitempty"`
}

func makeVerifyEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(verifyRequest)

		// TODO: should store within session
		ski := new(big.Int)
		ski.SetString(os.Getenv("SKi"), 16)

		r, err := s.Verify(ski, req.g, req.q)
		return verifyResponse{r: r, Err: err}, nil
	}
}

func makeMetadataEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(metadataRequest)
		domains, err := s.GetMetadata(req.cID, req.mac)
		return metadataResponse{domains: domains, Err: err}, nil
	}
}
