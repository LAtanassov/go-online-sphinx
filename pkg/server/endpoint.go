package server

import (
	"context"
	"math/big"

	"github.com/go-kit/kit/endpoint"
)

type expKRequest struct {
	username string
	cNonce   *big.Int
	b        *big.Int
	q        *big.Int
}

type expKResponse struct {
	sID    string
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

		sID, sNonce, bd, q0, kv, err := s.ExpK(req.username, req.b, req.q)
		return expKResponse{sID: sID, sNonce: sNonce, bd: bd, q0: q0, kv: kv, Err: err}, nil
	}
}

type registerRequest struct {
	username string
}

type registerResponse struct {
	Err error `json:"error,omitempty"`
}

func makeRegisterEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(registerRequest)
		err := s.Register(req.username)
		return registerResponse{Err: err}, nil
	}
}
