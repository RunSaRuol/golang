package gohttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func Post(url string, requestId string, phone int) string {
	method := "POST"
	payload := []byte(fmt.Sprintf(`{
		"requestId":"%s",
		"data": {
			"value": %d	
		}
	}`, requestId, phone))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println(err)
		return "INTERNAL SERVER ERROR"
	}
	req.Header.Add("x-api-key", "B5d4JtTU8u1ggV8gp7OF88gcCGxZls6T3f5PYZSa")
	req.Header.Add("Content-Type", "text/plain")

	res, err := client.Do(req)
	log.Default().Println("------------------------------res---------------", res)
	if err != nil {
		fmt.Println(err)
		return "Call API fail"
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	fmt.Println("CALL API SUCCESS:", string(responseBody))
	if err != nil {
		fmt.Println(err)
		return "Data fail"
	}
	//fmt.Println(string(body))

	return "SUCCESS"
}
