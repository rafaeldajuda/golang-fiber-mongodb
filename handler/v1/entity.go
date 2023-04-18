package v1

import "time"

type MsgError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type MsgOK struct {
	ID interface{} `json:"_id" bson:"_id"`
}

type Animal struct {
	ID        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Owner     string    `json:"owner" bson:"owner"`
	Type      string    `json:"type" bson:"type"`
	Age       int       `json:"age" bson:"age"`
	Castrated bool      `json:"castrated" bson:"castrated"`
	Surgery   time.Time `json:"surgery" bson:"surgery"`
}
