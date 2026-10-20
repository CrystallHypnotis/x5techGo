package main

import "fmt"

// Обезличенный разбор. Запуск из папки IDE: go run ./02_anonymous
// Ошибки поведения: примеры 1 и 2.
// Упрощения и читаемость: примеры 3–9. Допустимый вариант: пример 10.
// Имена функций добавлены для демонстрации, фрагменты обезличены.

// 1. ОШИБКА: исключён только ноль, а отрицательные числа тоже делятся на 3 и 5.
// Вопрос группе: что вернётся для -3 и -15? По условию нужно invalid.
func fizzBuzzNegativeBug(n int) string {
	switch {
	case n == 0:
		return "invalid"
	case n%15 == 0:
		return "FizzBuzz"
	case n%3 == 0:
		return "Fizz"
	case n%5 == 0:
		return "Buzz"
	case n > 0:
		return ""
	default:
		return "invalid"
	}
}

// Исправление: исключаем сразу все n <= 0.
func fizzBuzzFixed(n int) string {
	switch {
	case n <= 0:
		return "invalid"
	case n%15 == 0:
		return "FizzBuzz"
	case n%3 == 0:
		return "Fizz"
	case n%5 == 0:
		return "Buzz"
	default:
		return ""
	}
}

// 2. ОШИБКА: переводим каждый элемент в float64 до сложения.
// Вопрос: куда денется единица между большим положительным и отрицательным числом?
// На [9007199254740992, 1, -9007199254740992] получается 0 вместо 1/3.
func averageFloatBug(values []int64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, value := range values {
		sum += float64(value)
	}
	return sum / float64(len(values))
}

// Исправление: сумма по условию помещается в int64, преобразуем перед делением.
func averageFixed(values []int64) float64 {
	if len(values) == 0 {
		return 0
	}
	return float64(sumFixed(values)) / float64(len(values))
}

// Вариант по условию: начинаем с первого элемента, сравниваем с остальными.
func minFixed(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}
	minimum := values[0]
	for _, value := range values[1:] {
		if value < minimum {
			minimum = value
		}
	}
	return minimum
}

// 3. ЧИТАЕМОСТЬ: switch внутри switch работает правильно, но добавляет вложенность.
// Вопрос: можно ли сохранить поведение с одним switch? См. fizzBuzzFixed.
func fizzBuzzNested(n int) string {
	switch {
	case n > 0:
		switch {
		case n%15 == 0:
			return "FizzBuzz"
		case n%3 == 0:
			return "Fizz"
		case n%5 == 0:
			return "Buzz"
		default:
			return ""
		}
	default:
		return "invalid"
	}
}

// 4. УПРОЩЕНИЕ: проверки не меняют результат.
// У nil длина тоже 0. range по пустому вводу не выполнит тело цикла.
func sumWithExtraChecks(values []int64) int64 {
	if values == nil || len(values) == 0 {
		return 0
	}
	var sum int64 = 0
	for _, value := range values {
		sum += value
	}
	if sum == 0 {
		return 0
	}
	return sum
}

// Вариант проще: переменная уже начинается с 0, return вернёт любую сумму.
func sumFixed(values []int64) int64 {
	var sum int64
	for _, value := range values {
		sum += value
	}
	return sum
}

// 5. УПРОЩЕНИЕ: нижние границы уже исключены предыдущими ветками.
func categoryWithExtraBounds(amount int64) string {
	switch {
	case amount < 0:
		return "invalid"
	case amount >= 0 && amount <= 9999:
		return "small"
	case amount >= 10000 && amount <= 49999:
		return "regular"
	default:
		return "large"
	}
}

func categoryFixed(amount int64) string {
	switch {
	case amount < 0:
		return "invalid"
	case amount < 10_000:
		return "small"
	case amount < 50_000:
		return "regular"
	default:
		return "large"
	}
}

// 6. ЧИТАЕМОСТЬ: много переменных, имена не объясняют назначение.
// Результат правильный. Короткий вариант: averageFixed выше.
func averageWithTemps(values []int64) float64 {
	if len(values) == 0 {
		return 0
	}
	var tmp int64
	tmp = sumFixed(values)
	var res float64
	var tmp1 float64
	var tmp2 float64
	tmp2 = float64(tmp)
	tmp1 = float64(len(values))
	res = tmp2 / tmp1
	return res
}

// Ещё одно упрощение: количество элементов уже известно через len(values).
// Для сравнения со студенческим фрагментом здесь сохранён отдельный счётчик.
// Складываем в int64, чтобы обсуждать счётчик отдельно от ошибки примера 2.
func averageWithCounter(values []int64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum int64
	count := 0
	for _, value := range values {
		sum += value
		count++
	}
	return float64(sum) / float64(count)
}

