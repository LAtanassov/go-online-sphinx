package server

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

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
		encodeResponse,
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
		UID string `json:"uid"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return registerRequest{
		uID: body.UID,
	}, nil
}

func decodeExpKRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		UID string `json:"uid"`
		R   string `json:"r"`
		Q   string `json:"q"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	rHex, err := hex.DecodeString(body.R)
	if err != nil {
		return nil, err
	}

	qHex, err := hex.DecodeString(body.Q)
	if err != nil {
		return nil, err
	}

	return expKRequest{
		uID: body.UID,
		r:   big.NewInt(0).SetBytes(rHex),
		q:   big.NewInt(0).SetBytes(qHex),
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
