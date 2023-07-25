package models

type StatsPool struct {
	FastestFee  int64   `json:"fastestFee"`
	HalfHourFee int64   `json:"halfHourFee"`
	HourFee     int64   `json:"hourFee"`
	EconomyFee  int64   `json:"economyFee"`
	MinimumFee  int64   `json:"minimumFee"`
	BTCRaw      float64 `json:"btcRaw"`
	BTC         string  `json:"btc"`
	BlockHeight int64   `json:"blockHeight"`
}

func NewStatsPool(btcRaw float64, blockHeight int64) *StatsPool {
	return &StatsPool{
		BTCRaw:      btcRaw,
		BlockHeight: blockHeight,
	}
}

func (sp *StatsPool) Parse(fees map[string]interface{}) error {
	if fees["economyFee"] != nil {
		sp.EconomyFee = int64(fees["economyFee"].(float64))
	}
	if fees["fastestFee"] != nil {
		sp.FastestFee = int64(fees["fastestFee"].(float64))
	}
	if fees["halfHourFee"] != nil {
		sp.HalfHourFee = int64(fees["halfHourFee"].(float64))
	}
	if fees["hourFee"] != nil {
		sp.HourFee = int64(fees["hourFee"].(float64))
	}
	if fees["minimumFee"] != nil {
		sp.MinimumFee = int64(fees["minimumFee"].(float64))
	}

	return nil
}
