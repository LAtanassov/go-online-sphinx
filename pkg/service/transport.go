package service

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/LAtanassov/go-online-sphinx/pkg/contract"
	"github.com/pkg/errors"

	"github.com/gorilla/sessions"

	"github.com/go-kit/kit/log"
)

// TODO: not global and load key from ENV
// refactoring necessary to increase testability
var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

// ErrLoginRequired is return probably because of missing session
var ErrLoginRequired = errors.New("login required")

// HTTPTransport implements all HTTP Handler
type HTTPTransport struct {
	service Service
	logger  log.Logger
}

// NewHTTPTransport ...
func NewHTTPTransport(s Service, l log.Logger) *HTTPTransport {
	return &HTTPTransport{
		service: s,
		logger:  l,
	}
}

// MakeRegisterHandler ...
func (h *HTTPTransport) MakeRegisterHandler() http.Handler {
	return post("/v1/register", func(resp http.ResponseWriter, req *http.Request) {
		regReq, err := contract.UnmarshalRegisterRequest(req.Body)
		if err != nil {
			h.logger.Log("handler", "register", "error", fmt.Sprintf("+%v", errors.Wrapf(err, "UnmarshalRegisterRequest() failed")))
			contract.MarshalError(resp, err)
			return
		}
		defer req.Body.Close()

		err = h.service.Register(regReq.CID)
		if err != nil {
			h.logger.Log("handler", "register", "error", fmt.Sprintf("+%v", errors.Wrap(err, "Register() failed")))
			contract.MarshalError(resp, err)
			return
		}
		resp.WriteHeader(http.StatusCreated)
	})
}

// MakeExpKHandler ...
func (h *HTTPTransport) MakeExpKHandler() http.Handler {
	return post("/v1/login/expk", func(resp http.ResponseWriter, req *http.Request) {
		session, err := store.New(req, "online-sphinx")
		if err != nil {
			h.logger.Log("handler", "expk", "error", fmt.Sprintf("+%v", errors.Wrap(err, "session.Get() failed")))
			contract.MarshalError(resp, err)
			return
		}

		expkReq, err := contract.UnmarshalExpKRequest(req.Body)
		if err != nil {
			h.logger.Log("handler", "expk", "error", fmt.Sprintf("+%v", errors.Wrap(err, "UnmarshalExpKRequest() failed")))
			contract.MarshalError(resp, err)
			return
		}
		defer req.Body.Close()

		ski, sID, sNonce, bd, q0, kv, err := h.service.ExpK(expkReq.CID, expkReq.CNonce, expkReq.B, expkReq.Q)
		if err != nil {
			h.logger.Log("handler", "expk", "error", fmt.Sprintf("+%v", errors.Wrap(err, "ExpK() failed")))
			contract.MarshalError(resp, err)
			return
		}

		session.Values["sID"] = sID.Text(16)
		session.Values["cID"] = expkReq.CID.Text(16)
		session.Values["SKi"] = ski.Text(16)
		err = session.Save(req, resp)
		if err != nil {
			h.logger.Log("handler", "expk", "error", fmt.Sprintf("+%v", errors.Wrap(err, "session.Save() failed")))
			contract.MarshalError(resp, err)
			return
		}

		resp.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = contract.MarshalExpKResponse(resp, contract.ExpKResponse{SID: sID, SNonce: sNonce, BD: bd, Q0: q0, KV: kv})
		if err != nil {
			h.logger.Log("handler", "expk", "error", fmt.Sprintf("+%v", err))
			contract.MarshalError(resp, err)
			return
		}
	})
}

