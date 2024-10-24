package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetUsers(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		query := `SELECT id, username, email, created_at FROM users`
		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch users"})
		}
		defer rows.Close()
		users := []map[string]interface{}{}
		for rows.Next() {
			var id int
			var username, email string
			var created_at time.Time

			if err := rows.Scan(&id, &username, &email, &created_at); err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "faild to scan users"})
			}
			user := map[string]interface{}{
				"id":         id,
				"username":   username,
				"email":      email,
				"created_at": created_at,
			}
			users = append(users, user)
		}
		if len(users) == 0 {
			return c.JSON(http.StatusOK, echo.Map{"message": "No users found"})
		}
		return c.JSON(http.StatusOK, users)
	}
}

func GetUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		query := `SELECT username, email, created_at FROM users WHERE id=$1`

		var username, email string
		var created_at time.Time

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := db.QueryRowContext(ctx, query, id).Scan(&username, &email, &created_at)
		if err != nil {
			fmt.Printf("err: %v", err)
			return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
		}
		user := map[string]interface{}{
			"usename":    username,
			"email":      email,
			"created_at": created_at,
		}
		return c.JSON(http.StatusOK, user)
	}
}
