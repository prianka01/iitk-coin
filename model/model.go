package model

type User struct {
	Name          string `json:"name"`
	Password      string `json:"password"`
	Rollno        int    `json:"rollno"`
	Token         string `json:"token"`
	CanAccessPage bool   `json:"accesspage"`
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}