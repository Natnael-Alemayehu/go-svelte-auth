package routes

import (
	"database/sql"
	"gosvelte/internal/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, db *sql.DB) {
	e.GET("/api/v1/users", handlers.GetUsers(db))
	e.GET("/api/v1/user/:id", handlers.GetUser(db))

	e.POST("/api/v1/login", handlers.LoginUser(db))
	e.POST("/api/v1/register", handlers.RegisterUser(db))
}
