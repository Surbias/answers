package model

type Answer struct {
	Timestamp int64  `json:"-" bson:"timestamp"`
	Key       string `json:"key" bson:"key"`
	Value     string `json:"value" bson:"value"`
}

type UpdateAnswer struct {
	Value string `json:"value"`
}
