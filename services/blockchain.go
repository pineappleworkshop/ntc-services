package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"ntc-services/config"
	"time"
)

type BlockChain struct {
	APIKey    string
	BaseURL   string
	BaseURLV3 string
	Client    *http.Client
}

func NewBlockChain() (*BlockChain, error) {
	baseURL, err := config.GetBlockChainBaseURL()
	if err != nil {
		return nil, err
	}

	return &BlockChain{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: time.Second * 30,
		},
	}, nil
}

func (bis *BlockChain) GetBTCPrice() (float64, error) {
	resp, err := bis.get("/tickers/BTC-USD")
	if err != nil {
		return -1.0, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return -1.0, err
	}

	lastTradePrice, ok := result["last_trade_price"].(float64)
	if !ok {
		fmt.Println("Invalid type for last_trade_price")
		return -1.0, err
	}

	return lastTradePrice, nil
}

func (bis *BlockChain) get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", bis.BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	//req.Header.Set("x-api-key", fmt.Sprintf("Bearer %s", bis.APIKey))
	resp, err := bis.Client.Do(req)
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
