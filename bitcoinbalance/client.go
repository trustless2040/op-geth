package bitcoinbalance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"net/http"
	"strconv"
)

var btcBalanceModule *BtcBalanceModule = nil

type BtcBalanceModule struct {
	baseURL string
	client  *http.Client
}

// Make a GET request to the specified endpoint
func (c *BtcBalanceModule) get(endpoint string) ([]byte, error) {
	url := c.baseURL + endpoint
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Make a POST request to the specified endpoint with the given payload
func (c *BtcBalanceModule) post(endpoint string, payload interface{}) ([]byte, error) {
	url := c.baseURL + endpoint

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func init() {
	InitBitcoinBalance("https://sequencer-be.regtest.trustless.computer")
}

func InitBitcoinBalance(baseURL string) {
	btcBalanceModule = &BtcBalanceModule{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func GetBitcoinBalanceModule() *BtcBalanceModule {
	return btcBalanceModule
}

func (m *BtcBalanceModule) GetBalance(address []common.Address) (map[string]int64, error) {
	addressStr := []string{}
	for _, addr := range address {
		addressStr = append(addressStr, addr.String())
	}
	type Payload struct {
		Addresses []string `json:"addresses"`
	}
	balanceRes, err := m.post("/api/balances", Payload{addressStr})
	if err != nil {
		return nil, err
	}
	type Result struct {
		Balances []struct {
			Address          string `json:"address"`
			Balance          string `json:"balance"`
			AvailableBalance string `json:"available_balance"`
		} `json:"balances"`
	}
	result := Result{}
	err = json.Unmarshal(balanceRes, &result)
	if err != nil {
		return nil, err
	}
	returnRes := map[string]int64{}
	for _, balanceInfo := range result.Balances {
		availableBalance, _ := strconv.ParseInt(balanceInfo.AvailableBalance, 10, 64)
		returnRes[balanceInfo.Address] = availableBalance
	}
	return returnRes, nil
}

func (m *BtcBalanceModule) DeductBalance(address []common.Address, amount []uint64, txID []string, size []uint64, blockNumber uint64) error {
	type Payload struct {
		Tc_address string `json:"tc_address"`
		Amount     uint64 `json:"amount"`
		TX_ID      string `json:"tx_id"`
		Size       uint64 `json:"size"`
		L2BlockNum uint64 `json:"l2_block_num"`
	}
	payloadArr := []Payload{}
	for i, addr := range address {
		payloadArr = append(payloadArr, Payload{addr.String(), amount[i], txID[i], size[i], blockNumber})
	}
	response, err := m.post("/api/deduct", payloadArr)
	if err != nil {
		return err
	}

	result := map[string]string{}
	json.Unmarshal(response, &result)
	if result["result"] != "ok" {
		return fmt.Errorf("deduct balance failed")
	}
	return nil
}

type TxFeeAmount struct {
	TxID   string `json:"tx_id"`
	Amount uint64 `json:"amount"`
}

func (m *BtcBalanceModule) CommitBalance(txID []string, amount []uint64, info engine.InclusionPayloadAttrInfo) error {
	type Payload struct {
		TxList      []TxFeeAmount `json:"tx_list"`
		Batch_l1_tx string        `json:"batch_tx"`
		Btc_tx      string        `json:"btc_tx"`
		Btc_time    uint64        `json:"btc_time"`
		Btc_feerate float64       `json:"btc_fee_rate"`
		Btc_fee     int64         `json:"btc_fee"`
	}
	payloadArr := []TxFeeAmount{}
	for i, _ := range txID {
		payloadArr = append(payloadArr, TxFeeAmount{txID[i], amount[i]})
	}
	response, err := m.post("/api/commit", Payload{payloadArr, info.L1TxID, info.BTCTxId, info.BTCBlockTime, info.BTCFeeRate, info.BTCTotalFee})
	if err != nil {
		return err
	}

	result := map[string]string{}
	json.Unmarshal(response, &result)
	if result["result"] != "ok" {
		return fmt.Errorf("deduct balance failed")
	}
	return nil
}
