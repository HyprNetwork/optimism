package txmgr

import "encoding/json"

type DepositeClient struct {
	client *HTTP
}

func NewDepositeClient(rpc string) *DepositeClient {
	return &DepositeClient{client: NewHTTP(rpc)}
}

func (t *DepositeClient) IsDepositeExist() bool {
	return t.client.addr != ""
}

func (t *DepositeClient) GetDepositTx() ([]string, error) {
	data, err := t.client.Call("tx_getDepositTx")
	if err != nil {
		return nil, err
	}
	var result []string
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
