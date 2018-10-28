package contract

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math/big"
)

// ErrUnexpectedType is returned after a type cast failed.
var ErrUnexpectedType = errors.New("unexpected type")

// MarshalRegisterRequest ...
func MarshalRegisterRequest(w io.Writer, r RegisterRequest) error {
	body := struct {
		CID string `json:"CID"`
	}{
		CID: r.CID.Text(16),
	}
	return json.NewEncoder(w).Encode(body)
}

// UnmarshalRegisterRequest from json as byte array to struct
func UnmarshalRegisterRequest(r io.Reader) (RegisterRequest, error) {
	var body struct {
		CID string `json:"CID"`
	}

	if err := json.NewDecoder(r).Decode(&body); err != nil {
		return RegisterRequest{}, err
	}

	cID := new(big.Int)
	cID.SetString(body.CID, 16)

	return RegisterRequest{
		CID: cID,
	}, nil
}

// RegisterRequest ...
type RegisterRequest struct {
	CID *big.Int
}

// MarshalExpKRequest ...
func MarshalExpKRequest(w io.Writer, r ExpKRequest) error {
	body := struct {
		CID    string
		CNonce string
		B      string
		Q      string
	}{
		CID:    r.CID.Text(16),
		CNonce: r.CNonce.Text(16),
		B:      r.B.Text(16),
		Q:      r.Q.Text(16),
	}
	return json.NewEncoder(w).Encode(body)
}

// UnmarshalExpKRequest ...
func UnmarshalExpKRequest(r io.Reader) (ExpKRequest, error) {
	var body struct {
		CID    string `json:"cID"`
		CNonce string `json:"cNonce"`
		B      string `json:"b"`
		Q      string `json:"q"`
	}

	if err := json.NewDecoder(r).Decode(&body); err != nil {
		return ExpKRequest{}, err
	}

	cID := new(big.Int)
	cID.SetString(body.CID, 16)

	cNonce := new(big.Int)
	cNonce.SetString(body.CNonce, 16)

	b := new(big.Int)
	b.SetString(body.B, 16)

	q := new(big.Int)
	q.SetString(body.Q, 16)

	return ExpKRequest{
		CID:    cID,
		CNonce: cNonce,
		B:      b,
		Q:      q,
	}, nil
}

// ExpKRequest ...
type ExpKRequest struct {
	CID    *big.Int
	CNonce *big.Int
	B      *big.Int
	Q      *big.Int
}

// MarshalExpKResponse ...
func MarshalExpKResponse(w io.Writer, r ExpKResponse) error {

	body := struct {
		SID    string `json:"sID"`
		SNonce string `json:"sNonce"`
		BD     string `json:"bd"`
		Q0     string `json:"q0"`
		KV     string `json:"kv"`
	}{
		r.SID.Text(16),
		r.SNonce.Text(16),
		r.BD.Text(16),
		r.Q0.Text(16),
		r.KV.Text(16),
	}

	return json.NewEncoder(w).Encode(body)
}

// UnmarshalExpKResponse ...
func UnmarshalExpKResponse(r io.Reader) (ExpKResponse, error) {
	var body struct {
		SID    string `json:"sID"`
		SNonce string `json:"sNonce"`
		BD     string `json:"bd"`
		Q0     string `json:"q0"`
		KV     string `json:"kv"`
	}

	if err := json.NewDecoder(r).Decode(&body); err != nil {
		return ExpKResponse{}, err
	}

	sID := new(big.Int)
	sID.SetString(body.SID, 16)

	sNonce := new(big.Int)
	sNonce.SetString(body.SNonce, 16)

	bd := new(big.Int)
	bd.SetString(body.BD, 16)

	q0 := new(big.Int)
	q0.SetString(body.Q0, 16)

	kv := new(big.Int)
	kv.SetString(body.KV, 16)

	return ExpKResponse{
		SID:    sID,
		SNonce: sNonce,
		BD:     bd,
		Q0:     q0,
		KV:     kv,
	}, nil
}

// ExpKResponse ...
type ExpKResponse struct {
	SID    *big.Int
	SNonce *big.Int
	BD     *big.Int
	Q0     *big.Int
	KV     *big.Int
}

// MarshalChallengeRequest ...
func MarshalChallengeRequest(w io.Writer, r ChallengeRequest) error {
	body := struct {
		G string
		Q string
	}{
		G: r.G.Text(16),
		Q: r.Q.Text(16),
	}
	return json.NewEncoder(w).Encode(body)
}

// UnmarshalChallengeRequest ...
func UnmarshalChallengeRequest(r io.Reader) (ChallengeRequest, error) {
	var body struct {
		G string `json:"g"`
		Q string `json:"q"`
	}

	if err := json.NewDecoder(r).Decode(&body); err != nil {
		return ChallengeRequest{}, err
	}

	g := new(big.Int)
	g.SetString(body.G, 16)

	q := new(big.Int)
	q.SetString(body.Q, 16)

	return ChallengeRequest{
		G: g,
		Q: q,
	}, nil
}

// ChallengeRequest ...
type ChallengeRequest struct {
	G *big.Int
	Q *big.Int
}

// MarshalChallengeResponse ...
func MarshalChallengeResponse(w io.Writer, r ChallengeResponse) error {

	body := struct {
		R string `json:"r"`
	}{
		r.R.Text(16),
	}

	return json.NewEncoder(w).Encode(body)
}

// UnmarshalChallengeResponse ...
func UnmarshalChallengeResponse(r io.Reader) (ChallengeResponse, error) {
	var body struct {
		R string `json:"r"`
	}

	if err := json.NewDecoder(r).Decode(&body); err != nil {
		return ChallengeResponse{}, err
	}

	rv := new(big.Int)
	rv.SetString(body.R, 16)

	return ChallengeResponse{
		R: rv,
	}, nil
}

