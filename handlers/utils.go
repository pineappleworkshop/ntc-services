package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func parsePagination(c echo.Context) (int64, error) {
	var page int
	if c.QueryParam("page") != "" {
		var err error
		page, err = strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			c.Logger().Error(http.StatusInternalServerError, "Something went wrong")
		}
	} else {
		page = 1
	}

	//var limit int
	//if c.QueryParam("limit") != "" {
	//	var err error
	//	limit, err = strconv.Atoi(c.QueryParam("limit"))
	//	if err != nil {
	//		c.Logger().Error(http.StatusInternalServerError, "Something went wrong")
	//	}
	//} else {
	//	limit = 100
	//}

	return int64(page), nil
}
