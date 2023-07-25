package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ntc-services/services"
)

func GetInscriptions(c echo.Context) error {
	// TODO: Validate
	addr := c.Param("addr")

	inscriptions, err := services.BESTINSLOT.GetInscriptionsByWalletAddr(addr)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, inscriptions)
}
