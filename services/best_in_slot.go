package services

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"ntc-services/config"
	"ntc-services/models"
	"strconv"
	"time"
)

type BestInSlot struct {
	APIKey    string
	BaseURL   string
	BaseURLV3 string
	Client    *http.Client
}

func NewBestInSlot() (*BestInSlot, error) {
	apiKey, err := config.GetBestInSlotAPIKey()
	if err != nil {
		return nil, err
	}
	baseURL, err := config.GetBestInSlotBaseURL()
	if err != nil {
		return nil, err
	}
	baseURLV3, err := config.GetBestInSlotBaseURLV3()
	if err != nil {
		return nil, err
	}

	return &BestInSlot{
		APIKey:    apiKey,
		BaseURL:   baseURL,
		BaseURLV3: baseURLV3,
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

func (bis *BestInSlot) GetInscriptionsByWalletAddr(addr string, limit, page int64) (*models.BisInscriptions, error) {
	// TODO: Implement limit, note: BIS supports incrementals of 20
	offset := fmt.Sprintf("&offset=%s", strconv.Itoa(int((page-1)*100)))
	count := fmt.Sprintf("&count=%s", strconv.Itoa(int(limit)))
	url := fmt.Sprintf(
		"%s%s%s%s%s",
		"/wallet/inscriptions?address=",
		addr,
		"&sort_by=inscr_num&order=asc",
		count,
		offset,
	)

	resp, err := bis.getV3(url)
	if err != nil {
		return nil, err
	}

	var inscriptions *models.BisInscriptions
	if err := json.Unmarshal(resp, &inscriptions); err != nil {
		return nil, err
	}

	return inscriptions, nil
}

func (bis *BestInSlot) get(endpoint string) ([]byte, error) {
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

func (bis *BestInSlot) getV3(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", bis.BaseURLV3+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", fmt.Sprintf("Bearer %s", bis.APIKey))
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
