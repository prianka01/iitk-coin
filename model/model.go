package model

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Rollno   int    `json:"rollno"`
	Token    string `json:"token"`
	Access   string `json:"access"`
	Coins    int    `json:"coins"`
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}