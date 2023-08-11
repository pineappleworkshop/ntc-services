package models

type InscriptionListResp struct {
	Page         int64          `json:"page"`
	Limit        int64          `json:"limit"`
	BlockHeight  int64          `json:"blockHeight"`
	Inscriptions []*Inscription `json:"inscriptions"`
}

type Inscription struct {
	InscriptionName         *string `json:"inscription_name" bson:"inscription_name"`
	InscriptionID           string  `json:"inscription_id" bson:"inscription_id"`
	InscriptionNumber       int64   `json:"inscription_number" bson:"inscription_number"`
	Metadata                *string `json:"metadata" bson:"metadata"`
	OwnerWalletAddr         string  `json:"owner_wallet_addr" bson:"owner_wallet_addr"`
	MimeType                string  `json:"mime_type" bson:"mime_type"`
	LastSalePrice           *int64  `json:"last_sale_price" bson:"last_sale_price"`
	Slug                    *string `json:"slug" bson:"slug"`
	CollectionName          *string `json:"collection_name" bson:"collection_name"`
	Satpoint                string  `json:"satpoint" bson:"satpoint"`
	LastTransferBlockHeight int64   `json:"last_transfer_block_height" bson:"last_transfer_block_height"`
	ContentURL              string  `json:"content_url" bson:"content_url"`
	BisURL                  string  `json:"bis_url" bson:"bis_url"`
	FloorPrice              *int64  `json:"floor_price,omitempty" bson:"floor_price"`
	Confirmations           *int64  `json:"confirmations,omitempty" bson:"confirmations"`
}

func ParseBISInscription(b BisInscription) *Inscription {
	convertInterfaceToStringPointer := func(v interface{}) *string {
		if s, ok := v.(string); ok {
			return &s
		}
		return nil
	}

	return &Inscription{
		InscriptionName:         convertInterfaceToStringPointer(b.InscriptionName),
		InscriptionID:           b.InscriptionID,
		InscriptionNumber:       b.InscriptionNumber,
		Metadata:                convertInterfaceToStringPointer(b.Metadata),
		OwnerWalletAddr:         b.OwnerWalletAddr,
		MimeType:                b.MimeType,
		Slug:                    convertInterfaceToStringPointer(b.Slug),
		CollectionName:          convertInterfaceToStringPointer(b.CollectionName),
		Satpoint:                b.Satpoint,
		LastTransferBlockHeight: b.LastTransferBlockHeight,
		ContentURL:              b.ContentURL,
		BisURL:                  b.BisURL,
	}
}
