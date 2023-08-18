package models

type UTXO struct {
	TxHashBigEndian string  `json:"tx_hash_big_endian" bson:"tx_hash_big_endian"`
	TxHash          string  `json:"tx_hash" bson:"tx_hash"`
	TxOutputN       int64   `json:"tx_output_n" bson:"tx_output_n"`
	Script          string  `json:"script" bson:"script"`
	Value           int64   `json:"value" bson:"value"`
	ValueHex        string  `json:"value_hex" bson:"value_hex"`
	Confirmations   int64   `json:"confirmations" bson:"confirmations"`
	TxIndex         float64 `json:"tx_index" bson:"tx_index"`
	IsInscription   bool    `json:"is_inscription" bson:"is_inscription"`
}

func (utxo *UTXO) Parse(utxoM map[string]interface{}) error {
	if utxoM["tx_hash_big_endian"] != nil {
		utxo.TxHashBigEndian = utxoM["tx_hash_big_endian"].(string)
	}
	if utxoM["tx_hash"] != nil {
		utxo.TxHash = utxoM["tx_hash"].(string)
	}
	if utxoM["tx_output_n"] != nil {
		utxo.TxOutputN = int64(utxoM["tx_output_n"].(float64))
	}
	if utxoM["script"] != nil {
		utxo.Script = utxoM["script"].(string)
	}
	if utxoM["value"] != nil {
		utxo.Value = int64(utxoM["value"].(float64))
	}
	if utxoM["value_hex"] != nil {
		utxo.ValueHex = utxoM["value_hex"].(string)
	}
	if utxoM["confirmations"] != nil {
		utxo.Confirmations = int64(utxoM["confirmations"].(float64))
	}
	if utxoM["tx_index"] != nil {
		utxo.TxIndex = utxoM["tx_index"].(float64)
	}

	return nil
}
