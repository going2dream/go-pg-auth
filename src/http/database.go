package http

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type (
	DatabaseConfig struct {
		Host           string
		Port           string
		Database       string
		Username       string
		Password       string
		UsersTable     string
		UsernameColumn string
		PasswordColumn string
	}
)

func NewDBConnection() *pgxpool.Pool {
	configFile, err := ioutil.ReadFile("config/database.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config DatabaseConfig
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("error: %v", err)
	}

	databaseURL := "postgresql://" + config.Username + ":" + config.Password + "@" + config.Host + ":" + config.Port + "/" + config.Database

	conn, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	return conn
}
