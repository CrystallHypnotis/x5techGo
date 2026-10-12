package lesson2

// FizzBuzz возвращает FizzBuzz, Fizz, Buzz или пустую строку для обычного числа.
// Для n <= 0 возвращает invalid.
func FizzBuzz(n int) string {
	// TODO: Реализовать через switch без выражения.
	// 2. `FizzBuzz`: для положительного числа, кратного 15, верните `FizzBuzz`,
	// иначе для кратного 3: `Fizz`, кратного 5: `Buzz`. Для прочих положительных чисел верните пустую строку `""`,
	// для нуля и отрицательных: `invalid`. Проверка: `go test -run TestFizzBuzz -v`.
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

// Sum возвращает сумму. Для пустого среза и nil возвращает 0.
func Sum(values []int64) int64 {
	// TODO: Пройти по значениям циклом.
	if (len(values) == 0) || (values == nil) {
		return 0
	}
	var sum int64 = 0
	for _, value := range values {
		sum += value
	}
	return sum
}

// Min возвращает минимум. Для пустого среза и nil возвращает 0 по условию задачи.
func Min(values []int64) int64 {
	// TODO: Учесть пустой ввод, затем сравнить элементы.
	if (len(values) == 0) || (values == nil) {
		return 0
	}
	var min int64 = values[0]
	for _, value := range values {
		if value < min {
			min = value
		}
	}
	return min
}

// Max возвращает максимум. Для пустого среза и nil возвращает 0 по условию задачи.
func Max(values []int64) int64 {
	// TODO: Учесть пустой ввод, затем сравнить элементы.
	if (len(values) == 0) || (values == nil) {
		return 0
	}
	var max int64 = values[0]
	for _, value := range values {
		if value > max {
			max = value
		}
	}
	return max
}

// Average возвращает среднее как float64. Для пустого среза и nil возвращает 0.
func Average(values []int64) float64 {
	// TODO: Выполнить преобразования до деления.
	if (len(values) == 0) || (values == nil) {
		return 0
	}
	var sum int64 = 0
	for _, value := range values {
		sum += value
	}

	return float64(sum) / float64(len(values))
}

// Category возвращает invalid, small, regular или large по сумме в копейках.
func Category(amount int64) string {
	// TODO: Проверить границы через switch без выражения.
	switch {
	case 0 <= amount && amount <= 9999:
		return "small"
	case 10000 <= amount && amount <= 49999:
		return "regular"
	case 50000 <= amount:
		return "large"
	default:
		return "invalid"
	}
}

// RuneCount считает кодовые точки Unicode в строке через range.
func RuneCount(text string) int {
	// TODO: Не использовать len(text) как количество кодовых точек.
	var count int = 0
	for range text {
		count++
	}

	return count
}
