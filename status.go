package zcash

import (
	jrpc "github.com/whyrusleeping/jrpc"
)

type Info struct {
	Version         int
	ProtocolVersion int
	WalletVersion   int
	Balance         float64
	Blocks          int
	Timeoffset      int
	Connections     int
	Proxy           string
	Difficulty      float64
	Testnet         bool
	Keypoololdest   int
	Keypoolsize     int
	Paytxfee        float64
	Relayfee        float64
	Errors          string
}

func (c *Client) GetInfo() (*Info, error) {
	req := &jrpc.Request{
		Method: "getinfo",
	}

	resp := jrpc.Response{ResultType: Info{}}
	err := c.cli.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Result.(*Info), nil
}
