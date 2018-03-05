package main

import (
	"encoding/json"
	"fmt"
)

type SearchErrorResponse struct {
	Error string
}

var jsonStr = `{"Error":"ErrorBadOrderField"}`

func main() {
	data := []byte(jsonStr)

	u := &SearchErrorResponse{}
	json.Unmarshal(data, u)
	fmt.Printf("struct:\n\t%#v\n\n", u.Error)

}