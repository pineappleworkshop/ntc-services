package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"ntc-services/models"
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

	// TODO: revisit err propagation and ctx tree
	// TODO: try to find inscriptions by inscriptionID before query/possible? (Trying to limit API requests)

	bisInscriptions, err := services.BESTINSLOT.GetInscriptionsByWalletAddr(c, c.Param("addr"), limit, page)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	var inscriptions []*models.Inscription
	for _, bisInscription := range bisInscriptions.Data {
		inscription := models.ParseBISInscription(bisInscription)
		inscriptions = append(inscriptions, inscription)
	}

	resp := models.InscriptionListResp{
		Page:         page,
		Limit:        limit,
		BlockHeight:  bisInscriptions.BlockHeight,
		Inscriptions: inscriptions,
	}

	return c.JSON(http.StatusOK, resp)
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
