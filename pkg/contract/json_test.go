package contract

import (
	"bufio"
	"bytes"
	"math/big"
	"reflect"
	"testing"
)

func TestRegisterRequest(t *testing.T) {
	t.Run("marshal and unmarshal", func(t *testing.T) {
		want := RegisterRequest{
			CID: big.NewInt(1),
		}

		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		err := MarshalRegisterRequest(w, want)
		if err != nil {
			t.Errorf("MarshalRegisterRequest() error = %v", err)
			return
		}
		w.Flush()

		r := bufio.NewReader(&buf)
		got, err := UnmarshalRegisterRequest(r)
		if err != nil {
			t.Errorf("UnmarshalRegisterRequest() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("RegisterRequest = %v, want %v", got, want)
		}
	})
}

func TestExpKRequest(t *testing.T) {
	t.Run("marshal and unmarshal", func(t *testing.T) {
		want := ExpKRequest{
			CID:    big.NewInt(1),
			CNonce: big.NewInt(2),
			B:      big.NewInt(3),
			Q:      big.NewInt(4),
		}
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		err := MarshalExpKRequest(w, want)
		if err != nil {
			t.Errorf("MarshalExpKRequest() error = %v", err)
			return
		}
		w.Flush()

		r := bufio.NewReader(&buf)
		got, err := UnmarshalExpKRequest(r)
		if err != nil {
			t.Errorf("UnmarshalExpKRequest() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("ExpKRequest = %v, want %v", got, want)
		}
	})
}

func TestExpKResponse(t *testing.T) {
	want := ExpKResponse{
		SID:    big.NewInt(1),
		SNonce: big.NewInt(2),
		BD:     big.NewInt(3),
		Q0:     big.NewInt(4),
		KV:     big.NewInt(5),
	}
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err := MarshalExpKResponse(w, want)
	if err != nil {
		t.Errorf("MarshalExpKResponse() error = %v", err)
		return
	}
	w.Flush()

	r := bufio.NewReader(&buf)

	got, err := UnmarshalExpKResponse(r)
	if err != nil {
		t.Errorf("UnmarshalExpKResponse() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExpKResponse = %v, want %v", got, want)
	}
}

func TestChallengeRequest(t *testing.T) {
	want := ChallengeRequest{
		G: big.NewInt(1),
		Q: big.NewInt(2),
	}
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err := MarshalChallengeRequest(w, want)
	if err != nil {
		t.Errorf("MarshalChallengeRequest() error = %v", err)
		return
	}
	w.Flush()

	r := bufio.NewReader(&buf)

	got, err := UnmarshalChallengeRequest(r)
	if err != nil {
		t.Errorf("UnmarshalChallengeRequest() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ChallengeRequest = %v, want %v", got, want)
	}
}

func TestChallengeResponse(t *testing.T) {
	want := ChallengeResponse{
		R: big.NewInt(1),
	}
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err := MarshalChallengeResponse(w, want)
	if err != nil {
		t.Errorf("MarshalChallengeResponse() error = %v", err)
		return
	}
	w.Flush()

	r := bufio.NewReader(&buf)

	got, err := UnmarshalChallengeResponse(r)
	if err != nil {
		t.Errorf("UnmarshalChallengeResponse() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ChallengeResponse = %v, want %v", got, want)
	}
}

func TestMetadataRequest(t *testing.T) {
	want := MetadataRequest{
		MAC: []byte("mac"),
	}
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err := MarshalMetadataRequest(w, want)
	if err != nil {
		t.Errorf("MarshalMetadataRequest() error = %v", err)
		return
	}
	w.Flush()

	r := bufio.NewReader(&buf)

	got, err := UnmarshalMetadataRequest(r)
	if err != nil {
		t.Errorf("UnmarshalMetadataRequest() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MarshalMetadataRequest = %v, want %v", got, want)
	}
}

func TestMetadataResponse(t *testing.T) {
	want := MetadataResponse{
		Domains: []string{"domain"},
	}
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err := MarshalMetadataResponse(w, want)
	if err != nil {
		t.Errorf("MarshalMetadataResponse() error = %v", err)
		return
	}
	w.Flush()

	r := bufio.NewReader(&buf)

	got, err := UnmarshalMetadataResponse(r)
	if err != nil {
		t.Errorf("UnmarshalMetadataResponse() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MetadataResponse = %v, want %v", got, want)
	}
}

func TestAddRequest(t *testing.T) {
	want := AddRequest{
		MAC:    []byte("mac"),
		Domain: "domain",
	}
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err := MarshalAddRequest(w, want)
	if err != nil {
		t.Errorf("MarshalAddRequest() error = %v", err)
		return
	}
	w.Flush()

	r := bufio.NewReader(&buf)

	got, err := UnmarshalAddRequest(r)
	if err != nil {
		t.Errorf("UnmarshalAddRequest() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AddRequest = %v, want %v", got, want)
	}
}

func TestGetRequest(t *testing.T) {
	want := GetRequest{
		MAC:    []byte("mac"),
		Domain: "domain",
		BMK:    big.NewInt(2),
		Q:      big.NewInt(3),
	}
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err := MarshalGetRequest(w, want)
	if err != nil {
		t.Errorf("MarshalGetRequest() error = %v", err)
		return
	}
	w.Flush()

	r := bufio.NewReader(&buf)

	got, err := UnmarshalGetRequest(r)
	if err != nil {
		t.Errorf("UnmarshalGetRequest() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetRequest = %v, want %v", got, want)
	}
}

func TestUnmarshalGetResponse(t *testing.T) {
	want := GetResponse{
		Bj: big.NewInt(2),
		Qj: big.NewInt(3),
	}
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err := MarshalGetResponse(w, want)
	if err != nil {
		t.Errorf("MarshalGetResponse() error = %v", err)
		return
	}
	w.Flush()

	r := bufio.NewReader(&buf)

	got, err := UnmarshalGetResponse(r)
	if err != nil {
		t.Errorf("UnmarshalGetResponse() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetResponse = %v, want %v", got, want)
	}
}
