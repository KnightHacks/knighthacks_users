package model

import "encoding/json"

func (u User) String() string {
	marshal, _ := json.Marshal(u)
	return string(marshal)
}
