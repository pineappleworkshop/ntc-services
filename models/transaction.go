package models

type Tx struct {
	Hex           string  `json:"hex" bson:"hex"`
	TXID          string  `json:"txid" bson:"txid"`
	Hash          string  `json:"hash" bson:"hash"`
	Size          int64   `json:"size" bson:"size"`
	VSize         int64   `json:"vsize" bson:"vsize"`
	Version       int64   `json:"version" bson:"version"`
	LockTime      int64   `json:"lockTime" bson:"lockTime"`
	VIn           []*VIn  `json:"vin" bson:"vin"`
	VOut          []*VOut `json:"vout" bson:"vout"`
	BlockHash     string  `json:"blockhash" bson:"blockhash"`
	Confirmations int64   `json:"confirmations" bson:"confirmations"`
	Time          int64   `json:"time" bson:"time"`
	BlockTime     int64   `json:"blocktime" bson:"blocktime"`
}

type VIn struct {
	Sequence    int64      `json:"sequence" bson:"sequence"`
	Coinbase    *bool      `json:"coinbase" bson:"coinbase"`
	TXID        string     `json:"txid" bson:"txid"`
	VOut        int64      `json:"vout" bson:"vout"`
	ScriptSig   *ScriptSig `json:"script_sig" bson:"script_sig"`
	TXINWitness []string   `json:"txinwitness" bson:"txinwitness"`
}

type ScriptSig struct {
	ASM string `json:"asm" bson:"asm"`
	Hex string `json:"hex" bson:"hex"`
}

type VOut struct {
	Value        float64       `json:"value" bson:"value"`
	N            int64         `json:"n" bson:"n"`
	ScriptPubKey *ScriptPubKey `json:"script_pub_key" bson:"script_pub_key"`
}

type ScriptPubKey struct {
	ASM       string       `json:"asm" bson:"asm"`
	Hex       string       `json:"hex" bson:"hex"`
	ReqSigs   *interface{} `json:"req_sigs" bson:"req_sigs"`
	Type      string       `json:"type" bson:"type"`
	Addresses []string     `json:"addresses" bson:"addresses"`
	Address   string       `json:"address" json:"address"`
}
