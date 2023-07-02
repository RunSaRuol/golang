package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
	"strings"

	//"os"

	"example.com/m/dbconnect"
	"example.com/m/gohttp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

func CreateUser(ctx context.Context,
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
		})
		return resp, nil
	}
	// //Convert phone string to int
	// phone, err := strconv.Atoi(req.Data.Phone)

	// if err != nil {
	// 	resp.Body = ParseResponse(HttpResponse{
	// 		Uuid: req.RequestID,
	// 		Err:  err,
	// 	})
	// 	return resp, nil
	// }
	// log.Default().Println("------START------")
	// statusRes := gohttp.Post(url, req.RequestID, phone)

	// if statusRes == "SUCCESS" {
	// 	log.Default().Println("----SUCCESS-----", statusRes)
	// }
	secretKey := os.Getenv("SECRET_KEY")
	payload := req.RequestID + req.Data.Phone + req.Data.User + secretKey
	h := sha256.New()
	h.Write([]byte(payload))
	sum := hex.EncodeToString(h.Sum(nil))

	log.Default().Println("------Payload: ------", payload)
	log.Default().Println("------signature: ------", sum)

	if strings.Compare(sum, req.Signature) != 0 {
		resp.Body = gohttp.ParseResponse(gohttp.HttpResponse{
			Uuid:    req.RequestID,
			Err:     err,
			Message: "Signature fail",
			Code:    "01",
		})
		resp.StatusCode = 200
		log.Default().Println("------signature faill ------")
		return resp, nil

	}
	log.Default().Println("------signature success ------")

	db, err := dbconnect.InitPostgres()
	if err != nil {
		resp.Body = gohttp.ParseResponse(gohttp.HttpResponse{
			Uuid: req.RequestID,
			Err:  err,
		})
		resp.StatusCode = 500
		return resp, nil
	}

	// check user nam exist
	exists, err := checkUsernameExists(req.Data.User)
	if err != nil {
		resp.Body = gohttp.ParseResponse(gohttp.HttpResponse{
			Uuid:    req.RequestID,
			Err:     err,
			Code:    "01",
			Message: "Cannot access db",
		})
		resp.StatusCode = 200
		return resp, nil
	}

	if exists {
		resp.Body = gohttp.ParseResponse(gohttp.HttpResponse{
			Uuid:    req.RequestID,
			Err:     err,
			Code:    "01",
			Message: "User name is exist",
		})
		resp.StatusCode = 200
		return resp, nil
	}
	log.Default().Println("------User name is exist ------")

	// set http-code 200
	resp.StatusCode = 200
	// save new user
	err = db.Debug().Exec(`insert into ruonnv1(username,name,phone) values(?,?,?)`, req.Data.User, req.Data.Name, req.Data.Phone).Error
	if err != nil {
		resp.Body = gohttp.ParseResponse(gohttp.HttpResponse{Uuid: req.RequestID, Err: err, Code: "01", Message: "create user db fail"})
		return resp, nil
	}
	resp.Body = gohttp.ParseResponse(gohttp.HttpResponse{Uuid: req.RequestID, Err: err, Code: "00", Message: "Success"})
	defer func() {
		dbIns, _ := db.DB()
		dbIns.Close()
	}()
	resp.StatusCode = 500
	return resp, nil
}

func checkUsernameExists(username string) (bool, error) {
	db, err := dbconnect.InitPostgres()
	var user = gohttp.UserDTO{}
	err = db.Table("ruonnv1").Where("username = ?", username).First(&user).Error
	defer func() {
		dbIns, _ := db.DB()
		dbIns.Close()
	}()
	if err != nil {
		return false, err
	}
	return len(user.Name) > 0, nil

}
func main() {
	lambda.Start(CreateUser)
	// Create(context.Background(), events.APIGatewayProxyRequest{})
}
