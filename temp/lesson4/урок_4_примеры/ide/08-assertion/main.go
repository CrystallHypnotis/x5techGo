package main

import "fmt"

type Notifier interface {
	Send(string)
}

type Email struct {
	Address string
}

func (Email) Send(string) {

}

type Push struct{}

func (Push) Send(string) {

}

func main() {
	for _, n := range []Notifier{
		Email{
			Address: "student@example.com",
		},
		Push{},
		&Email{
			Address: "pointer@example.com",
		},
	} {
		e, ok := n.(*Email)
		fmt.Printf("%T: ok=%v, address=%q\n", n, ok, e)
	}
	var n Notifier = Email{}
	_, ok := n.(interface {
		Send(string)
	})
	fmt.Println("подходит другому интерфейсу:", ok)
}
