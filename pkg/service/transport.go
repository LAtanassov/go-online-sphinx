package service

import (
	"context"
	"encoding/hex"
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

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		CID string `json:"CID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	cID := new(big.Int)
	cID.SetString(body.CID, 16)

	return registerRequest{
		cID: cID,
	}, nil
}

// MakeExpKHandler returns a handler for the handling service.
func MakeExpKHandler(s Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(getOrCreateSession),
		kithttp.ServerAfter(setSession),
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

func decodeExpKRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		CID    string `json:"cID"`
		CNonce string `json:"cNonce"`
		B      string `json:"b"`
		Q      string `json:"q"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	cID := new(big.Int)
	cID.SetString(body.CID, 16)

	cNonce := new(big.Int)
	cNonce.SetString(body.CNonce, 16)

	b := new(big.Int)
	b.SetString(body.B, 16)

	q := new(big.Int)
	q.SetString(body.Q, 16)

	return expKRequest{
		cID:    cID,
		cNonce: cNonce,
		b:      b,
		q:      q,
	}, nil
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

	body := struct {
		SID    string `json:"sID"`
		SNonce string `json:"sNonce"`
		BD     string `json:"bd"`
		Q0     string `json:"q0"`
		KV     string `json:"kv"`
		Err    error  `json:"error,omitempty"`
	}{
		r.sID.Text(16),
		r.sNonce.Text(16),
		r.bd.Text(16),
		r.q0.Text(16),
		r.kv.Text(16),
		r.Err,
	}

	return json.NewEncoder(w).Encode(body)
}

// MakeChallengeHandler returns a handler for the handling service.
func MakeChallengeHandler(s Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(getOrCreateSession),
		kithttp.ServerAfter(setSession),
	}

	challengeHandler := kithttp.NewServer(
		makeChallengeEndpoint(s),
		decodeChallengeRequest,
		encodeChallengeResponse,
		opts...,
	)

	r.Handle("/v1/login/challenge", challengeHandler).Methods("POST")

	return r
}

func decodeChallengeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		G string `json:"g"`
		Q string `json:"q"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	g := new(big.Int)
	g.SetString(body.G, 16)

	q := new(big.Int)
	q.SetString(body.Q, 16)

	return challengeRequest{
		g: g,
		q: q,
	}, nil
}

func encodeChallengeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if err, ok := response.(errorer); ok && err.error() != nil {
		encodeError(ctx, err.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp, ok := response.(challengeResponse)
	if !ok {
		encodeError(ctx, ErrUnexpectedType, w)
		return nil
	}

	body := struct {
		R   string `json:"r"`
		Err error  `json:"error,omitempty"`
	}{
		resp.r.Text(16),
		resp.Err,
	}

	return json.NewEncoder(w).Encode(body)
}

// MakeMetadataHandler ...
func MakeMetadataHandler(s Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(getOrCreateSession),
		kithttp.ServerAfter(setSession),
	}

	metadataHandler := kithttp.NewServer(
		makeMetadataEndpoint(s),
		decodeMetadataRequest,
		encodeMetadataResponse,
		opts...,
	)

	r.Handle("/v1/metadata", metadataHandler).Methods("GET")

	return r
}

func decodeMetadataRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		CID string `json:"cID"`
		MAC string `json:"mac"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	cID := new(big.Int)
	cID.SetString(body.CID, 16)

	mac, err := hex.DecodeString(body.MAC)
	if err != nil {
		return nil, err
	}

	return metadataRequest{
		cID: cID,
		mac: mac,
	}, nil
}

func encodeMetadataResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if err, ok := response.(errorer); ok && err.error() != nil {
		encodeError(ctx, err.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp, ok := response.(metadataResponse)
	if !ok {
		encodeError(ctx, ErrUnexpectedType, w)
		return nil
	}

	body := struct {
		Domains []string `json:"domains"`
		Err     error    `json:"error,omitempty"`
	}{
		resp.domains,
		resp.Err,
	}

	return json.NewEncoder(w).Encode(body)
}

// MakeAddHandler ...
func MakeAddHandler(s Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(getOrCreateSession),
		kithttp.ServerAfter(setSession),
	}

	addHandler := kithttp.NewServer(
		makeAddEndpoint(s),
		decodeAddRequest,
		encodeAddResponse,
		opts...,
	)

	r.Handle("/v1/add", addHandler).Methods("POST")

	return r
}

func decodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Domain string `json:"domain"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return addRequest{
		domain: body.Domain,
	}, nil
}

func encodeAddResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if err, ok := response.(errorer); ok && err.error() != nil {
		encodeError(ctx, err.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp, ok := response.(addResponse)
	if !ok {
		encodeError(ctx, ErrUnexpectedType, w)
		return nil
	}

	body := struct {
		Err error `json:"error,omitempty"`
	}{
		resp.Err,
	}

	return json.NewEncoder(w).Encode(body)
}

// MakeGetHandler ...
func MakeGetHandler(s Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(getOrCreateSession),
		kithttp.ServerAfter(setSession),
	}

	addHandler := kithttp.NewServer(
		makeGetEndpoint(s),
		decodeGetRequest,
		encodeGetResponse,
		opts...,
	)

	r.Handle("/v1/get", addHandler).Methods("GET")

	return r
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Domain string `json:"domain"`
		BMK    string `json:"bmk"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	bmk := new(big.Int)
	_, ok := bmk.SetString(body.BMK, 16)
	if !ok {
		return getRequest{}, ErrUnexpectedType
	}

	return getRequest{
		domain: body.Domain,
		bmk:    bmk,
	}, nil
}

func encodeGetResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if err, ok := response.(errorer); ok && err.error() != nil {
		encodeError(ctx, err.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp, ok := response.(addResponse)
	if !ok {
		encodeError(ctx, ErrUnexpectedType, w)
		return nil
	}

	body := struct {
		Err error `json:"error,omitempty"`
	}{
		resp.Err,
	}

	return json.NewEncoder(w).Encode(body)
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

func getOrCreateSession(ctx context.Context, r *http.Request) context.Context {
	return ctx
}

func setSession(ctx context.Context, w http.ResponseWriter) context.Context {
	return ctx
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
