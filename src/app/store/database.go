package store

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
		DBHost         string `yaml:"db_host"`
		DBPort         string `yaml:"db_port"`
		DBUsername     string `yaml:"db_username"`
		DBPassword     string `yaml:"db_password"`
		Database       string `yaml:"database"`
		UsersTable     string `yaml:"users_table"`
		LoginColumn    string `yaml:"login_column"`
		PasswordColumn string `yaml:"password_column"`
	}
)

func NewPoolInstance() *pgxpool.Pool {
	configFile, err := ioutil.ReadFile("config/database.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config DatabaseConfig
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("error: %v", err)
	}

	databaseURL := "postgresql://" + config.DBUsername + ":" + config.DBPassword + "@" + config.DBHost + ":" + config.DBPort + "/" + config.Database

	pool, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return pool
}
