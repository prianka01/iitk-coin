package model

type User struct {
	Name               string  `json:"name"`
	Password           string  `json:"password"`
	Rollno             int     `json:"rollno"`
	Token              string  `json:"token"`
	Access             string  `json:"access"`
	Coins              float64 `json:"coins"`
	EventsParticipated int     `json:"events"`
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

type Transaction struct {
	Type      string  `json:"type"`
	Sender    int     `json:"sender"`
	Reciever  int     `json:"reciever"`
	Amount    int     `json:"amount"`
	Tax       float64 `json:"tax"`
	TimeStamp string  `json:"timestamp"`
	Info      string  `json:"info"`
}
type RedeemRequest struct {
	Status string `json:"status"`
	Item   string `json:"item"`
	Sender int    `json:"sender"`
	Info   string `json:"info"`
}

type StoreItems struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}