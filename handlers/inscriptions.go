package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ntc-services/services"
)

func GetInscriptions(c echo.Context) error {
	page, err := parsePagination(c)
	if err != nil {
		// TODO: handle proper err
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// TODO: Validate
	addr := c.Param("addr")

	inscriptions, err := services.BESTINSLOT.GetInscriptionsByWalletAddr(addr, page)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, inscriptions)
}
