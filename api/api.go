package api

type Sign struct {
	PictureName string
	SignName    string
	Value       int
}

type Sign2 struct {
	Link         string `json:"link"`
	Base64string string `json:"base64string"`
}

type Signs2 []Sign2

type QueryStats struct {
	AvgWriteQueryTimeMs  int64
	AvgReadQueryTimeMs   int64
	LastWriteQueryTimeMs int64
	LastReadQueryTimeMs  int64
	MaxWriteQueryTimeMs  int64
	MaxReadQueryTimeMs   int64
}

type Signs []Sign

type User struct {
	Key            int            `json:"key"`
	Login          string         `json:"login"`
	Password       string         `json:"password"`
	Subscriptions  map[string]int `json:"subscriptions"`
	SubscribersIP  []string       `json:"subscribersIP"`
	PreviusResults []string       `json:"previusResults"`
}

type Users []User

// Test structure for DB insert-retreive testing
type TestStruct struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
