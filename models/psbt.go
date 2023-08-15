package models

type PSBT struct {
	MakerPayment     int64           `json:"maker_payment" bson:"maker_payment"`
	TakerPayment     int64           `json:"taker_payment" bson:"taker_payment"`
	MakerChange      int64           `json:"maker_change" bson:"maker_change"`
	TakerChange      int64           `json:"taker_change" bson:"taker_change"`
	PlatformFee      int64           `json:"platform_fee" bson:"platform_fee"`
	MinerFee         int64           `json:"miner_fee" bson:"miner_fee"`
	Inputs           map[int]*Input  `json:"inputs" bson:"inputs"`
	Outputs          map[int]*Output `json:"outputs" bson:"outputs"`
	PreMinerFeePSBT  *PSBTSerialized `json:"pre_miner_fee_psbt" bson:"pre_miner_fee_psbt"`
	PostMinerFeePSBT *PSBTSerialized `json:"post_miner_fee_psbt" bson:"post_miner_fee_psbt"`
	FinalizedPSBT    *PSBTSerialized `json:"finalized_psbt" bson:"finalized_psbt"`
}

type PSBTSerialized struct {
	//Bytes  []byte `json:"bytes" bson:"bytes"`
	Base64 string `json:"base64" bson:"base64"`
}

type Input struct {
	Addr        string `json:"addr" bson:"addr"`
	PublicKey   string `json:"public_key" bson:"public_key"`
	Type        string `json:"type" bson:"type"`
	TxID        string `json:"tx_id" bson:"tx_id"`
	Index       int64  `json:"index" bson:"index"`
	WitnessUTXO struct {
		Script string `json:"script" bson:"script"`
		Amount int64  `json:"amount" bson:"amount"`
	} `json:"witness_utxo" bson:"witness_utxo"`
}

type Output struct {
	Addr          string `json:"addr" bson:"addr"`
	Value         int64  `json:"value" bson:"value"`
	IsInscription bool   `json:"is_inscription" bson:"is_inscription"`
	IsPayment     bool   `json:"is_payment" bson:"is_payment"`
	IsChange      bool   `json:"is_change" bson:"is_change"`
	IsFee         bool   `json:"is_fee" bson:"is_fee"`
}