// MakeChallengeHandler ...
func (h *HTTPTransport) MakeChallengeHandler() http.Handler {
	return post("/v1/login/challenge", func(resp http.ResponseWriter, req *http.Request) {
		session, err := store.Get(req, "online-sphinx")
		if err != nil {
			h.logger.Log("handler", "challenge", "error", fmt.Sprintf("+%v", errors.Wrap(err, "session.Get() failed")))
			contract.MarshalError(resp, errors.Wrap(ErrLoginRequired, "called /v1/login/challenge without session"))
			return
		}

		skiHex, ok := session.Values["SKi"].(string)
		if !ok {
			h.logger.Log("handler", "challenge", "error", fmt.Sprintf("+%v", errors.Wrap(ErrLoginRequired, "session.Values() retrieved SKi failed")))
			contract.MarshalError(resp, errors.Wrap(ErrLoginRequired, "called /v1/login/challenge without session"))
			return
		}
		ski := new(big.Int)
		ski.SetString(skiHex, 16)

		challReq, err := contract.UnmarshalChallengeRequest(req.Body)
		if err != nil {
			h.logger.Log("handler", "challenge", "error", fmt.Sprintf("+%v", errors.Wrap(err, "UnmarshalChallengeRequest failed")))
			contract.MarshalError(resp, err)
			return
		}
		defer req.Body.Close()

		r, err := h.service.Challenge(ski, challReq.G, challReq.Q)
		if err != nil {
			h.logger.Log("handler", "challenge", "error", fmt.Sprintf("+%v", errors.Wrap(err, "Challenge() failed")))
			contract.MarshalError(resp, err)
			return
		}

		resp.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = contract.MarshalChallengeResponse(resp, contract.ChallengeResponse{R: r})
		if err != nil {
			h.logger.Log("handler", "challenge", "error", fmt.Sprintf("+%v", errors.Wrap(err, "MarshalChallengeResponse() failed")))
			contract.MarshalError(resp, err)
			return
		}
	})
}

// MakeLogoutHandler ...
func (h *HTTPTransport) MakeLogoutHandler() http.Handler {
	return post("/v1/logout", func(resp http.ResponseWriter, req *http.Request) {
		session, err := store.Get(req, "online-sphinx")
		if err != nil {
			resp.WriteHeader(http.StatusOK)
			return
		}
		defer req.Body.Close()

		session.Options.MaxAge = -1
		session.Save(req, resp)

		resp.WriteHeader(http.StatusOK)
	})
}

// MakeMetadataHandler ...
func (h *HTTPTransport) MakeMetadataHandler() http.Handler {
	return post("/v1/metadata", func(resp http.ResponseWriter, req *http.Request) {
		session, err := store.Get(req, "online-sphinx")
		if err != nil {
			h.logger.Log("handler", "metadata", "error", fmt.Sprintf("+%v", errors.Wrap(err, "session.Get() failed")))
			contract.MarshalError(resp, err)
			return
		}

		cIDHex, ok := session.Values["cID"].(string)
		if !ok {
			h.logger.Log("handler", "metadata", "error", fmt.Sprintf("+%v", errors.Wrap(ErrLoginRequired, "session.Values() retrieve cID failed")))
			contract.MarshalError(resp, ErrLoginRequired)
			return
		}
		cID := new(big.Int)
		cID.SetString(cIDHex, 16)

		skiHex, ok := session.Values["SKi"].(string)
		if !ok {
			h.logger.Log("handler", "metadata", "error", fmt.Sprintf("+%v", errors.Wrap(ErrLoginRequired, "session.Values() retrieve SKi failed")))
			contract.MarshalError(resp, ErrLoginRequired)
			return
		}
		ski := new(big.Int)
		ski.SetString(skiHex, 16)

		metaReq, err := contract.UnmarshalMetadataRequest(req.Body)
		if err != nil {
			h.logger.Log("handler", "metadata", "error", fmt.Sprintf("+%v", errors.Wrap(err, "UnmarshalMetadataRequest() failed")))
			contract.MarshalError(resp, err)
			return
		}
		defer req.Body.Close()

		err = h.service.VerifyMAC(metaReq.MAC, ski, []byte("metadata"))
		if err != nil {
			h.logger.Log("handler", "metadata", "error", fmt.Sprintf("+%v", errors.Wrap(err, "VerifyMAC() failed")))
			contract.MarshalError(resp, err)
			return
		}

		domains, err := h.service.GetMetadata(cID)
		if err != nil {
			h.logger.Log("handler", "metadata", "error", fmt.Sprintf("+%v", errors.Wrap(err, "GetMetadata() failed")))
			contract.MarshalError(resp, err)
			return
		}

		resp.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = contract.MarshalMetadataResponse(resp, contract.MetadataResponse{Domains: domains})
		if err != nil {
			h.logger.Log("handler", "metadata", "error", fmt.Sprintf("+%v", errors.Wrap(err, "MarshalMetadataResponse() failed")))
			contract.MarshalError(resp, err)
			return
		}
	})
}

