package zcash

import (
	"fmt"

	jrpc "github.com/whyrusleeping/jrpc"
)

type Client struct {
	cli *jrpc.Client
}

func NewZcashClient(host, user, pass string) *Client {
	return &Client{
		cli: &jrpc.Client{
			Host: host,
			User: user,
			Pass: pass,
		},
	}
}

type TxResult struct {
	Hex      string
	Complete bool
	Errors   []interface{}
}

func (c *Client) SignRawTransaction(rawtx string) (string, error) {
	req := &jrpc.Request{
		Method: "signrawtransaction",
		Params: []string{rawtx},
	}

	sigout := jrpc.Response{ResultType: TxResult{}}
	err := c.cli.Do(req, &sigout)
	if err != nil {
		return "", err
	}
	if sigout.Error != nil {
		return "", sigout.Error
	}

	txr := sigout.Result.(*TxResult)
	if len(txr.Errors) > 0 {
		return "", fmt.Errorf("%v", txr.Errors[0])
	}

	return txr.Hex, nil
}

func (c *Client) SendRawTransaction(sigtx string) (string, error) {
	req := &jrpc.Request{
		Method: "sendrawtransaction",
		Params: []string{sigtx},
	}

	resp := jrpc.Response{ResultType: ""}
	err := c.cli.Do(req, &resp)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}

	return *resp.Result.(*string), nil
}

type infoTxOut struct {
	Address  string
	Amount   float64
	Category string
	Vout     int
	Size     int
}

type infoTransaction struct {
	Amount        float64
	Confirmations int
	Blockhash     string
	Blockindex    int
	Blocktime     int64
	Details       []infoTxOut
}

func (c *Client) GetTransactionValue(tx, addr string) (float64, error) {
	req := &jrpc.Request{
		Method: "gettransaction",
		Params: []string{tx},
	}

	var out struct {
		Result infoTransaction
	}
	err := c.cli.Do(req, &out)
	if err != nil {
		return 0, err
	}

	for _, res := range out.Result.Details {
		if res.Address == addr && res.Category == "receive" {
			return res.Amount, nil
		}
	}
	return 0, fmt.Errorf("no outputs for given address found")
}

type UnspentTx struct {
	Txid   string `json:"txid"`
	Vout   int
	Amount float64
}

func (c *Client) GetUnspents() ([]*UnspentTx, error) {
	req := &jrpc.Request{
		Method: "listunspent",
		Params: []interface{}{0},
	}

	resp := new(jrpc.Response)
	resp.ResultType = []*UnspentTx{}
	if err := c.cli.Do(req, resp); err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	return *resp.Result.(*[]*UnspentTx), nil
}
