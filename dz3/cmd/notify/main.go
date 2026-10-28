package main

import (
	"errors"
	"fmt"
	"lesson5/student/practice"
)

func main() {
	// TODO: домашняя работа. Покажите успешную рассылку, пустой текст,
	// ошибку настройки и продолжение после BrokenNotifier.
	// Импорт учебного пакета: "lesson5/student/practice".

	email := &practice.EmailNotifier{BaseNotifier: practice.BaseNotifier{Channel: "Email"}, Address: "kekw@test.com"}
	push := &practice.PushNotifier{BaseNotifier: practice.BaseNotifier{Channel: "Push"}, DeviceID: "88005553555"}
	err := practice.Dispatch([]practice.Notifier{email, push}, "Заказ готов")
	fmt.Println("Ошибки нет:", err)
	fmt.Println("История email:", email.Messages)
	fmt.Println("История push:", push.Messages)

	err = practice.Dispatch([]practice.Notifier{email, push}, "")
	if errors.Is(err, practice.ErrEmptyMessage) {
		fmt.Println("Ошибка распознана:", err)
	}
	fmt.Println("История email:", email.Messages)
	fmt.Println("История push:", push.Messages)

	invalid_push := &practice.PushNotifier{BaseNotifier: practice.BaseNotifier{Channel: "Push"}}
	err = practice.Dispatch([]practice.Notifier{invalid_push}, "kekw")
	var send_err *practice.SendError
	if errors.As(err, &send_err) {
		fmt.Printf("Ошибка: %q:%q\n", send_err.Channel, send_err.Reason)
	}

	valid_push := &practice.PushNotifier{BaseNotifier: practice.BaseNotifier{Channel: "Push"}, DeviceID: "device-2"}
	broken_push := &practice.BrokenNotifier{BaseNotifier: practice.BaseNotifier{Channel: "broken"}}
	errs := practice.SafeDispatch([]practice.Notifier{broken_push, valid_push}, "kekw")
	for _, err := range errs {
		if errors.Is(err, practice.ErrNotifierPanic) {
			fmt.Println("Ошибка распознана:", err)
		}
	}
	fmt.Println("valid_push:", valid_push.Messages)
}
