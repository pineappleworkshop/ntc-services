package models

type BisBRC20s struct {
	Data        []BisBRC20 `json:"data"`
	BlockHeight int        `json:"block_height"`
}

type BisBRC20 struct {
	Ticker               string  `json:"ticker"`
	OverallBalance       string  `json:"overall_balance"`
	AvailableBalance     string  `json:"available_balance"`
	TransferrableBalance string  `json:"transferrable_balance"`
	ImageURL             *string `json:"image_url"`
}