// ChallengeResponse ...
type ChallengeResponse struct {
	R *big.Int
}

// MarshalMetadataRequest ...
func MarshalMetadataRequest(w io.Writer, r MetadataRequest) error {

	body := struct {
		MAC string `json:"mac"`
	}{
		hex.EncodeToString(r.MAC),
	}

	return json.NewEncoder(w).Encode(body)
}

// UnmarshalMetadataRequest ...
func UnmarshalMetadataRequest(r io.Reader) (MetadataRequest, error) {
	var body struct {
		MAC string `json:"mac"`
	}

	if err := json.NewDecoder(r).Decode(&body); err != nil {
		return MetadataRequest{}, err
	}

	mac, err := hex.DecodeString(body.MAC)
	if err != nil {
		return MetadataRequest{}, err
	}

	return MetadataRequest{
		MAC: mac,
	}, nil
}

// MetadataRequest ...
type MetadataRequest struct {
	MAC []byte
}

// MarshalMetadataResponse ...
func MarshalMetadataResponse(w io.Writer, r MetadataResponse) error {
	body := struct {
		Domains []string `json:"domains"`
	}{
		r.Domains,
	}

	return json.NewEncoder(w).Encode(body)
}

// UnmarshalMetadataResponse ...
func UnmarshalMetadataResponse(r io.Reader) (MetadataResponse, error) {
	var body struct {
		Domains []string `json:"domains"`
	}

	if err := json.NewDecoder(r).Decode(&body); err != nil {
		return MetadataResponse{}, err
	}

	return MetadataResponse{
		Domains: body.Domains,
	}, nil
}

// MetadataResponse ...
type MetadataResponse struct {
	Domains []string
}

// MarshalAddRequest ...
func MarshalAddRequest(w io.Writer, r AddRequest) error {
	body := struct {
		MAC    string `json:"mac"`
		Domain string `json:"domain"`
	}{
		hex.EncodeToString(r.MAC),
		r.Domain,
	}

	return json.NewEncoder(w).Encode(body)
}

// UnmarshalAddRequest ...
func UnmarshalAddRequest(r io.Reader) (AddRequest, error) {
	var body struct {
		MAC    string `json:"mac"`
		Domain string `json:"domain"`
	}

	if err := json.NewDecoder(r).Decode(&body); err != nil {
		return AddRequest{}, err
	}

	mac, err := hex.DecodeString(body.MAC)
	if err != nil {
		return AddRequest{}, err
	}

	return AddRequest{
		MAC:    mac,
		Domain: body.Domain,
	}, nil
}

// AddRequest ...
type AddRequest struct {
	MAC    []byte
	Domain string
}

// MarshalGetRequest ...
func MarshalGetRequest(w io.Writer, r GetRequest) error {
	body := struct {
		MAC    string `json:"mac"`
		Domain string `json:"domain"`
		BMK    string `json:"bmk"`
		Q      string `json:"q"`
	}{
		hex.EncodeToString(r.MAC),
		r.Domain,
		r.BMK.Text(16),
		r.Q.Text(16),
	}

	return json.NewEncoder(w).Encode(body)
}

// UnmarshalGetRequest ...
func UnmarshalGetRequest(r io.Reader) (GetRequest, error) {
	var body struct {
		MAC    string `json:"mac"`
		Domain string `json:"domain"`
		BMK    string `json:"bmk"`
		Q      string `json:"q"`
	}

	if err := json.NewDecoder(r).Decode(&body); err != nil {
		return GetRequest{}, err
	}

	bmk := new(big.Int)
	_, ok := bmk.SetString(body.BMK, 16)
	if !ok {
		return GetRequest{}, ErrUnexpectedType
	}

	q := new(big.Int)
	_, ok = q.SetString(body.Q, 16)
	if !ok {
		return GetRequest{}, ErrUnexpectedType
	}

	mac, err := hex.DecodeString(body.MAC)
	if err != nil {
		return GetRequest{}, err
	}

	return GetRequest{
		MAC:    mac,
		Domain: body.Domain,
		BMK:    bmk,
		Q:      q,
	}, nil
}

// GetRequest ...
type GetRequest struct {
	MAC    []byte
	Domain string
	BMK    *big.Int
	Q      *big.Int
}

// MarshalGetResponse ...
func MarshalGetResponse(w io.Writer, r GetResponse) error {

	body := struct {
		Bj string `json:"bj"`
		Qj string `json:"qj"`
	}{
		r.Bj.Text(16),
		r.Qj.Text(16),
	}

	return json.NewEncoder(w).Encode(body)
}

// UnmarshalGetResponse ...
func UnmarshalGetResponse(r io.Reader) (GetResponse, error) {
	var body struct {
		Bj string `json:"bj"`
		Qj string `json:"qj"`
	}

	if err := json.NewDecoder(r).Decode(&body); err != nil {
		return GetResponse{}, err
	}

	bj := new(big.Int)
	_, ok := bj.SetString(body.Bj, 16)
	if !ok {
		return GetResponse{}, ErrUnexpectedType
	}

	qj := new(big.Int)
	_, ok = qj.SetString(body.Qj, 16)
	if !ok {
		return GetResponse{}, ErrUnexpectedType
	}

	return GetResponse{
		Bj: bj,
		Qj: qj,
	}, nil
}

// GetResponse ...
type GetResponse struct {
	Bj *big.Int
	Qj *big.Int
}
