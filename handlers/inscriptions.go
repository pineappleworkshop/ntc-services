package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"ntc-services/services"
)

func GetInscriptions(c echo.Context) error {
	page, limit, err := parsePagination(c)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	if valid := validateBTCAddress(c.Param("addr")); !valid {
		err := errors.New("BTC address is not valid")
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	// TODO: revisit err propogation and ctx tree
	// TODO: discuss the resp structure.. currently a clone of BIS
	inscriptions, err := services.BESTINSLOT.GetInscriptionsByWalletAddr(c, c.Param("addr"), limit, page)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, inscriptions)
}

func GetBRC20s(c echo.Context) error {
	if valid := validateBTCAddress(c.Param("addr")); !valid {
		err := errors.New("BTC address is not valid")
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	brc20s, err := services.BESTINSLOT.GetBRC20sByWalletAddr(c, c.Param("addr"), 0, 0)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, brc20s)
}
