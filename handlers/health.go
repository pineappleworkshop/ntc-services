package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"omnisat-api/config"
	"omnisat-api/models"
)

func HealthHandler(c echo.Context) error {
	health := new(models.Health)
	health.Service = config.SERVICE_NAME
	health.Status = http.StatusOK
	health.Version = config.VERSION

	return c.JSON(http.StatusOK, health)
}
