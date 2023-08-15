package models

type BisInscriptions struct {
	BlockHeight int64            `json:"block_height"`
	Data        []BisInscription `json:"data"`
}

type BisInscriptionsSingle struct {
	BlockHeight int64          `json:"block_height"`
	Data        BisInscription `json:"data"`
}

type BisInscription struct {
	BisURL                  string      `json:"bis_url"`
	CollectionName          interface{} `json:"collection_name"`
	ContentURL              string      `json:"content_url"`
	GenesisHeight           int64       `json:"genesis_height"`
	InscriptionID           string      `json:"inscription_id"`
	InscriptionName         interface{} `json:"inscription_name"`
	InscriptionNumber       int64       `json:"inscription_number"`
	LastTransferBlockHeight int64       `json:"last_transfer_block_height"`
	Metadata                interface{} `json:"metadata"`
	MimeType                string      `json:"mime_type"`
	OwnerWalletAddr         string      `json:"owner_wallet_addr"`
	Satpoint                string      `json:"satpoint"`
	Slug                    interface{} `json:"slug"`
}

//type BISWalletIx struct {
//	IsHTML                  bool        `json:"isHtml"`
//	Validity                interface{} `json:"validity"`
//	IsBRC                   bool        `json:"isBrc"`
//	IsSNS                   bool        `json:"isSns"`
//	Name                    interface{} `json:"name"`
//	Amount                  interface{} `json:"amount"`
//	IsValidTransfer         interface{} `json:"isValidTransfer"`
//	Operation               interface{} `json:"operation"`
//	Ticker                  interface{} `json:"ticker"`
//	IsJSON                  bool        `json:"isJson"`
//	Content                 interface{} `json:"content"` // can be either string or map[string]interface{}
//	InscriptionName         interface{} `json:"inscription_name"`
//	InscriptionID           string      `json:"inscription_id"`
//	InscriptionNumber       int         `json:"inscription_number"`
//	Metadata                interface{} `json:"metadata"`
//	OwnerWalletAddr         string      `json:"owner_wallet_addr"`
//	MimeType                string      `json:"mime_type"`
//	LastSalePrice           interface{} `json:"last_sale_price"`
//	Slug                    interface{} `json:"slug"`
//	CollectionName          interface{} `json:"collection_name"`
//	ContentURL              string      `json:"content_url"`
//	BISURL                  string      `json:"bis_url"`
//	Wallet                  *string     `json:"wallet"`
//	LastTransferBlockHeight int         `json:"last_transfer_block_height"`
//	MediaLength             *int        `json:"media_length"`
//	GenesisTS               *int        `json:"genesis_ts"`
//	GenesisHeight           *int        `json:"genesis_height"`
//	GenesisFee              *int        `json:"genesis_fee"`
//	OutputValue             *int        `json:"output_value"`
//	SATPoint                *string     `json:"satpoint"`
//	CollectionSlug          *string     `json:"collection_slug"`
//	Confirmations           *int        `json:"confirmations"`
//}
