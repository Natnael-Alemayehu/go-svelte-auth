package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-secret-key")

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func RegisterUser(db *sql.DB) echo.HandlerFunc {

	return func(c echo.Context) error {
		var req RegisterRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Password encryption failed"})
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		query := `
				INSERT INTO users (username, email, password_hash, created_at)
				VALUES ($1,$2,$3,$4)`
		args := []any{req.Username, req.Email, string(hashedPassword), time.Now()}
		_, err = db.ExecContext(ctx, query, args...)
		if err != nil {
			fmt.Printf("Err: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create a user"})
		}
		return c.JSON(http.StatusCreated, echo.Map{"message": "User created successfully"})
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func LoginUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		var id int
		var passwordHash string

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		query := `
				SELECT id, password_hash FROM users WHERE email=$1`

		err := db.QueryRowContext(ctx, query, req.Email).Scan(&id, &passwordHash)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
		}
		token := jwt.New(jwt.SigningMethodES256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = id
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, err := token.SignedString([]byte("AjFdFok7Ya1H1YBGtjOdg27HqbD7ocuq5yPm0LB6jVM"))
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to generate token"})
		}
		return c.JSON(http.StatusOK, echo.Map{"token": t})
	}
}
