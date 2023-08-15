package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"ntc-services/config"
	"ntc-services/models"
	"time"
)

type NtcPSBT struct {
	BaseURL string
	Client  *http.Client
}

func NewNtcPSBT() (*NtcPSBT, error) {
	baseURL, err := config.GetNtcBSPT()
	if err != nil {
		return nil, err
	}

	return &NtcPSBT{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: time.Second * 30,
		},
	}, nil
}

func (ntcp *NtcPSBT) PostPSBT(reqBody *PSBTReqBody) (*models.PSBTSerialized, error) {
	bodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := ntcp.post("/psbt", bodyJSON)
	if err != nil {
		return nil, err
	}

	var respBody *models.PSBTSerialized
	if err := json.Unmarshal(resp, &respBody); err != nil {
		return nil, err
	}

	return respBody, nil
}

func (ntcp *NtcPSBT) post(endpoint string, body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s%s", ntcp.BaseURL, endpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ntcp.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("ntc_psbt response was not 200")
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
