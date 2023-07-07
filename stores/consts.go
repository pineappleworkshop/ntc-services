package stores

const (
	DB_NAME = "omnisat-mongo-indexing"

	MONGO_ERR_NOT_FOUND = "mongo: no documents in result"
	ID_HEX_INVALID_ERR  = "the provided hex string is not a valid ObjectID"

	DB_COLLECTION_BRCS = "brcs"

	DB_COLLECTION_TICKERS   = "brc20_tickers"
	DB_COLLECTION_DEPLOYS   = "brc20_deploys"
	DB_COLLECTION_MINTS     = "brc20_mints"
	DB_COLLECTION_TRANSFERS = "brc20_transfers"
	DB_COLLECTION_BALANCES  = "brc20_user_balances"

	DB_COLLECTION_BRCS20_INVALIDS = "brc20_invalids"
)
