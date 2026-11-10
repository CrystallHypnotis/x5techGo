package main

import "fmt"

type Email struct {
	Address string
}

func (e Email) Label() string {
	return "email: " + e.Address
}

func (e Email) Rename(address string) {
	e.Address = address
}

func main() {
	e := Email{Address: "old@example.com"}
	fmt.Println(e.Label())
	e.Rename("new@example.com")
	fmt.Println(e.Label())
}
