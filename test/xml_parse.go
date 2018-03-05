package main

import (
	"encoding/xml"
	"io"
	"fmt"
	"bufio"
	"os"
	"strings"
	"sort"
	"text/tabwriter"
)

type UserRow struct {
	Id     int    `xml:"id"`
	Name   string `xml:"first_name"`
	Age    int    `xml:"age"`
	About  string `xml:"about"`
	Gender string `xml:"gender"`
}

type byId []*UserRow

func (x byId) Len() int  { return len(x)}
func (x byId) Less(i,j int) bool { return x[i].Id < x[j].Id}
func (x byId) Swap(i,j int) {x[i],x[j] = x[j],x[i]}

type byName []*UserRow

func (x byName) Len() int  { return len(x)}
func (x byName) Less(i,j int) bool { return x[i].Name < x[j].Name}
func (x byName) Swap(i,j int) {x[i],x[j] = x[j],x[i]}

type byAge []*UserRow

func (x byAge) Len() int  { return len(x)}
func (x byAge) Less(i,j int) bool { return x[i].Age < x[j].Age}
func (x byAge) Swap(i,j int) {x[i],x[j] = x[j],x[i]}



func UserDecoderFilter(inp io.Reader,f func(*UserRow) bool) []*UserRow  {
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

    var	usersPointers []*UserRow
    for i:=0;i<len(users);i++{
    	usersPointers = append(usersPointers,&users[i])
	}
	return usersPointers

}

func printUsers(users []*UserRow) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Id", "Name", "Age", "About", "Gender")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range users {
		fmt.Fprintf(tw, format, t.Id, t.Name, t.Age, t.About, t.Gender)
	}
	tw.Flush() // calculate column widths and print table
}

func main() {
	f,err:=os.Open("/home/vad/GO/src/coursera/golang1/hw4_test_coverage/dataset.xml")
	if err!= nil{
		fmt.Println("err")
	}
	query:=""
	users:=UserDecoderFilter(f,func(r *UserRow) bool{
		qr:=query
		if strings.Contains(r.Name,qr) || strings.Contains(r.About,qr) {
			return true
		}else {
			return false
		}
	})

	printUsers(users)
	sort.Sort(byName(users))
	printUsers(users)

}
