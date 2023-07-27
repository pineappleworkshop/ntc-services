package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"ntc-services/config"
	"time"
)

type Ordex struct {
	BaseURL string
	Client  *http.Client
}

func NewOrdex() (*Ordex, error) {
	baseURL, err := config.GetOrdexBaseURL()
	if err != nil {
		return nil, err
	}

	return &Ordex{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: time.Second * 30,
		},
	}, nil
}

func (o *Ordex) GetInscriptionById(id string) (interface{}, error) {
	resp, err := o.get("/inscriptions/" + id)
	if err != nil {
		return nil, err
	}

	var body interface{}
	if err := json.Unmarshal(resp, &body); err != nil {
		return nil, err
	}

	return body, nil
}

func (o *Ordex) get(endpoint string) ([]byte, error) {
	fmt.Println("baseURL:", o.BaseURL)
	req, err := http.NewRequest("GET", o.BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := o.Client.Do(req)
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