// MakeAddHandler ...
func (h *HTTPTransport) MakeAddHandler() http.Handler {
	return post("/v1/add", func(resp http.ResponseWriter, req *http.Request) {

		session, err := store.Get(req, "online-sphinx")
		if err != nil {
			h.logger.Log("handler", "add", "error", fmt.Sprintf("+%v", errors.Wrap(err, "session.Get() failed")))
			contract.MarshalError(resp, err)
			return
		}

		cIDHex, ok := session.Values["cID"].(string)
		if !ok {
			h.logger.Log("handler", "add", "error", fmt.Sprintf("+%v", errors.Wrap(ErrLoginRequired, "session.Values() retrieve cID failed")))
			contract.MarshalError(resp, ErrLoginRequired)
			return
		}
		cID := new(big.Int)
		cID.SetString(cIDHex, 16)

		skiHex, ok := session.Values["SKi"].(string)
		if !ok {
			h.logger.Log("handler", "add", "error", fmt.Sprintf("+%v", errors.Wrap(ErrLoginRequired, "session.Values() retrieve SKi failed")))
			contract.MarshalError(resp, ErrLoginRequired)
			return
		}
		ski := new(big.Int)
		ski.SetString(skiHex, 16)

		addReq, err := contract.UnmarshalAddRequest(req.Body)
		if err != nil {
			h.logger.Log("handler", "add", "error", fmt.Sprintf("+%v", errors.Wrap(err, "UnmarshalAddRequest() failed")))
			contract.MarshalError(resp, err)
			return
		}
		defer req.Body.Close()

		err = h.service.VerifyMAC(addReq.MAC, ski, []byte(addReq.Domain))
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}

		err = h.service.Add(cID, addReq.Domain)
		if err != nil {
			h.logger.Log("handler", "add", "error", fmt.Sprintf("+%v", errors.Wrap(err, "Add() failed")))
			contract.MarshalError(resp, err)
			return
		}

		resp.WriteHeader(http.StatusCreated)
	})
}

// MakeGetHandler ...
func (h *HTTPTransport) MakeGetHandler() http.Handler {
	return post("/v1/get", func(resp http.ResponseWriter, req *http.Request) {

		session, err := store.Get(req, "online-sphinx")
		if err != nil {
			h.logger.Log("handler", "get", "error", fmt.Sprintf("+%v", errors.Wrap(err, "session.Get() failed")))
			contract.MarshalError(resp, err)
			return
		}

		cIDHex, ok := session.Values["cID"].(string)
		if !ok {
			h.logger.Log("handler", "get", "error", fmt.Sprintf("+%v", errors.Wrap(ErrLoginRequired, "session.Values() retrieve cID failed")))
			contract.MarshalError(resp, ErrLoginRequired)
			return
		}
		cID := new(big.Int)
		cID.SetString(cIDHex, 16)

		skiHex, ok := session.Values["SKi"].(string)
		if !ok {
			h.logger.Log("handler", "get", "error", fmt.Sprintf("+%v", errors.Wrap(ErrLoginRequired, "session.Values() retrieve SKi failed")))
			contract.MarshalError(resp, ErrLoginRequired)
			return
		}
		ski := new(big.Int)
		ski.SetString(skiHex, 16)

		getReq, err := contract.UnmarshalGetRequest(req.Body)
		if err != nil {
			h.logger.Log("handler", "get", "error", fmt.Sprintf("+%v", errors.Wrap(err, "UnmarshalGetRequest() failed")))
			contract.MarshalError(resp, err)
			return
		}
		defer req.Body.Close()

		err = h.service.VerifyMAC(getReq.MAC, ski, getReq.BMK.Bytes())
		if err != nil {
			h.logger.Log("handler", "get", "error", fmt.Sprintf("+%v", errors.Wrap(err, "VerifyMAC() failed")))
			contract.MarshalError(resp, err)
			return
		}
		bj, qj, err := h.service.Get(cID, getReq.Domain, getReq.BMK, getReq.Q)
		if err != nil {
			h.logger.Log("handler", "get", "error", fmt.Sprintf("+%v", errors.Wrap(err, "Get() failed")))
			contract.MarshalError(resp, err)
			return
		}

		err = contract.MarshalGetResponse(resp, contract.GetResponse{Bj: bj, Qj: qj})
		if err != nil {
			h.logger.Log("handler", "get", "error", fmt.Sprintf("+%v", errors.Wrap(err, "Get() failed")))
			contract.MarshalError(resp, err)
			return
		}

	})
}

// MakeLivenessHandler returns liveness handler
func (h *HTTPTransport) MakeLivenessHandler() http.Handler {
	return get("/_status/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}

// MakeReadinessHandler return readiness handler
func (h *HTTPTransport) MakeReadinessHandler() http.Handler {
	return get("/_status/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
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

func post(path string, f http.HandlerFunc) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(path, f).Methods("POST")
	return r
}

func get(path string, f http.HandlerFunc) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(path, f).Methods("GET")
	return r
}
