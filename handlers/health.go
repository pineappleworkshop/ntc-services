package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ntc-services/config"
	"ntc-services/models"
)

func HealthHandler(c echo.Context) error {
	health := new(models.Health)
	health.Service = config.SERVICE_NAME
	health.Status = http.StatusOK
	health.Version = config.VERSION
	state := models.NewState("")
	if err := state.Read(); err != nil {
	}

	health.State = state

	return c.JSON(http.StatusOK, health)
}
