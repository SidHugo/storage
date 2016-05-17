package api

type Sign struct {
	PictureName string
	SignName    string
	Value       int
}

type Signs []Sign

type User struct {
	key		string
	login		string
	password	string
	subscriptions	map[string]string
	subscribersIP	[]string
	previusResults	[]string
}

type Users []User