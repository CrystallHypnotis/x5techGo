package main

import "fmt"

type Notifier interface {
	Send(string)
	GetAddress() string
}

type Email struct {
	address string
}

func (Email) Send(message string) {
	fmt.Println("email:", message)
}
func (e Email) GetAddress() string {
	return e.address
}

type Push struct {
	address string
}

type SMS struct {
	phoneNumber string
}

func (SMS) Send(message string) {
	fmt.Println("sms:", message)
}

func (s SMS) GetAddress() string {
	return ""
}

func (Push) Send(message string) {
	fmt.Println("push:", message)
}

func (p Push) GetAddress() string {
	return p.address
}

func Notify(n Notifier, message string) {
	n.Send(message)
}

var _ Notifier = Email{}

var _ Notifier = Push{}

func main() {
	Notify(Email{}, "Заказ готов")
	Notify(Push{}, "Заказ готов")
}
