package client

import (
	"encoding/json"
	"io"
	"math/big"
)

type registerRequest struct {
	Username string `json:"username"`
}

type expKRequest struct {
	CID    string `json:"cID"`
	CNonce string `json:"cNonce"`
	B      string `json:"b"`
	Q      string `json:"q"`
}

type expKResponse struct {
	SID    string `json:"sID"`
	SNonce string `json:"sNonce"`
	BD     string `json:"bd"`
	Q0     string `json:"Q0"`
	KV     string `json:"kv"`
	Err    error  `json:"error"`
}

type verifyRequest struct {
	G string `json:"g"`
	Q string `json:"q"`
}

type verifyResponse struct {
	R   string `json:"r"`
	Err error  `json:"error"`
}

type metadataRequest struct {
	CID string `json:"cID"`
	MAC string `json:"mac"`
}

type metadataResponse struct {
	Domains []Domain `json:"domains"`
	Err     error    `json:"error"`
}

func marshalRegisterRequest(username string) ([]byte, error) {
	return json.Marshal(&registerRequest{Username: username})
}

func marshalExpKRequest(cID, cNonce, b, q *big.Int) ([]byte, error) {
	return json.Marshal(&expKRequest{
		CID:    cID.Text(16),
		CNonce: cNonce.Text(16),
		B:      b.Text(16),
		Q:      q.Text(16),
	})
}

func unmarsalExpKResponse(r io.Reader) (sID, sNonce, bd, kv, q0 *big.Int, err error) {

	resp := expKResponse{}
	err = json.NewDecoder(r).Decode(&resp)
	if err != nil {
		return
	}

	if resp.Err != nil {
		err = resp.Err
		return
	}

	bd = new(big.Int)
	bd.SetString(resp.BD, 16)

	q0 = new(big.Int)
	q0.SetString(resp.Q0, 16)

	kv = new(big.Int)
	kv.SetString(resp.KV, 16)

	sID = new(big.Int)
	sID.SetString(resp.SID, 16)

	sNonce = new(big.Int)
	sNonce.SetString(resp.SNonce, 16)
	return
}

func marshalExpKResponse(sID, sNonce, bd, kv, q0 *big.Int) ([]byte, error) {
	return json.Marshal(&expKResponse{
		SID:    sID.Text(16),
		SNonce: sNonce.Text(16),
		BD:     bd.Text(16),
		KV:     kv.Text(16),
		Q0:     q0.Text(16),
	})
}

func marshalVerifyRequest(g, q *big.Int) ([]byte, error) {
	return json.Marshal(&verifyRequest{
		G: g.Text(16),
		Q: q.Text(16),
	})
}

func marshalVerifyResponse(r *big.Int, err error) ([]byte, error) {
	return json.Marshal(&verifyResponse{
		R:   r.Text(16),
		Err: err,
	})
}

func unmarsalVerifyResponse(rd io.Reader) (*big.Int, error) {

	resp := verifyResponse{}
	err := json.NewDecoder(rd).Decode(&resp)
	if err != nil {
		return nil, err
	}

	if resp.Err != nil {
		return nil, resp.Err
	}

	r := new(big.Int)
	r.SetString(resp.R, 16)

	return r, nil
}

func marshalMetadataRequest(cID, mac string) ([]byte, error) {
	return json.Marshal(&metadataRequest{
		CID: cID,
		MAC: mac,
	})
}

func marshalMetadataResponse(domains []Domain, err error) ([]byte, error) {
	return json.Marshal(&metadataResponse{
		Domains: domains,
		Err:     err,
	})
}

func unmarsalMetadataResponse(r io.Reader) ([]Domain, error) {

	resp := metadataResponse{}
	err := json.NewDecoder(r).Decode(&resp)
	if err != nil {
		return nil, err
	}

	if resp.Err != nil {
		return nil, resp.Err
	}

	return resp.Domains, nil
}
