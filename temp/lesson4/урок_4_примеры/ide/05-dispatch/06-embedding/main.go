package main

import "fmt"

type Base struct {
	Channel string
}

func (b Base) Log(message string) {
	fmt.Printf("[%s] %s\n", b.Channel, message)
}

func (b Base) Name() string {
	return "base"
}

func (b Base) Report() {
	fmt.Println("Report:", b.Name())
}

type Email struct {
	Base
	Address string
}

func (Email) Name() string {
	return "email"
}

func main() {
	e := Email{
		Base: Base{
			Channel: "email",
		},
		Address: "student@example.com",
	}
	e.Log("Заказ готов")
	e.Base.Log("Полная запись")
	fmt.Println("Name:", e.Name())
	e.Report()
	// var b Base = e // Встраивание не делает Email значением типа Base.
}
