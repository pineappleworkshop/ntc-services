package handlers

import (
	"net/http"
	"ntc-services/config"
	"ntc-services/models"
	"ntc-services/services"

	"github.com/labstack/echo/v4"
)

func HealthHandler(c echo.Context) error {
	// todo: test ordex API wit ha simple http request
	health := new(models.Health)
	health.Service = config.SERVICE_NAME
	health.Status = http.StatusOK
	health.Version = config.VERSION
	health.BestInSlotStatus = http.StatusNotImplemented
	health.OrdexStatus = http.StatusNotImplemented
	inscription, err := services.BESTINSLOT.GetInscriptionById(c, "4e80d14abdb35ce193758cfd69ae8ce67f8036368ac75b729ef2fd3e0c6bad2fi0")
	if err != nil {
		c.Logger().Error(err)
	}
	if inscription != nil {
		health.BestInSlotStatus = http.StatusOK
	}
	ordexInscription, err := services.ORDEX.GetInscriptionById("4e80d14abdb35ce193758cfd69ae8ce67f8036368ac75b729ef2fd3e0c6bad2fi0")
	if err != nil {
		c.Logger().Error(err)
	}
	if ordexInscription != nil {
		health.OrdexStatus = http.StatusOK
	}
	// TODO: path react to env
	// state, err := models.GetState()
	// if err != nil {
	// 	c.Logger().Error(err)
	// }

	// health.State = state

	return c.JSON(http.StatusOK, health)
}
