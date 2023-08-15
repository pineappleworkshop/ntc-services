package services

import (
	"bytes"
	"encoding/json"
	"errors"
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

func (o *Ordex) GetInscriptionByNumber(num int64) (interface{}, error) {
	resp, err := o.get(fmt.Sprintf("/inscriptions/byNumber?number=%v", num))
	if err != nil {
		return nil, err
	}

	var body interface{}
	if err := json.Unmarshal(resp, &body); err != nil {
		return nil, err
	}

	return body, nil
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

func (o *Ordex) GetInscriptionsByIds(ids []string) (interface{}, error) {
	requestBody := struct {
		Ids []string `json:"ids"`
	}{
		Ids: ids,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := o.post("/inscriptions/byIds", bodyBytes)
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

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Ordex response was not 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
func (o *Ordex) post(endpoint string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", o.BaseURL+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := o.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Ordex response was not 200")
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
