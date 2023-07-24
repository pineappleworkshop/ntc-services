package models

import "fmt"

type StatsPool struct {
	FastestFee  int64 `json:"fastestFee"`
	HalfHourFee int64 `json:"halfHourFee"`
	HourFee     int64 `json:"hourFee"`
	EconomyFee  int64 `json:"economyFee"`
	MinimumFee  int64 `json:"minimumFee"`
	Computed    struct {
		FastestFee  string `json:"fastestFee"`
		HalfHourFee string `json:"halfHourFee"`
		HourFee     string `json:"hourFee"`
		EconomyFee  string `json:"economyFee"`
		MinimumFee  string `json:"minimumFee"`
	} `json:"computed"`
	BTCRaw float64 `json:"btcRaw"`
	BTC    string  `json:"btc"`
}

func NewStatsPool(btcRaw float64) *StatsPool {
	return &StatsPool{
		Computed: struct {
			FastestFee  string `json:"fastestFee"`
			HalfHourFee string `json:"halfHourFee"`
			HourFee     string `json:"hourFee"`
			EconomyFee  string `json:"economyFee"`
			MinimumFee  string `json:"minimumFee"`
		}{},
		BTCRaw: btcRaw,
	}
}

func (sp *StatsPool) Parse(fees map[string]interface{}) error {
	if fees["economyFee"] != nil {
		sp.EconomyFee = int64(fees["economyFee"].(float64))
		computed, err := sp.compute(fees["economyFee"].(float64))
		if err != nil {
			return err
		}
		sp.Computed.EconomyFee = fmt.Sprintf("%v", computed)
	}
	if fees["fastestFee"] != nil {
		sp.FastestFee = int64(fees["fastestFee"].(float64))
		computed, err := sp.compute(fees["fastestFee"].(float64))
		if err != nil {
			return err
		}
		sp.Computed.FastestFee = fmt.Sprintf("%v", computed)
	}
	if fees["halfHourFee"] != nil {
		sp.HalfHourFee = int64(fees["halfHourFee"].(float64))
		computed, err := sp.compute(fees["halfHourFee"].(float64))
		if err != nil {
			return err
		}
		sp.Computed.HalfHourFee = fmt.Sprintf("%v", computed)
	}
	if fees["hourFee"] != nil {
		sp.HourFee = int64(fees["hourFee"].(float64))
		computed, err := sp.compute(fees["hourFee"].(float64))
		if err != nil {
			return err
		}
		sp.Computed.HourFee = fmt.Sprintf("%v", computed)
	}
	if fees["minimumFee"] != nil {
		sp.MinimumFee = int64(fees["minimumFee"].(float64))
		computed, err := sp.compute(fees["minimumFee"].(float64))
		if err != nil {
			return err
		}
		sp.Computed.MinimumFee = fmt.Sprintf("%v", computed)
	}
	sp.BTC = fmt.Sprintf("%v", sp.BTCRaw)

	return nil
}

func (sp *StatsPool) compute(num float64) (float64, error) {
	computed := num * 140 * (sp.BTCRaw / 100000000.0)

	return computed, nil
}
