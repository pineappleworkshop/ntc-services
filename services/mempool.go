package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"ntc-services/config"
	"time"
)

type Mempool struct {
	BaseURL string
	Client  *http.Client
}

func NewMempool() (*Mempool, error) {
	baseURL, err := config.GetMempoolBaseURL()
	if err != nil {
		return nil, err
	}

	return &Mempool{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: time.Second * 30,
		},
	}, nil
}

func (m *Mempool) GetRecommendedFees() (interface{}, error) {
	resp, err := m.get("/v1/fees/recommended")
	if err != nil {
		return nil, err
	}

	var body interface{}
	if err := json.Unmarshal(resp, &body); err != nil {
		return nil, err
	}

	return body, nil
}

func (m *Mempool) GetBlockHeight() (int64, error) {
	resp, err := m.get("/blocks/tip/height")
	if err != nil {
		return -1, err
	}

	var height int64
	if err := json.Unmarshal(resp, &height); err != nil {
		return -1, err
	}

	return height, nil
}

func (m *Mempool) get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", m.BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.Client.Do(req)
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
