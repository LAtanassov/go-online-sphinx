package client

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
