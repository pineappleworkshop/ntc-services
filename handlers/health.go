package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ntc-services/config"
	"ntc-services/models"
	"ntc-services/services"
)

func HealthHandler(c echo.Context) error {
	health := new(models.Health)
	health.Service = config.SERVICE_NAME
	health.Status = http.StatusOK
	health.Version = config.VERSION

	// TODO: path react to env
	state := models.NewState(services.STATE_PATH)
	if err := state.Read(); err != nil {
		c.Logger().Error(err)
	}

	health.State = state

	return c.JSON(http.StatusOK, health)
}
