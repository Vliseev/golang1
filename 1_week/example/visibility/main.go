package main

import (
	"mail_course/1_week/example/visibility/person"
	"fmt"
)

func main() {
	p := person.NewPerson(1, "rvasily", "secret")

	// p.secret undefined (cannot refer to unexported field or method secret)
	// fmt.Printf("main.PrintPerson: %+v\n", p.secret)

	secret := person.GetSecret(p)
	fmt.Println("GetSecret", secret)
}
