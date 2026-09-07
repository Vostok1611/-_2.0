package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var task string

type requestBody struct {
	Task string `json:"task"`
}

func postHandler(c echo.Context) error {
	var body requestBody

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	task = body.Task
	return c.JSON(http.StatusOK, map[string]string{"message": "task saved"})
}

func getHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "hello, "+task)
}

func main() {
	e := echo.New()

	e.GET("/task", getHandler)
	e.POST("/task", postHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
