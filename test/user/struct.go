package user

import (
	"fmt"
)

type User struct {
	ID       int `json:"user_id"`
	Username string
	Phone    string `json:"phone"`
}

var jsonStr = `{"user_id": 42, "username": "rvasily", "phone": "1052"}`


func Test(){

	data := []byte(jsonStr)
	u := &User{}
	u.UnmarshalJSON(data)

	fmt.Println(u)
}