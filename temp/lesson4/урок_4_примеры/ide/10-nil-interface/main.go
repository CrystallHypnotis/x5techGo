package main

import "fmt"

type Notifier interface {
	Send(string)
}

type Email struct {
	Address string
}

func (e *Email) Send(message string) {
	if e == nil {
		fmt.Println("получатель nil")
		return
	}
	fmt.Println(e.Address, message)
}

func main() {
	var n Notifier
	fmt.Println("пустой интерфейс:", n == nil)
	fmt.Printf("тип интерфейса: %T\n", n)
	var p *Email
	n = p
	fmt.Println("указатель:", p == nil)
	fmt.Println("интерфейс с указателем:", n == nil)
	fmt.Printf("динамический тип: %T\n", n)
	n.Send("Заказ готов")
}
