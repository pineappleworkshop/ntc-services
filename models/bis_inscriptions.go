package models

type BisInscriptions struct {
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
	Wallet                  string      `json:"wallet"`
	MediaLength             int64       `json:"media_length"`
	GenesisTimestamp        int64       `json:"genesis_ts"`
	GenesisFee              int64       `json:"genesis_fee"`
	OutputValue             int64       `json:"output_value"`
	LastSalePrice           int64       `json:"last_sale_price"`
	CollectionSlug          interface{} `json:"collection_slug"`
	CollectionFloorPrice    int64       `json:"collection_floor_price"`
	MinPrice                int64       `json:"min_price"`
	OrdswapPrice            int64       `json:"ordswap_price"`
	MagicedenPrice          int64       `json:"magiceden_price"`
	OrdinalsWalletPrice     int64       `json:"ordinalswallet_price"`
	GammaioPrice            int64       `json:"gammaio_price"`
	NostrPrice              int64       `json:"nostr_price"`
	OdynalsPrice            int64       `json:"odynals_price"`
	UnisatPrice             int64       `json:"unisat_price"`
	OrdinalsMarketPrice     int64       `json:"ordinalsmarket_price"`
	OkxPrice                int64       `json:"okx_price"`
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
