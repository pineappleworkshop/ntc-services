package stores

const (
	DB_NAME = "ntc_protocol"

	MONGO_ERR_NOT_FOUND = "mongo: no documents in result"
	ID_HEX_INVALID_ERR  = "the provided hex string is not a valid ObjectID"

	DB_COLLECTION_BLOCKS_RAW = "blocks_raw"
	DB_COLLECTION_TXS_RAW    = "txs_raw"
	DB_COLLECTION_TRADES     = "trades"
	DB_COLLECTION_STATE      = "state"
	DB_COLLECTION_TXS        = "txs"
	DB_COLLECTION_VINS       = "vins"
	DB_COLLECTION_VOUTS      = "vouts"
)
