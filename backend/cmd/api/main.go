package main

import (
	"flag"
	"gosvelte/config"
	"gosvelte/internal/db"
	"gosvelte/internal/routes"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	var cfg config.Config

	flag.StringVar(&cfg.Port, "port", ":8080", "API server port. Example :8080")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development | staging | production)")

	// Setting the DSN for the database
	flag.StringVar(&cfg.DB.DSN, "db-dsn", "", "PostgreSQL DSN")

	// Setting connection pool settings from the command line.
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns8080", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.Parse()

	// Initialize database connection
	dbConn, err := db.OpenDB(cfg)
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}
	defer dbConn.Close()

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register routes from routes package
	routes.RegisterRoutes(e, dbConn)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))

}
