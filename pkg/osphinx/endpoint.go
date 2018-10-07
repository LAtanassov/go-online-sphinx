package osphinx

import (
	"context"
	"math/big"

	"github.com/go-kit/kit/endpoint"
)

type expKRequest struct {
	R *big.Int
	Q *big.Int
}

type expKResponse struct {
	B0  *big.Int
	Err error `json:"error,omitempty"`
}

func (r expKResponse) error() error { return r.Err }

func makeExpKEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(expKRequest)
		b0, err := s.ExpK(req.R, req.Q)
		return expKResponse{B0: b0, Err: err}, nil
	}
}
