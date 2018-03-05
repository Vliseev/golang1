package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	badToken      = "badAccesToken"
	unknownErrResp = "unknownErrResp"
	timeoutError  = "timeoutError"
	unknownError  = "unknownError"
	internalError = "StatusInternalServerError"
	errUnpack     = "errUnpackJson"
	badOrderField = "ErrorBadOrderField"
	cantUnpackRes = "cantUnpackResultJson"
)

type (
	sortType func(data sort.Interface)
	TestCase struct {
		ID      string
		Client  *SearchClient
		Request *SearchRequest
		Result  *SearchResponse
		IsError bool
	}
)

type (
	UserRow struct {
		Id     int    `xml:"id"`
		Name   string `xml:"first_name"`
		Age    int    `xml:"age"`
		About  string `xml:"about"`
		Gender string `xml:"gender"`
	}
)

type byId []*UserRow

func (x byId) Len() int           { return len(x) }
func (x byId) Less(i, j int) bool { return x[i].Id < x[j].Id }
func (x byId) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type byName []*UserRow

func (x byName) Len() int           { return len(x) }
func (x byName) Less(i, j int) bool { return x[i].Name < x[j].Name }
func (x byName) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type byAge []*UserRow

func (x byAge) Len() int           { return len(x) }
func (x byAge) Less(i, j int) bool { return x[i].Age < x[j].Age }
func (x byAge) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func UserDecoderFilter(inp io.Reader, f func(*UserRow) bool) []*UserRow {
	input := bufio.NewReader(inp)
	decoder := xml.NewDecoder(input)
	var users []UserRow
	var user UserRow
	for {
		tok, tokenErr := decoder.Token()
		if tokenErr != nil && tokenErr != io.EOF {
			fmt.Println("error happend", tokenErr)
			break
		} else if tokenErr == io.EOF {
			break
		}
		if tok == nil {
			fmt.Println("t is nil break")
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			if tok.Name.Local == "row" {
				if err := decoder.DecodeElement(&user, &tok); err != nil {
					fmt.Println("error happend", err)
				}
				if f(&user) == true {
					users = append(users, user)
				}
			}
		}
	}

	var usersPointers []*UserRow
	for i := 0; i < len(users); i++ {
		usersPointers = append(usersPointers, &users[i])
	}
	return usersPointers
}

func (sr *SearchRequest) GetParam(r *http.Request) {
	val := r.FormValue("limit")
	sr.Limit, _ = strconv.Atoi(val)

	val = r.FormValue("offset")
	sr.Offset, _ = strconv.Atoi(val)

	val = r.FormValue("query")
	sr.Query = val

	val = r.FormValue("order_field")
	sr.OrderField = val

	val = r.FormValue("OrderBy")
	sr.OrderBy, _ = strconv.Atoi(val)
}

