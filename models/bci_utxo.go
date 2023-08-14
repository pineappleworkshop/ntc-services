package models

type BciUTXO struct {
	Confirmations   int64  `json:"confirmations"`
	Script          string `json:"script"`
	TxHash          string `json:"tx_hash"`
	TxHashBigEndian string `json:"tx_hash_big_endian"`
	TxIndex         int64  `json:"tx_index"`
	TxOutputN       uint32 `json:"tx_output_n"`
	Value           int64  `json:"value"`
	ValueHex        string `json:"value_hex"`
}

func ParseJSONToBciUTXO(BciUTXOM map[string]interface{}) (*BciUTXO, error) {
	utxo := new(BciUTXO)
	for k, v := range BciUTXOM {
		if k == "confirmations" {
			utxo.Confirmations = int64(v.(float64))
		}
		if k == "script" {
			utxo.Script = v.(string)
		}
		if k == "tx_hash" {
			utxo.TxHash = v.(string)
		}
		if k == "tx_hash_big_endian" {
			utxo.TxHashBigEndian = v.(string)
		}
		if k == "tx_index" {
			utxo.TxIndex = int64(v.(float64))
		}
		if k == "tx_output_n" {
			utxo.TxOutputN = uint32(v.(float64))
		}
		if k == "value" {
			utxo.Value = int64(v.(float64))
		}
		if k == "value_hex" {
			utxo.ValueHex = v.(string)
		}
	}

	return utxo, nil
}
