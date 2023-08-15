package services

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ntc-services/models"
	"testing"
	"time"
)

var now = time.Now().Unix()

var makerWallet = &models.Wallet{
	ID:                 primitive.NewObjectID(),
	Type:               "unisat",
	CardinalAddr:       "",
	TapRootAddr:        "bc1p0t40pgryukh88rhwx4ffzt28cjmhxnpm56s3382vyy22zek5wpmq8rps2l",
	SegwitAddr:         "",
	CreatedAt:          now,
	UpdatedAt:          &now,
	LastConnectedAt:    &now,
	LastConnectedBlock: &now,
}

var makerUTXOsAll = []*models.UTXO{
	{
		Confirmations:   2675,
		Script:          "51207aeaf0a064e5ae738eee3552912d47c4b7734c3ba6a1189d4c2114a166d47076",
		TxHash:          "9fb77d422cad2a8ffcb9ed5f8270d7d1050569e222ae3c48c0f89d2a1943641e",
		TxHashBigEndian: "1e6443192a9df8c0483cae22e2690505d1d770825fedb9fc8f2aad2c427db79f",
		TxIndex:         5619533738251685,
		TxOutputN:       0,
		Value:           10000,
		ValueHex:        "2710",
	},
	{
		Confirmations:   3776,
		Script:          "51207aeaf0a064e5ae738eee3552912d47c4b7734c3ba6a1189d4c2114a166d47076",
		TxHash:          "fc5b52cff7b78fbade0c64a8046ff812f44e03d663ff97b0f4b20bd7f28e1ed6",
		TxHashBigEndian: "d61e8ef2d70bb2f4b097ff63d6034ef412f86f04a8640cdeba8fb7f7cf525bfc",
		TxIndex:         8879013170706161,
		TxOutputN:       0,
		Value:           546,
		ValueHex:        "0222",
	},
}

var floorPrice = int64(10000000)
var makerInscriptions = []*models.Inscription{
	{
		InscriptionName:         nil,
		InscriptionID:           "d61e8ef2d70bb2f4b097ff63d6034ef412f86f04a8640cdeba8fb7f7cf525bfci0",
		InscriptionNumber:       18284663,
		Metadata:                nil,
		OwnerWalletAddr:         "bc1p0t40pgryukh88rhwx4ffzt28cjmhxnpm56s3382vyy22zek5wpmq8rps2l",
		MimeType:                "image/png",
		LastSalePrice:           nil,
		Slug:                    nil,
		CollectionName:          nil,
		Satpoint:                "d61e8ef2d70bb2f4b097ff63d6034ef412f86f04a8640cdeba8fb7f7cf525bfc:0:0",
		LastTransferBlockHeight: 799444,
		ContentURL:              "https://bis-ord-content.fra1.cdn.digitaloceanspaces.com/ordinals/d61e8ef2d70bb2f4b097ff63d6034ef412f86f04a8640cdeba8fb7f7cf525bfci0",
		BisURL:                  "https://bestinslot.xyz/ordinals/inscription/d61e8ef2d70bb2f4b097ff63d6034ef412f86f04a8640cdeba8fb7f7cf525bfci0",
		FloorPrice:              &floorPrice,
		Confirmations:           nil,
	},
}

var makerInscriptionsAll = []*models.Inscription{
	{
		InscriptionName:         nil,
		InscriptionID:           "d61e8ef2d70bb2f4b097ff63d6034ef412f86f04a8640cdeba8fb7f7cf525bfci0",
		InscriptionNumber:       18284663,
		Metadata:                nil,
		OwnerWalletAddr:         "bc1p0t40pgryukh88rhwx4ffzt28cjmhxnpm56s3382vyy22zek5wpmq8rps2l",
		MimeType:                "image/png",
		LastSalePrice:           nil,
		Slug:                    nil,
		CollectionName:          nil,
		Satpoint:                "d61e8ef2d70bb2f4b097ff63d6034ef412f86f04a8640cdeba8fb7f7cf525bfc:0:0",
		LastTransferBlockHeight: 799444,
		ContentURL:              "https://bis-ord-content.fra1.cdn.digitaloceanspaces.com/ordinals/d61e8ef2d70bb2f4b097ff63d6034ef412f86f04a8640cdeba8fb7f7cf525bfci0",
		BisURL:                  "https://bestinslot.xyz/ordinals/inscription/d61e8ef2d70bb2f4b097ff63d6034ef412f86f04a8640cdeba8fb7f7cf525bfci0",
		FloorPrice:              &floorPrice,
		Confirmations:           nil,
	},
}

