package main

import "fmt"

type Email struct {
	Address string
}

type Push struct {
	DeviceID string
}

func Describe(value any) string {
	switch v := value.(type) {
	case Email:
		return "email: " + v.Address
	case *Email:
		if v == nil {
			return "nil *Email"
		}
		return "email pointer: " + v.Address
	case Push:
		return "push: " + v.DeviceID
	case nil:
		return "nil"
	default:
		return "неизвестный тип"
	}
}

func main() {
	for _, v := range []any{
		Email{
			"a@example.com",
		},
		&Email{
			"b@example.com",
		},
		Push{
			"phone",
		},
		42,
		"42",
		nil,
	} {
		fmt.Println(Describe(v))
	}
}
