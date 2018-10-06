package osphinx

import (
	"context"
	"math/big"

	"github.com/go-kit/kit/endpoint"
)

type loginRequest struct {
	R *big.Int
	Q *big.Int
}

type loginResponse struct {
	B0  *big.Int
	Err error `json:"error,omitempty"`
}

func (r loginResponse) error() error { return r.Err }

func makeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(loginRequest)
		b0, err := s.Login(req.R, req.Q)
		return loginResponse{B0: b0, Err: err}, nil
	}
}
