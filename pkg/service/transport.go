package service

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"

	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

// ErrUnexpectedType is returned after a type cast failed.
var ErrUnexpectedType = errors.New("unexpected type")

// MakeRegisterHandler returns a handler
func MakeRegisterHandler(s Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	registerHandler := kithttp.NewServer(
		makeRegisterEndpoint(s),
		decodeRegisterRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/v1/register", registerHandler).Methods("POST")

	return r
}

// MakeExpKHandler returns a handler for the handling service.
func MakeExpKHandler(s Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	expKHandler := kithttp.NewServer(
		makeExpKEndpoint(s),
		decodeExpKRequest,
		encodeExpKResponse,
		opts...,
	)

	r.Handle("/v1/login/expk", expKHandler).Methods("POST")

	return r
}

// MakeAccessControl sets Header for access control
func MakeAccessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

// MakeLivenessHandler returns liveness handler
func MakeLivenessHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
}

// MakeReadinessHandler return readiness handler
func MakeReadinessHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
}

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return registerRequest{
		username: body.Username,
	}, nil
}

func decodeExpKRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Username string `json:"username"`
		CNonce   string `json:"cNonce"`
		B        string `json:"b"`
		Q        string `json:"q"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	cNonce := new(big.Int)
	cNonce.SetString(body.CNonce, 16)

	b := new(big.Int)
	b.SetString(body.B, 16)

	q := new(big.Int)
	q.SetString(body.Q, 16)

	return expKRequest{
		username: body.Username,
		cNonce:   cNonce,
		b:        b,
		q:        q,
	}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeExpKResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	r, ok := response.(expKResponse)
	if !ok {
		encodeError(ctx, ErrUnexpectedType, w)
		return nil
	}

	rjson := struct {
		sID    string
		sNonce string
		bd     string
		q0     string
		kv     string
		Err    error `json:"error,omitempty"`
	}{
		r.sID,
		r.sNonce.Text(16),
		r.bd.Text(16),
		r.q0.Text(16),
		r.kv.Text(16),
		r.Err,
	}

	return json.NewEncoder(w).Encode(rjson)
}

type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
