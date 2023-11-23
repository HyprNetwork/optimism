package txmgr

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/valyala/fasthttp"
)

type RPCRequest struct {
	Version string          `json:"jsonrpc"`
	ID      uint64          `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type RPCResponse struct {
	Version string       `json:"jsonrpc"`
	ID      uint64       `json:"id"`
	Result  string       `json:"result,omitempty"`
	Error   *ErrorObject `json:"error,omitempty"`
}

type ErrorObject struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *ErrorObject) Error() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("jsonrpc.internal marshal error: %v", err)
	}
	return string(data)
}

type HTTP struct {
	addr   string
	client *fasthttp.Client
}

func NewHTTP(addr string) *HTTP {
	return &HTTP{
		addr: addr,
		client: &fasthttp.Client{
			Dial: func(addr string) (net.Conn, error) {
				return fasthttp.DialTimeout(addr, time.Minute)
			},
		},
	}
}

func (h *HTTP) Call(method string, params ...interface{}) (string, error) {
	request := RPCRequest{
		Method:  method,
		Version: "2.0",
	}
	if len(params) > 0 {
		data, err := json.Marshal(params)
		if err != nil {
			return "", err
		}
		// fmt.Printf("data %s\n", data)
		request.Params = data
	}
	raw, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(h.addr)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetBody(raw)

	if err := h.client.Do(req, res); err != nil {
		return "", err
	}

	var response RPCResponse
	if err := json.Unmarshal(res.Body(), &response); err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", errors.New(response.Error.Error())
	}
	return response.Result, nil
}

type DAConfig struct {
	Rpc string
}

func NewDAConfig(rpc string) *DAConfig {
	return &DAConfig{
		Rpc: rpc,
	}
}

func (cfg *DAConfig) SetDA(value string) ([]byte, error) {
	res, err := NewHTTP(cfg.Rpc).Call("da_setDA", value)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(res)
}

func (cfg *DAConfig) GetDA(key []byte) (string, error) {
	return NewHTTP(cfg.Rpc).Call("da_getDA", hex.EncodeToString(key))
}
