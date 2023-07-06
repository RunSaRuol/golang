package gohttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Data struct {
	Value int
}

type RequestBody struct {
	RequestId string
	Data      Data
}

type Response struct {
	Result string
}

func Post(url string, bodyReq []byte, authKey string) string {
	timeout := os.Getenv("TIME_OUT")
	timeoutInt, err := strconv.Atoi(timeout)
	method := "POST"
	// payload := []byte(fmt.Sprintf(`{
	// 	"requestId":"%s",
	// 	"data": {
	// 		"value": %d
	// 	}
	// }`, requestId, phone))

	client := &http.Client{
		Timeout: time.Duration(timeoutInt) * time.Second,
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyReq))

	if err != nil {
		fmt.Println(err)
		return "INTERNAL SERVER ERROR"
	}
	req.Header.Add("x-api-key", authKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	log.Default().Println("---------------------<<url>> :---------------", url)
	log.Default().Println("------------------------------res---------------", res)
	if err != nil {
		fmt.Println(err)
		return "Call API fail"
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	fmt.Println("CALL API SUCCESS:", string(responseBody))
	if err != nil {
		log.Default().Println("------------------------------errr---------------", err)
		return "Data fail"
	}
	//fmt.Println(string(body))

	return string(responseBody)
}
