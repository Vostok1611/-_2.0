package main

import (
	"GOBAR/internal/db"
	"GOBAR/internal/handlers"
	"GOBAR/internal/taskService"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {

	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("could not connect to BD: %v", err)
	}
	e := echo.New()

	repo := taskService.NewTaskRepository(database)
	service := taskService.NewTaskService(repo)
	handler := handlers.NewTaskHandler(service)

	e.GET("/task", handler.GetTask)
	e.POST("/task", handler.PostTask)
	e.PATCH("/task/:id", handler.PatchTask)
	e.DELETE("/task/:id", handler.DeleteTask)

	e.Logger.Fatal(e.Start(":8080"))
}