// 7. ИМЕНА: values до цикла является слайсом, в теле цикла является элементом.
// Код работает. Лучше назвать элемент value. Области видимости разберём позже.
func minWithSameName(values []int64) (mint int64) {
	if len(values) == 0 {
		return 0
	}
	mint = values[0]
	for _, values := range values {
		mint = min(mint, values)
	}
	return
}

// Имя len скрывает встроенную функцию. Лучше использовать count.
func runeCountNamedLen(text string) int {
	var len int
	for range text {
		len++
	}
	return len
}

func runeCountFixed(text string) int {
	var count int
	for range text {
		count++
	}
	return count
}

// 8. ЧИСТОТА КОДА: если все ветки возвращают результат, следующий return лишний.
// Исходный фрагмент оставлен в комментарии, чтобы go vet не ругался на него:
//
// switch {
// case n <= 0:
//     return "invalid"
// ...
// default:
//     return ""
// }
// return "" // Сюда выполнение уже не дойдёт.
//
// Выполненные TODO тоже удаляем: они больше не описывают состояние работы.

// 9. ФОРМАТИРОВАНИЕ: исходная запись показана в комментарии.
// До:    if (len(values) == 0) {return 0}
// После: обычный блок if, как в minFixed.
// Команда перед сдачей: gofmt -w tasks.go

// 10. ДОПУСТИМЫЙ ВАРИАНТ: встроенные min/max внутри собственного цикла.
// Начинаем с первого элемента, затем сами обходим входные значения.
// Вопрос группе: как записать то же сравнение через if?
func minWithBuiltin(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}
	res := values[0]
	for _, value := range values {
		res = min(res, value)
	}
	return res
}

func maxWithBuiltin(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}
	res := values[0]
	for _, value := range values {
		res = max(res, value)
	}
	return res
}

// НЕ ОШИБКИ:
// values[1:] при длине 1 безопасен после проверки непустого ввода.
// Сравнение первого элемента с самим собой в Min/Max корректно.
// Индексный цикл в статистике допустим; range обязателен для RuneCount.
// n%3 == 0 && n%5 == 0 эквивалентно кратности 15.
// Отдельный if n <= 0 перед switch допустим.

func main() {
	fmt.Println("1. Отрицательные кратные числа: ошибка и исправление")
	for _, n := range []int{-3, -5, -15} {
		fmt.Printf("n=%d: исходный=%q исправленный=%q\n", n, fizzBuzzNegativeBug(n), fizzBuzzFixed(n))
	}

	fmt.Println("\n2. Потеря единицы при сложении в float64")
	values := []int64{9007199254740992, 1, -9007199254740992}
	fmt.Printf("%v: исходный=%g исправленный=%g\n", values, averageFloatBug(values), averageFixed(values))

	fmt.Println("\n3. Вложенный switch: результаты одинаковые")
	for _, n := range []int{-15, 15, 3, 5, 7} {
		fmt.Printf("n=%d: вложенный=%q простой=%q\n", n, fizzBuzzNested(n), fizzBuzzFixed(n))
	}

	fmt.Println("\n4. Лишние проверки в Sum")
	for _, values := range [][]int64{nil, {}, {5, -5}, {1, 2}} {
		fmt.Printf("%v: с проверками=%d без проверок=%d\n", values, sumWithExtraChecks(values), sumFixed(values))
	}

	fmt.Println("\n5. Повторные границы Category")
	for _, amount := range []int64{-1, 0, 9999, 10000, 49999, 50000} {
		fmt.Printf("%d: исходный=%q упрощённый=%q\n", amount, categoryWithExtraBounds(amount), categoryFixed(amount))
	}

	fmt.Println("\n6. Лишние переменные и счётчик")
	values = []int64{1, 2}
	fmt.Printf("%v: с переменными=%g со счётчиком=%g простой=%g\n",
		values, averageWithTemps(values), averageWithCounter(values), averageFixed(values))

	fmt.Println("\n7. Неясные имена: результат остаётся правильным")
	fmt.Println("Минимум:", minWithSameName([]int64{-8, -2, -5}))
	fmt.Println("Кодовых точек:", runeCountNamedLen("Goшник"), runeCountFixed("Goшник"))
	fmt.Println("\n8–9. Недостижимый код, TODO и форматирование: смотрим комментарии выше")

	fmt.Println("\n10. Допустимый вариант: встроенные min/max внутри цикла")
	for _, values := range [][]int64{nil, {42}, {-8, -2, -5}} {
		fmt.Printf("%v: минимум=%d максимум=%d\n", values, minWithBuiltin(values), maxWithBuiltin(values))
	}
}