var maker = &models.Side{
	ID:                 primitive.NewObjectID(),
	WalletID:           makerWallet.ID,
	Wallet:             makerWallet,
	InscriptionNumbers: nil,
	BTC:                0,
	Inscriptions:       makerInscriptions,
	//UTXOs:              makerUTXOs,
	CreatedAt: now,
	UpdatedAt: &now,
}

var takerWallet = &models.Wallet{
	ID:                 primitive.NewObjectID(),
	Type:               "unisat",
	CardinalAddr:       "",
	TapRootAddr:        "bc1pxy8gsmgu5zzv0ytj7ae4pgnqkcdwaqas7xmc4szcg70zqwsj4rxss2tvhu",
	SegwitAddr:         "",
	CreatedAt:          now,
	UpdatedAt:          &now,
	LastConnectedAt:    &now,
	LastConnectedBlock: &now,
}

var takerUTXOsAll = []*models.UTXO{
	{
		Confirmations:   2968,
		Script:          "5120310e886d1ca084c79172f77350a260b61aee83b0f1b78ac058479e203a12a8cd",
		TxHash:          "381b5d5ea418f44183a1971798e7fbfa6f3d1d6fd21852f052d762a258e58f1a",
		TxHashBigEndian: "1a8fe558a262d752f05218d26f1d3d6ffafbe7981797a18341f418a45e5d1b38",
		TxIndex:         1974085816189726,
		TxOutputN:       0,
		Value:           300000,
		ValueHex:        "0493e0",
	},
}

var takerInscriptions = []*models.Inscription{}

var takerInscriptionsAll = []*models.Inscription{}

var taker = &models.Side{
	ID:                 primitive.NewObjectID(),
	WalletID:           takerWallet.ID,
	Wallet:             takerWallet,
	InscriptionNumbers: nil,
	BTC:                100000,
	Inscriptions:       nil,
	//UTXOs:              takerUTXOs,
	CreatedAt: now,
	UpdatedAt: &now,
}

var trade = &models.Trade{
	ID:              primitive.NewObjectID(),
	MakerID:         maker.ID,
	Maker:           maker,
	TakerID:         &taker.ID,
	Taker:           taker,
	PSBT:            nil,
	FeeRate:         nil,
	MaxFee:          10,
	PlatformFee:     nil,
	TxID:            nil,
	Status:          "OFFER_ACCEPTED", // something else
	StatusChangedAt: &now,
	CreatedAt:       now,
	UpdatedAt:       &now,
}

func TestPBST(t *testing.T) {
	var err error
	psbt := NewPBST(trade, makerUTXOsAll, takerUTXOsAll, makerInscriptionsAll, takerInscriptionsAll)
	assert.Len(t, psbt.MakerUTXOsAll, 2)
	assert.Len(t, psbt.TakerUTXOsAll, 1)
	assert.Len(t, psbt.MakerInscriptionsAll, 1)
	assert.Len(t, psbt.TakerInscriptionsAll, 0)

	err = psbt.selectInscriptionsUTXOs()
	assert.Nil(t, err)
	assert.Len(t, psbt.MakerInscriptionUTXOs, 1)
	assert.Len(t, psbt.MakerPaymentUTXOs, 0)
	assert.Len(t, psbt.MakerUTXOsAll, 1)
	assert.Len(t, psbt.TakerInscriptionUTXOs, 0)
	assert.Len(t, psbt.TakerPaymentUTXOs, 0)
	assert.Len(t, psbt.TakerUTXOsAll, 1)

	err = psbt.calculatePlatformFee()
	assert.Nil(t, err)
	assert.Equal(t, int64(101000), psbt.PlatformFee)
}
