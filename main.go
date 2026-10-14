package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=gobar port=5433 sslmode=disable"
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Task{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
}

type Task struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Task      string         `json:"task"`
	IsDone    bool           `json:"is_done"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type requestBody struct {
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
}

func postHandler(c echo.Context) error {
	var body requestBody

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if body.Task == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Task cannot be empty"})
	}

	newTask := Task{
		Task:   body.Task,
		IsDone: body.IsDone,
	}

	if err := db.Create(&newTask).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not add task"})
	}

	return c.JSON(http.StatusCreated, newTask)
}

func getHandler(c echo.Context) error {
	var tasks []Task

	if err := db.Find(&tasks).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get tasks"})
	}
	return c.JSON(http.StatusOK, tasks)
}

func patchHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	var body requestBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if body.Task == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Task cannot be empty"})
	}

	var task Task
	if err := db.First(&task, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
	}

	task.Task = body.Task
	task.IsDone = body.IsDone

	if err := db.Save(&task).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update task"})
	}

	return c.JSON(http.StatusOK, task)
}

func deleteHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	var task Task
	if err := db.First(&task, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
	}

	if err := db.Delete(&task).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete task"})
	}

	return c.NoContent(http.StatusNoContent)
}

func main() {
	initDB()
	e := echo.New()

	e.GET("/task", getHandler)
	e.POST("/task", postHandler)
	e.PATCH("/task/:id", patchHandler)
	e.DELETE("/task/:id", deleteHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
