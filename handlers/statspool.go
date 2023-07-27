package handlers

import (
	"net/http"
	"ntc-services/models"
	"ntc-services/services"

	"github.com/labstack/echo/v4"
)

func GetStatsPool(c echo.Context) error {
	price, err := services.BLOCKCHAIN.GetBTCPrice()
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

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

	resp := models.NewStatsPool(price, blockHeight)
	if err := resp.Parse(recommendedFees.(map[string]interface{})); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, resp)
}
