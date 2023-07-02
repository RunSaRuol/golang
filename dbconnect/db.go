package dbconnect

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getEnv(key, fallback string) string {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// 	user :=  "yairsggo"//os.Getenv("DB_USER")
// 	dbpass := "MbuwvGgJcC-nXskeCQnhunp8C93XC2-p"//os.Getenv("DB_PASS")
// 	dbhost := "rajje.db.elephantsql.com"//os.Getenv("DB_HOST")
// 	dbservice := "yairsggo"//os.Getenv("DB_SERVICE")

type Postgres struct {
	Username    string `yaml:"username" mapstructure:"username"`
	Password    string `yaml:"password" mapstructure:"password"`
	Database    string `yaml:"database" mapstructure:"database"`
	Host        string `yaml:"host" mapstructure:"host"`
	Port        int    `yaml:"port" mapstructure:"port"`
	Schema      string `yaml:"schema" mapstructure:"schema"`
	MaxIdleConn int    `yaml:"max_idle_conn" mapstructure:"max_idle_conn"`
	MaxOpenConn int    `yaml:"max_open_conn" mapstructure:"max_open_conn"`
}

func loadConfig() Postgres {
	user := getEnv("DB_USER", "yairsggo")
	dbpass := getEnv("DB_PASS", "MbuwvGgJcC-nXskeCQnhunp8C93XC2-p")
	dbhost := getEnv("DB_HOST", "rajje.db.elephantsql.com")
	dbservice := getEnv("DB_SERVICE", "yairsggo")
	return Postgres{
		Username: user,
		Password: dbpass,
		Database: dbservice,
		Host:     dbhost,
		Port:     5432,
	}
}

// create database postgres instance
func InitPostgres() (*gorm.DB, error) {
	config := loadConfig()
	log.Default().Println("connecting postgres database")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d ", config.Host, config.Username, config.Password, config.Database, config.Port)
	log.Default().Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	sqlDB, err := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		log.Default().Println("connect postgres err:", err)
		return db, err
	}
	log.Default().Println("connect postgres successfully")
	return db, err
}
