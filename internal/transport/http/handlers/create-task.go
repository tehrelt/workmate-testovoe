package handlers

import "github.com/labstack/echo/v4"

func CreateTask() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, echo.Map{
			"message": "task created",
		})
	}
}
