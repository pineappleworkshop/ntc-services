package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"ntc-services/config"
	"time"
)

type BlockChainInfo struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

func NewBlockChainInfo() (*BlockChainInfo, error) {
	baseURL, err := config.GetBlockChainInfoBaseURL()
	if err != nil {
		return nil, err
	}

	return &BlockChainInfo{
		APIKey:  "",
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: time.Second * 30,
		},
	}, nil
}

func (bci *BlockChainInfo) GetUTXOsForAddr(addr string) (map[string]interface{}, error) {
	resp, err := bci.get(fmt.Sprintf("/unspent?active=%s", addr))
	if err != nil {
		fmt.Printf("Could not assemble URL: %+v \n", err)
		return nil, err
	}

	var result map[string]interface{}
	if err = json.Unmarshal(resp, &result); err != nil {
		fmt.Println("Error unmarshaling JSON:", err) // TODO: better logging
		return nil, err
	}

	return result, nil
}

func (bci *BlockChainInfo) get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", bci.BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := bci.Client.Do(req)
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
