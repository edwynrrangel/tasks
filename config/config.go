package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/edwynrrangel/tasks/logger"
)

type (
	Configuration struct {
		Port               string
		CorsAllowedOrigins string
		PostgreSQLClient   *sqlx.DB
		JWTExpireTimeInMin uint
		JWTSecret          string
	}
)

func LoadConfig() *Configuration {
	loadEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	corsAllowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsAllowedOrigins == "" {
		corsAllowedOrigins = "*"
	}

	clientDB, err := getPostgresqlConfig()
	if err != nil {
		logger.Fatal(err)
	}

	jwtExpireTimeInMin, err := strconv.ParseUint(os.Getenv("JWT_EXPIRE_TIME_IN_MIN"), 10, 32)
	if err != nil {
		jwtExpireTimeInMin = 15
	}

	return &Configuration{
		Port:               port,
		CorsAllowedOrigins: corsAllowedOrigins,
		PostgreSQLClient:   clientDB,
		JWTExpireTimeInMin: uint(jwtExpireTimeInMin),
		JWTSecret:          os.Getenv("JWT_SECRET"),
	}
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		logger.Warn("error loading .env file")
	}
}

func getPostgresqlConfig() (*sqlx.DB, error) {
	client, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		),
	)
	if err != nil {
		return nil, err
	}

	err = client.Ping()
	if err != nil {
		return nil, err
	}

	logger.Info("Connected to PostgreSQL!")
	return client, nil
}

var Config *Configuration
