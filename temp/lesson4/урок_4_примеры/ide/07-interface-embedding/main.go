package main

import "fmt"

type Sender interface {
	Send(string)
}

type Named interface {
	Name() string
}

type NamedSender interface {
	Sender
	Named
}

type Console struct{}

func (Console) Send(message string) {
	fmt.Println("console:", message)
}

func (Console) Name() string {
	return "console"
}

func Show(n NamedSender) {
	fmt.Println(n.Name())
	n.Send("Заказ готов")
}

func main() {
	Show(Console{})
}
