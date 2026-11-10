package main

import "fmt"

type Notifier interface {
	Send(string) error
}

type Email struct{}

func (Email) Send(message string) error {
	fmt.Println("email:", message)
	return nil
}

type Push struct{}

func (Push) Send(message string) error {
	fmt.Println("push:", message)
	return nil
}

type Console struct{}

func (Console) Send(message string) error {
	fmt.Println("console:", message)
	return nil
}

func Dispatch(all []Notifier, message string) error {
	for _, n := range all {
		if err := n.Send(message); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	all := []Notifier{
		Email{},
		Push{},
		Console{},
	}
	fmt.Println("ошибка:", Dispatch(all, "Заказ готов"))
}
