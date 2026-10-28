package homework

import "strings"

// PairInput соответствует объекту nums/target из условия контеста.
type PairInput struct {
	Nums   []int
	Target int
}

// longestPrefix возвращает самый длинный общий префикс строк.
func longestPrefix(words []string) string {
	// TODO: реализовать поиск префикса по кодовым точкам Unicode.
	if len(words) == 0 {
		return ""
	}
	sortedWords := make([]string, len(words))
	copy(sortedWords, words)
	stack := make([][2]int, 1, len(sortedWords))
	stack[0] = [2]int{0, len(sortedWords) - 1}

	for len(stack) != 0 {
		last_idx := len(stack) - 1
		borders := stack[last_idx]
		stack = stack[:last_idx]

		pivot_idx := (borders[1] + borders[0]) / 2
		pivot := sortedWords[pivot_idx]

		i, j := borders[0], borders[1]
		for i <= j {
			for sortedWords[i] < pivot {
				i += 1
			}
			for sortedWords[j] > pivot {
				j -= 1
			}
			if i <= j {
				sortedWords[i], sortedWords[j] = sortedWords[j], sortedWords[i]
				i += 1
				j -= 1
			}
		}

		if borders[0] < j {
			stack = append(stack, [2]int{borders[0], j})
		}

		if borders[1] > i {
			stack = append(stack, [2]int{i, borders[1]})
		}
	}

	s1 := []rune(sortedWords[0])
	s2 := []rune(sortedWords[len(sortedWords)-1])

	min_len := min(len(s1), len(s2))
	idx := 0
	for idx < min_len && s1[idx] == s2[idx] {
		idx++
	}

	return string(s1[:idx])
}

// pairIndices возвращает индексы двух разных элементов по возрастанию.
func pairIndices(input PairInput) [2]int {
	// TODO: найти единственную пару с суммой input.Target.
	m := make(map[int]int) //[требуемое дополнение]индекс в nums

	for idx, num := range input.Nums {
		m[input.Target-num] = idx
	}

	for ind, num := range input.Nums {
		idx, ok := m[num]
		if ok == true && ind != idx {
			return [2]int{ind, idx}
		}
	}

	return [2]int{}
}

// hasDuplicate сообщает, встречается ли какое-либо значение повторно.
func hasDuplicate(nums []int) bool {
	// TODO: проверить наличие повторов.
	m := make(map[int]int) //[key]value

	for _, num := range nums {
		_, ok := m[num]
		if ok == true {
			return true
		}
		m[num] = 1
	}

	return false
}

// compactSorted создаёт независимый слайс уникальных значений.
func compactSorted(nums []int) []int {
	// TODO: сохранить порядок и не изменять исходный слайс.
	if len(nums) == 0 {
		return []int{}
	}
	ans := make([]int, 1, len(nums))
	ans[0] = nums[0]
	ans_len := 0
	for _, vl := range nums {
		if ans[ans_len] != vl {
			ans = append(ans, vl)
			ans_len++
		}
	}
	return ans
}

// balancedBrackets проверяет пары и порядок закрытия скобок.
func balancedBrackets(text string) bool {
	// TODO: обработать строку из символов ()[]{}.
	if len(text) == 0 {
		return true
	}
	m := make(map[rune]rune)
	m[')'] = '('
	m['}'] = '{'
	m[']'] = '['
	s1 := []rune(text)
	stack := make([]rune, 0, len(text))
	for _, sym := range s1 {
		if strings.Contains("([{", string(sym)) {
			stack = append(stack, sym)
			continue
		}
		pair := m[sym]
		if len(stack) == 0 || stack[len(stack)-1] != pair {
			return false
		}
		stack = stack[:len(stack)-1]

	}
	if len(stack) != 0 {
		return false
	}
	return true
}
