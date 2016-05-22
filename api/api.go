package api

type Sign struct {
	PictureName string
	SignName    string
	Value       int
}

type Signs []Sign

type User struct {
	Key		string			`json:"key"`
	Login		string			`json:"login"`
	Password	string			`json:"password"`
	Subscriptions	map[string]string	`json:"subscriptions"`
	SubscribersIP	[]string		`json:"subscribersIP"`
	PreviusResults	[]string		`json:"previusResults"`
}

type Users []User