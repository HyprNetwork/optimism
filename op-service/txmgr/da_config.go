package txmgr

import (
	"encoding/hex"
	"encoding/json"
)

type DAConfig struct {
	client *HTTP
}

func NewDAConfig(rpc string) *DAConfig {
	return &DAConfig{
		client: NewHTTP(rpc),
	}
}

func (cfg *DAConfig) SetDA(value string) ([]byte, error) {
	data, err := cfg.client.Call("da_setDA", value)
	if err != nil {
		return nil, err
	}
	var result string
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return hex.DecodeString(result)
}

func (cfg *DAConfig) GetDA(key []byte) (string, error) {
	data, err := cfg.client.Call("da_getDA", hex.EncodeToString(key))
	if err != nil {
		return "", err
	}
	var result string
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}
	return result, nil
}
