package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func parsePagination(c echo.Context) (int64, int64, error) {
	var page int
	if c.QueryParam("page") != "" {
		var err error
		page, err = strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			c.Logger().Error(http.StatusInternalServerError, "Paginating page failed")
		}
	} else {
		page = 1
	}
	if page < 1 {
		err := errors.New("pagination page cannot be less than 1")
		c.Logger().Error(http.StatusInternalServerError, err)
		return -1, -1, err
	}

	var limit int
	if c.QueryParam("limit") != "" {
		var err error
		limit, err = strconv.Atoi(c.QueryParam("limit"))
		if err != nil {
			c.Logger().Error(http.StatusInternalServerError, "Paginating limit failed")
		}
	} else {
		limit = 100
	}

	if limit > 100 {
		err := errors.New("pagination limit cannot be greater the 100")
		c.Logger().Error(http.StatusInternalServerError, err)
		return -1, -1, err
	}
	if limit%20 != 0 {
		err := errors.New("pagination limit must be in increments of 20")
		c.Logger().Error(http.StatusInternalServerError, err)
		return -1, -1, err
	}

	return int64(page), int64(limit), nil
}