func ErrHandler(w http.ResponseWriter, r *http.Request) {

	switch r.FormValue("query") {
	case timeoutError:
		time.Sleep(time.Second)
		return

	}
	if r.FormValue("query") == timeoutError {
		time.Sleep(time.Second)
		return
	}
	if r.FormValue("query") == internalError {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if r.FormValue("query") == errUnpack {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error"))
		return
	}
	if r.FormValue("query") == badOrderField {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Error":"ErrorBadOrderField"}`))
		return
	}
	if r.FormValue("query") == unknownErrResp {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Error":"unknownError"}`))
		return
	}

	sr := SearchRequest{}
	sr.GetParam(r)
	if r.Header.Get("AccessToken") == badToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println(sr)
}

func sortUsers(sr *SearchRequest, users []*UserRow) error {
	var sortQuery sort.Interface

	switch sr.OrderField {
	case "Id":
		sortQuery = byId(users)
	case "Age":
		sortQuery = byAge(users)
	case "Name":
		sortQuery = byName(users)
	case "":
		sortQuery = byName(users)
	default:
		return fmt.Errorf("OrderField is invalid")
	}

	switch sr.OrderBy {
	case -1:
		sortQuery = sort.Reverse(sortQuery)
	case 1:
	default:
		return fmt.Errorf("OrderBy is invalid")
	}

	sort.Sort(sortQuery)
	return nil
}

func NormalHandler(w http.ResponseWriter, r *http.Request) {
	sr := SearchRequest{}
	sr.GetParam(r)

	//f, err := os.Open("/home/vad/GO/src/coursera/golang1/hw4_test_coverage/dataset.xml")
	f, err := os.Open(g"dataset.xml")
	defer f.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Error`))
	}

	users := UserDecoderFilter(f, func(r *UserRow) bool {
		query := sr.Query
		if strings.Contains(r.Name, query) || strings.Contains(r.About, query) {
			return true
		} else {
			return false
		}
	})

	if sr.OrderBy != 0 {
		err := sortUsers(&sr, users)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errStr, _ := json.Marshal(SearchErrorResponse{err.Error()})
			w.Write(errStr)
		}
	}

	if sr.Offset > 0 {
		if sr.Offset >= len(users) {
			users = []*UserRow{}
		} else {
			users = users[sr.Offset:]
		}

	}
	if sr.Limit < len(users) {
		users = users[:sr.Limit]
	}

	usersJson, _ := json.Marshal(users)

	data := []User{}

	json.Unmarshal(usersJson, &data)

	w.Write(usersJson)
}

func TestSearchClient_Errors(t *testing.T) {
	cases := []TestCase{
		TestCase{
			ID:      "Timeot",
			Client:  &SearchClient{AccessToken: "tok1"},
			Request: &SearchRequest{Query: timeoutError},
			IsError: true,
		},
		TestCase{
			ID:      "unknown error",
			Client:  &SearchClient{URL: "errURL"},
			Request: &SearchRequest{Query: unknownError},
			IsError: true,
		},
		TestCase{
			ID:      "LimitErr",
			Client:  &SearchClient{},
			Request: &SearchRequest{Limit:-1},
			IsError: true,
		},
		TestCase{
			ID:      "LimitErr",
			Client:  &SearchClient{},
			Request: &SearchRequest{Offset:-1},
			IsError: true,
		},
		TestCase{
			ID:      "bad token",
			Client:  &SearchClient{AccessToken: badToken},
			Request: &SearchRequest{},
			IsError: true,
		},
		TestCase{
			ID:      "StatusInternalServerError",
			Client:  &SearchClient{},
			Request: &SearchRequest{Query: internalError},
			IsError: true,
		},
		TestCase{
			ID:      "errUnpackJson",
			Client:  &SearchClient{},
			Request: &SearchRequest{Query: errUnpack},
			IsError: true,
		},
		TestCase{
			ID:      "badOrderField",
			Client:  &SearchClient{},
			Request: &SearchRequest{Query: badOrderField},
			IsError: true,
		},
		TestCase{
			ID:      "CantUnpackResult",
			Client:  &SearchClient{},
			Request: &SearchRequest{Query: cantUnpackRes},
			IsError: true,
		},
		TestCase{
			ID:      "unknownErrResp",
			Client:  &SearchClient{},
			Request: &SearchRequest{Query: unknownErrResp},
			IsError: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(ErrHandler))
	defer ts.Close()

	for _, testCase := range cases {
		client := testCase.Client

		if client.URL == "" {
			client.URL = ts.URL
		}
		_, err := client.FindUsers(*testCase.Request)
		if err != nil && testCase.IsError!=true{
			t.Errorf(err.Error())
		}
	}
}

func TestSearchClient_FindUsers(t *testing.T) {

	cases := []TestCase{
		TestCase{
			ID:      "NextPage",
			Client:  &SearchClient{},
			Request: &SearchRequest{Limit: 27, Query: "", OrderField: "byAge", OrderBy: -1},
			IsError: false,
		},
		TestCase{
			ID:      "OnePage",
			Client:  &SearchClient{},
			Request: &SearchRequest{Limit: 25, Query: "Aliquip", OrderField: "byName", OrderBy: 1},
			IsError: false,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(NormalHandler))
	defer ts.Close()

	for _, testCase := range cases {
		client := testCase.Client

		if client.URL == "" {
			client.URL = ts.URL
		}
		_, err := client.FindUsers(*testCase.Request)
		if err != nil && testCase.IsError!=true{
			t.Errorf(err.Error())
		}
	}
}
