package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ntc-services/models"
	"ntc-services/services"
)

func GetStatsPool(c echo.Context) error {
	// TODO: use another source to query for current BTC price
	// https://www.blockchain.com/explorer/api/blockchain_api

	//price, err := services.BESTINSLOT.GetBTCPrice()
	//if err != nil {
	//	c.Logger().Error(err)
	//	return c.JSON(http.StatusInternalServerError, err)
	//}

	blockHeight, err := services.MEMPOOL.GetBlockHeight()
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	recommendedFees, err := services.MEMPOOL.GetRecommendedFees()
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp := models.NewStatsPool(-1.0, blockHeight)
	if err := resp.Parse(recommendedFees.(map[string]interface{})); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, resp)
}
