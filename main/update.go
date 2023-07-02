package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	//"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Response events.APIGatewayProxyResponse

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
	user := "yairsggo"                           //os.Getenv("DB_USER")
	dbpass := "MbuwvGgJcC-nXskeCQnhunp8C93XC2-p" //os.Getenv("DB_PASS")
	dbhost := "rajje.db.elephantsql.com"         //os.Getenv("DB_HOST")
	dbservice := "yairsggo"                      //os.Getenv("DB_SERVICE")
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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Default().Println("connect postgres err:", err)
		return db, err
	}
	log.Default().Println("connect postgres successfully")
	return db, err
}

type RequestBodyAPIGW struct {
	RequestID string  `json:"requestId"`
	Data      UserDTO `json:"data"`
}

type UserDTO struct {
	ID    int    `json:"userId"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func UpdateUser(ctx context.Context, eventReq events.APIGatewayProxyRequest) (Response, error) {
	var (
		req  = RequestBodyAPIGW{}
		resp = Response{
			StatusCode:      400,
			IsBase64Encoded: false,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
	)
	err := json.Unmarshal([]byte(eventReq.Body), &req)
	if err != nil {
		resp.Body = ParseResponse(HttpResponse{
			Uuid: req.RequestID,
			Err:  err,
		})
		return resp, nil
	}
	db, err := InitPostgres()

	if err != nil {
		resp.Body = ParseResponse(HttpResponse{
			Uuid: req.RequestID,
			Err:  err,
		})
		resp.StatusCode = 500
		return resp, nil
	}
	resp.StatusCode = 200
	err = db.Exec(`update ruonnv1 set name = ?, phone = ? where id = ?`, req.Data.Email, req.Data.Phone, req.Data.ID).Error
	if err != nil {
		resp.Body = ParseResponse(HttpResponse{Uuid: req.RequestID, Err: err})
		return resp, nil
	}
	resp.Body = ParseResponse(HttpResponse{Uuid: req.RequestID, Data: nil})
	return resp, nil
}

func main() {
	lambda.Start(UpdateUser)
}

type HttpResponse struct {
	Uuid string // uuid, indicator per api
	Err  error
	Time string // time tracing
	Data interface{}
}

func ParseResponse(respBody HttpResponse) string {
	respBody.Time = time.Now().Format("2006-01-02T15:04:05.000-07:00")
	if respBody.Err != nil {
		return responseErr(respBody)
	}
	return responseOk(respBody)
}

func responseOk(respBody HttpResponse) string {
	var buf bytes.Buffer
	mapRes := map[string]interface{}{
		"responseId":      respBody.Uuid,
		"responseMessage": "successfully",
		"responseTime":    respBody.Time,
	}
	if respBody.Data != nil {
		mapRes["data"] = respBody.Data
	}
	body, errMarshal := json.Marshal(mapRes)
	if errMarshal != nil {
		log.Default().Println("marshal response err", errMarshal)
	}
	json.HTMLEscape(&buf, body)
	return buf.String()
}

func responseErr(respBody HttpResponse) string {
	var buf bytes.Buffer
	mapRes := map[string]interface{}{
		"responseId":      respBody.Uuid,
		"responseMessage": respBody.Err.Error(),
		"responseTime":    respBody.Time,
	}

	body, errMarshal := json.Marshal(mapRes)
	if errMarshal != nil {
		log.Default().Println("marshal response err", errMarshal)
	}
	json.HTMLEscape(&buf, body)
	return buf.String()
}
