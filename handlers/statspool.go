package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ntc-services/models"
	"ntc-services/services"
)

func GetStatsPool(c echo.Context) error {
	price, err := services.BESTINSLOT.GetBTCPrice()
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	recommendedFees, err := services.MEMPOOL.GetRecommendedFees()
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp := models.NewStatsPool(price)
	if err := resp.Parse(recommendedFees.(map[string]interface{})); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, resp)
}
