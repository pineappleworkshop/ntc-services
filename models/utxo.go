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
}
