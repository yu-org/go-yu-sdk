package go_yu_sdk

import (
	"encoding/json"
	"github.com/yu-org/yu/core/protocol"
	"github.com/yu-org/yu/core/types"
	"io"
	"net/http"
	"path/filepath"
)

func (c *YuClient) StopChain() error {
	_, err := http.Get(filepath.Join(c.url, protocol.AdminApiPath, "stop"))
	return err
}

func (c *YuClient) GetReceipts() ([]*types.Receipt, error) {
	resp, err := http.Get(filepath.Join(c.url, protocol.RootApiPath, "receipts"))
	if err != nil {
		return nil, err
	}
	byt, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var apiResp protocol.APIResponse
	err = json.Unmarshal(byt, &apiResp)
	if err != nil {
		return nil, err
	}
	if apiResp.IsSuccess() {
		return apiResp.Data.([]*types.Receipt), nil
	}
	return nil, apiResp.Error()
}

func (c *YuClient) GetReceiptCount() (int, error) {
	resp, err := http.Get(filepath.Join(c.url, protocol.RootApiPath, "receipts_count"))
	if err != nil {
		return 0, err
	}
	byt, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var apiResp protocol.APIResponse
	err = json.Unmarshal(byt, &apiResp)
	if err != nil {
		return 0, err
	}
	if apiResp.IsSuccess() {
		return apiResp.Data.(int), nil
	}
	return 0, apiResp.Error()
}
