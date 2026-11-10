package main

import "fmt"

type Email struct{}

func (Email) Send(message string) {
	fmt.Println("email:", message)
}

type Push struct{}

func (Push) Send(message string) {
	fmt.Println("push:", message)
}

func NotifyEmail(n Email, message string) {
	n.Send(message)
}

func main() {
	NotifyEmail(Email{}, "Заказ готов")
	NotifyEmail(Push{}, "Заказ готов") // Другой конкретный тип не подходит.
}
