package discount

const discountAmount = 10

func Discount(price int) int {
	return price - discountAmount
}
