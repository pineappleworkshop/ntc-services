package services

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"ntc-services/config"
	"time"
)

type BestInSlot struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

func NewBestInSlot() (*BestInSlot, error) {
	baseURL, err := config.GetBestInSlotBaseURL()
	if err != nil {
		return nil, err
	}
	apiKey, err := config.GetBestInSlotAPIKey()
	if err != nil {
		return nil, err
	}

	return &BestInSlot{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: time.Second * 30,
		},
	}, nil
}

func (bis *BestInSlot) GetBTCPrice() (float64, error) {
	resp, err := bis.get("/btc_price")

	if err != nil {
		return -1.0, err
	}

	var i uint64
	buf := bytes.NewReader(resp)
	if err := binary.Read(buf, binary.LittleEndian, &i); err != nil {
		return -1.0, err
	}

	price := math.Float64frombits(i)
	return price, nil
}

func (bis *BestInSlot) get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", bis.BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bis.APIKey))
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
