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
	Username string `json:"username"`
	CNonce   string `json:"cNonce"`
	B        string `json:"b"`
	Q        string `json:"q"`
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
	VNonce string `json:"vNonce"`
}

type verifyResponse struct {
	WNonce string `json:"wNonce"`
}

type metadataRequest struct {
}

func marshalRegisterRequest(username string) ([]byte, error) {
	return json.Marshal(&registerRequest{Username: username})
}

func marshalExpKRequest(cNonce, b, q *big.Int) ([]byte, error) {
	return json.Marshal(&expKRequest{
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
