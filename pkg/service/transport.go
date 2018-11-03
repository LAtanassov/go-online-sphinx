package service

import (
	"math/big"
	"net/http"

	"github.com/LAtanassov/go-online-sphinx/pkg/contract"
	"github.com/pkg/errors"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// TODO: not global and load key from ENV
// refactoring necessary to increase testability
var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

// ErrLoginRequired is return probably because of missing session
var ErrLoginRequired = errors.New("login required")

// MakeRegisterHandler returns a handler
func MakeRegisterHandler(s Service) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/v1/register", func(resp http.ResponseWriter, req *http.Request) {
		regReq, err := contract.UnmarshalRegisterRequest(req.Body)
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}
		defer req.Body.Close()

		err = s.Register(regReq.CID)
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}
		resp.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	return r
}

// MakeExpKHandler returns a handler for the handling service.
func MakeExpKHandler(s Service) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/v1/login/expk", func(resp http.ResponseWriter, req *http.Request) {

		session, err := store.Get(req, "online-sphinx")
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}

		expkReq, err := contract.UnmarshalExpKRequest(req.Body)
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}
		defer req.Body.Close()

		ski, sID, sNonce, bd, q0, kv, err := s.ExpK(expkReq.CID, expkReq.CNonce, expkReq.B, expkReq.Q)
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}

		session.Values["sID"] = sID.Text(16)
		session.Values["cID"] = expkReq.CID.Text(16)
		session.Values["SKi"] = ski.Text(16)
		session.Save(req, resp)

		resp.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = contract.MarshalExpKResponse(resp, contract.ExpKResponse{SID: sID, SNonce: sNonce, BD: bd, Q0: q0, KV: kv})
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}

	}).Methods("POST")

	return r
}

// MakeChallengeHandler returns a handler for the handling service.
func MakeChallengeHandler(s Service) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/v1/login/challenge", func(resp http.ResponseWriter, req *http.Request) {

		session, err := store.Get(req, "online-sphinx")
		if err != nil {
			contract.MarshalError(resp, ErrLoginRequired)
			return
		}

		skiHex, ok := session.Values["SKi"].(string)
		if !ok {
			contract.MarshalError(resp, ErrLoginRequired)
			return
		}
		ski := new(big.Int)
		ski.SetString(skiHex, 16)

		challReq, err := contract.UnmarshalChallengeRequest(req.Body)
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}
		defer req.Body.Close()

		r, err := s.Challenge(ski, challReq.G, challReq.Q)
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}

		resp.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = contract.MarshalChallengeResponse(resp, contract.ChallengeResponse{R: r})
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}

	}).Methods("POST")

	return r
}

// MakeMetadataHandler ...
func MakeMetadataHandler(s Service) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/v1/metadata", func(resp http.ResponseWriter, req *http.Request) {

		session, err := store.Get(req, "online-sphinx")
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}

		cIDHex, ok := session.Values["cID"].(string)
		if !ok {
			contract.MarshalError(resp, ErrLoginRequired)
			return
		}
		cID := new(big.Int)
		cID.SetString(cIDHex, 16)

		skiHex, ok := session.Values["SKi"].(string)
		if !ok {
			contract.MarshalError(resp, ErrLoginRequired)
			return
		}
		ski := new(big.Int)
		ski.SetString(skiHex, 16)

		metaReq, err := contract.UnmarshalMetadataRequest(req.Body)
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}
		defer req.Body.Close()

		err = s.VerifyMAC(metaReq.MAC, ski, []byte("metadata"))
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}

		domains, err := s.GetMetadata(cID)
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}

		resp.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = contract.MarshalMetadataResponse(resp, contract.MetadataResponse{Domains: domains})
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}

	}).Methods("POST")

	return r
}

// MakeAddHandler ...
func MakeAddHandler(s Service) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/v1/add", func(resp http.ResponseWriter, req *http.Request) {

		session, err := store.Get(req, "online-sphinx")
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}

		cIDHex, ok := session.Values["cID"].(string)
		if !ok {
			contract.MarshalError(resp, ErrLoginRequired)
			return
		}
		cID := new(big.Int)
		cID.SetString(cIDHex, 16)

		skiHex, ok := session.Values["SKi"].(string)
		if !ok {
			contract.MarshalError(resp, ErrLoginRequired)
			return
		}
		ski := new(big.Int)
		ski.SetString(skiHex, 16)

		addReq, err := contract.UnmarshalAddRequest(req.Body)
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}
		defer req.Body.Close()

		err = s.VerifyMAC(addReq.MAC, ski, []byte(addReq.Domain))
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}

		err = s.Add(cID, addReq.Domain)
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}

		resp.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	return r
}

// MakeGetHandler ...
func MakeGetHandler(s Service) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/v1/get", func(resp http.ResponseWriter, req *http.Request) {
		session, err := store.Get(req, "online-sphinx")
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}

		cIDHex, ok := session.Values["cID"].(string)
		if !ok {
			contract.MarshalError(resp, ErrLoginRequired)
			return
		}
		cID := new(big.Int)
		cID.SetString(cIDHex, 16)

		skiHex, ok := session.Values["SKi"].(string)
		if !ok {
			contract.MarshalError(resp, ErrLoginRequired)
			return
		}
		ski := new(big.Int)
		ski.SetString(skiHex, 16)

		getReq, err := contract.UnmarshalGetRequest(req.Body)
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}
		defer req.Body.Close()

		err = s.VerifyMAC(getReq.MAC, ski, getReq.BMK.Bytes())
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}
		bj, qj, err := s.Get(cID, getReq.Domain, getReq.BMK, getReq.Q)
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}

		err = contract.MarshalGetResponse(resp, contract.GetResponse{Bj: bj, Qj: qj})
		if err != nil {
			contract.MarshalError(resp, err)
			return
		}

	}).Methods("POST")

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
