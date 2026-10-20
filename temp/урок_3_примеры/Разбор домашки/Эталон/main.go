package main

import "fmt"

// FizzBuzz возвращает FizzBuzz, Fizz, Buzz или пустую строку для обычного числа.
// Для n <= 0 возвращает invalid.
func FizzBuzz(n int) string {
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

// Sum возвращает сумму. Для пустого среза и nil возвращает 0.
func Sum(values []int64) int64 {
	var sum int64
	for _, value := range values {
		sum += value
	}
	return sum
}

// Min возвращает минимум. Для пустого среза и nil возвращает 0 по условию задачи.
func Min(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, value := range values[1:] {
		if value < min {
			min = value
		}
	}
	return min
}

// Max возвращает максимум. Для пустого среза и nil возвращает 0 по условию задачи.
func Max(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, value := range values[1:] {
		if value > max {
			max = value
		}
	}
	return max
}

// Average возвращает среднее как float64. Для пустого среза и nil возвращает 0.
func Average(values []int64) float64 {
	if len(values) == 0 {
		return 0
	}
	return float64(Sum(values)) / float64(len(values))
}

// Category возвращает invalid, small, regular или large по сумме в копейках.
func Category(amount int64) string {
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

// RuneCount считает кодовые точки Unicode в строке через range.
func RuneCount(text string) int {
	var count int
	for range text {
		count++
	}
	return count
}

func main() {
	fmt.Println("FizzBuzz: порядок условий")
	for _, n := range []int{0, -15, 15, 7} {
		fmt.Printf("FizzBuzz(%d) = %q\n", n, FizzBuzz(n))
	}

	fmt.Println("\nСтатистика: пустой ввод, один элемент, отрицательные числа")
	for _, values := range [][]int64{nil, {}, {42}, {-8, -2, -5}, {1, 2}} {
		fmt.Printf("%v: сумма=%d минимум=%d максимум=%d среднее=%g\n",
			values, Sum(values), Min(values), Max(values), Average(values))
	}

	values := []int64{1, 2}
	fmt.Println("\nПреобразование перед делением:", Average(values))
	fmt.Println("Преобразование после целочисленного деления:", float64(Sum(values)/int64(len(values))))

	values = []int64{9007199254740992, 1, -9007199254740992}
	fmt.Println("Среднее при компенсации больших чисел:", Average(values))

	fmt.Println("\nCategory: границы диапазонов")
	for _, amount := range []int64{-1, 0, 9999, 10000, 49999, 50000} {
		fmt.Printf("Category(%d) = %q\n", amount, Category(amount))
	}

	fmt.Println("\nRuneCount: байты и кодовые точки")
	for _, text := range []string{"Goшник", "🙂", "e\u0301"} {
		fmt.Printf("%q: байтов=%d кодовых точек=%d\n", text, len(text), RuneCount(text))
	}
}
