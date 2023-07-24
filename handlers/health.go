package handlers

import (
	"net/http"
	"ntc-services/config"
	"ntc-services/models"

	"github.com/labstack/echo/v4"
)

func HealthHandler(c echo.Context) error {
	health := new(models.Health)
	health.Service = config.SERVICE_NAME
	health.Status = http.StatusOK
	health.Version = config.VERSION

	// TODO: path react to env
	// state, err := models.GetState()
	// if err != nil {
	// 	c.Logger().Error(err)
	// }

	// health.State = state

	return c.JSON(http.StatusOK, health)
}
