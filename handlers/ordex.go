package handlers

import (
	"fmt"
	"net/http"
	"ntc-services/services"

	"github.com/labstack/echo/v4"
)

func OrdexHandler(c echo.Context) error {
	InscriptionId := c.Param("id")
	ordexInscription, err := services.ORDEX.GetInscriptionById(InscriptionId)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	fmt.Println("---------ORDEX INSCRIPTION: ", ordexInscription)

	return c.JSON(http.StatusOK, ordexInscription)
}
