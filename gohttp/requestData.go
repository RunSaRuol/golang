package gohttp

type RequestBodyAPIGW struct {
	RequestID   string  `json:"requestId"`
	RequestTime string  `json:"requestTime"`
	Signature   string  `json:"signature"`
	Data        UserDTO `json:"data"`
}

type UserDTO struct {
	Name  string `json:"name"`
	User  string `json:"userName"`
	Phone string `json:"phone"`
}
