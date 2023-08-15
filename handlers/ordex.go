package handlers

import (
	"fmt"
	"net/http"
	"ntc-services/services"

	"github.com/labstack/echo/v4"
)

type ReqBody struct {
	InscriptionNumbers []string `json:"inscription_numbers" bson:"inscription_numbers"`
}

func NewReqBody() *ReqBody {

	return &ReqBody{}
}

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

func OrdexGetInscriptionsByIds(c echo.Context) error {
	ReqBody := NewReqBody()
	if err := c.Bind(ReqBody); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	response, err := services.ORDEX.GetInscriptionsByIds(ReqBody.InscriptionNumbers)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, response)
}
