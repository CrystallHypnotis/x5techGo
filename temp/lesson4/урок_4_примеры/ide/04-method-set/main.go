package main

import "fmt"

type Notifier interface {
	Send(string)
}

type Counter struct {
	Calls int
}

func (c *Counter) Send(string) {
	c.Calls++
}

var _ Notifier = (*Counter)(nil)

func main() {
	c := Counter{}
	c.Send("прямой вызов")
	var n Notifier = &c
	n.Send("через интерфейс")
	fmt.Println("вызовы:", c.Calls)
	var bad Notifier = c // Send есть в наборе методов *Counter.
}
