package go_yu_sdk

import (
	"bytes"
	"encoding/json"
	"github.com/HyperService-Consortium/go-hexutil"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core"
	"github.com/yu-org/yu/core/keypair"
	"io"
	"net/http"
)

type YuClient struct {
	url     string
	privkey keypair.PrivKey
	pubkey  keypair.PubKey
}

func NewClient(url string) *YuClient {
	return &YuClient{url: url}
}

func (c *YuClient) WithKeys(privkey keypair.PrivKey, pubkey keypair.PubKey) *YuClient {
	c.privkey, c.pubkey = privkey, pubkey
	return c
}

func (c *YuClient) WriteChain(tripodName, funcName string, params any, leiPrice, tips uint64) error {
	paramsByt, err := json.Marshal(params)
	if err != nil {
		return err
	}
	wrCall := &common.WrCall{
		TripodName: tripodName,
		FuncName:   funcName,
		Params:     string(paramsByt),
		LeiPrice:   leiPrice,
		Tips:       tips,
	}

	byt, err := json.Marshal(wrCall)
	if err != nil {
		panic(err)
	}
	msgHash := common.BytesToHash(byt)
	sig, err := c.privkey.SignData(msgHash.Bytes())
	if err != nil {
		panic(err)
	}
	postBody := &core.WritingPostBody{
		Pubkey:    c.pubkey.StringWithType(),
		Signature: hexutil.Encode(sig),
		Call:      wrCall,
	}

	bodyByt, err := json.Marshal(postBody)
	if err != nil {
		return err
	}

	logrus.Debug("wrCall: ", c.url)

	_, err = http.Post(c.url, "application/json", bytes.NewReader(bodyByt))
	return err
}

func (c *YuClient) ReadChain(tripodName, funcName string, params any) ([]byte, error) {
	paramsByt, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	rdCall := &common.RdCall{
		TripodName: tripodName,
		FuncName:   funcName,
		Params:     string(paramsByt),
		BlockHash:  "",
	}
	byt, err := json.Marshal(rdCall)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(c.url, "application/json", bytes.NewReader(byt))
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}
