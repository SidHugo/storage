package api

type Sign struct {
	PictureName string
	SignName    string
	Value       int
}

type QueryStats struct {
	AvgWriteQueryTimeMs  int64
	AvgReadQueryTimeMs   int64
	LastWriteQueryTimeMs int64
	LastReadQueryTimeMs  int64
}

type Signs []Sign
