package models

import "encoding/json"

type BciUTXO struct {
	Confirmations   int64  `json:"confirmations"`
	Script          string `json:"script"`
	TxHash          string `json:"tx_hash"`
	TxHashBigEndian string `json:"tx_hash_big_endian"`
	TxIndex         int64  `json:"tx_index"`
	TxOutputN       int64  `json:"tx_output_n"`
	Value           int64  `json:"value"`
	ValueHex        string `json:"value_hex"`
}

func ParseJSONToBciUTXO(jsonStr string) (*BciUTXO, error) {
	var utxo BciUTXO
	if err := json.Unmarshal([]byte(jsonStr), &utxo); err != nil {
		return nil, err
	}

	return &utxo, nil
}
