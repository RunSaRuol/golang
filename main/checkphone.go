package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"example.com/m/gohttp"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

type Response events.APIGatewayProxyResponse

func getEnv1(key, fallback string) string {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var authKey = getEnv1("API_KEY", "B5d4JtTU8u1ggV8gp7OF88gcCGxZls6T3f5PYZSa")
var urlCreateUser = getEnv1("URL_CREATE_USER", "https://jqsr6098w3.execute-api.us-east-1.amazonaws.com/dev/main/create")
var urlCheckPhone = getEnv1("URL_CHECK_PHONE", "https://1g1zcrwqhj.execute-api.ap-southeast-1.amazonaws.com/dev/testapi")

func checkphone(ctx context.Context,
	eventReq events.APIGatewayProxyRequest) (Response, error) {
	var (
		req  = gohttp.RequestBodyAPIGW{}
		resp = Response{
			StatusCode:      200,
			IsBase64Encoded: false,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
	)
	log.Default().Println("------request API CreateUser: ------", req)
	err := json.Unmarshal([]byte(eventReq.Body), &req)
	if err != nil {
		resp.Body = gohttp.ParseResponse(gohttp.HttpResponse{
			Uuid: req.RequestID,
			Err:  err,
			Code: "01",
		})
		return resp, nil
	}
	//Convert phone string to int
	phone, err := strconv.Atoi(req.Data.Phone)

	if err != nil {
		resp.Body = gohttp.ParseResponse(gohttp.HttpResponse{
			Uuid: req.RequestID,
			Err:  err,
			Code: "01",
		})
		return resp, nil
	}

	//Call api checkphone
	log.Default().Println("------START call api check phone------")
	payload := []byte(fmt.Sprintf(`{
		"requestId":"%s",
		"data": {
			"value": %d
		}
	}`, req.RequestID, phone))

	DataRes := gohttp.Post(urlCheckPhone, payload, authKey)
	var apiRes gohttp.APIResponse

	err1 := json.Unmarshal([]byte(DataRes), &apiRes)
	if err1 != nil {
		resp.Body = gohttp.ParseResponse(gohttp.HttpResponse{
			Uuid: req.RequestID,
			Err:  err,
			Code: "01",
		})
		return resp, nil
	}

	if apiRes.ResponseCode == "00" {
		log.Default().Println("----apiRes-----", apiRes)
		log.Default().Println("------START call api createUser------")
		now := time.Now()
		payloadCR := []byte(fmt.Sprintf(`{
		"requestId":"%s",
		"requestTime": "%s",
		"signature": "%s",
		"data": {
			"username": "%s",
            "name": "%s",
            "phone": "%s"
		}
	}`, req.RequestID, now, req.Signature, req.Data.User, req.Data.Name, req.Data.Phone))

		DataRes1 := gohttp.Post(urlCreateUser, payloadCR, authKey)
		var apiRes1 gohttp.APIResponse
		err2 := json.Unmarshal([]byte(DataRes1), &apiRes1)
		log.Default().Println("----apiRes-----", apiRes1)
		if err2 != nil {
			resp.Body = gohttp.ParseResponse(gohttp.HttpResponse{
				Uuid:    req.RequestID,
				Err:     err,
				Code:    "01",
				Message: "Call Api create user fail",
			})
			log.Default().Println("------------------------------createUser-err2: ---------------", resp)
			return resp, nil
		}
		if apiRes1.ResponseCode != "00" {
			resp.Body = gohttp.ParseResponse(gohttp.HttpResponse{
				Uuid:    apiRes1.ResponseId,
				Err:     err,
				Code:    apiRes1.ResponseCode,
				Message: apiRes1.ResponseMessage,
			})
			return resp, nil
		}
		resp.Body = gohttp.ParseResponse(gohttp.HttpResponse{
			Uuid:    req.RequestID,
			Err:     err,
			Code:    "00",
			Message: "CreateUser Successfully!",
		})
		log.Default().Println("------------------------------createUser-success: ---------------", resp)
		log.Default().Println("------------------------------END call api createUser ---------------")
		return resp, nil
	}
	resp.Body = gohttp.ParseResponse(gohttp.HttpResponse{
		Uuid:    req.RequestID,
		Err:     err,
		Code:    "01",
		Message: "Your phone number is wrong!!",
	})
	log.Default().Println("------------------------------createUser end: ---------------", resp)
	return resp, nil
}

func validate(value string) bool {
	log.Default().Println("------------------------------Begin with : ---------------", value)
	if len(strings.TrimSpace(value)) == 0 {
		return false
	}
	return true
}

func main() {
	lambda.Start(checkphone)
	// Create(context.Background(), events.APIGatewayProxyRequest{})
}
