package discount

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePercent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		percent     int
		expectedErr error
	}{
		{name: "`0` — допустимо", percent: 0},
		{name: "`10` — допустимо", percent: 10},
		{name: "`100` — допустимо", percent: 100},
		{name: "`-1` — ошибка", percent: -1, expectedErr: ErrInvalidPercent},
		{name: "`101` — ошибка", percent: 101, expectedErr: ErrInvalidPercent},
	}

	for _, subTest := range tests {
		t.Run(subTest.name, func(t *testing.T) {
			t.Parallel()
			assert.ErrorIs(t, ValidatePercent(subTest.percent), subTest.expectedErr)
		})
	}
}

func TestPercent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		price         int
		percent       int
		expectedValue int
		expectedErr   error
	}{
		{name: "обычная скидка 25", price: 100, percent: 25, expectedValue: 75},
		{name: "скидка 0%", price: 67, percent: 0, expectedValue: 67},
		{name: "скидка 100%", price: 67, percent: 100, expectedValue: 0},
		{name: "отрицательная цена — ошибка", price: -1, percent: 10, expectedErr: ErrNegativePrice},
		{name: "процент меньше 0 — ошибка", price: 67, percent: -1, expectedErr: ErrInvalidPercent},
		{name: "процент больше 100 — ошибка", price: 67, percent: 101, expectedErr: ErrInvalidPercent},
	}

	for _, subTest := range tests {
		t.Run(subTest.name, func(t *testing.T) {
			t.Parallel()
			value, err := Apply(subTest.price, subTest.percent)
			assert.ErrorIs(t, err, subTest.expectedErr)
			assert.Equal(t, subTest.expectedValue, value)
		})
	}
}

func TestFinalPrice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		items         []Item
		expectedValue int
		expectedErr   error
	}{
		{name: "nil список", items: nil, expectedValue: 0},
		{name: "пустой список", items: []Item{}, expectedValue: 0},
		{
			name: "несколько корректных товаров",
			items: []Item{
				{Price: 100, DiscountPercent: 10},
				{Price: 200, DiscountPercent: 25},
			},
			expectedValue: 240,
		},
		{
			name: "товар с некорректной ценой",
			items: []Item{
				{Price: 67, DiscountPercent: 67},
				{Price: -1, DiscountPercent: 0},
			},
			expectedErr: ErrNegativePrice,
		},
		{
			name:        "товар с некорректной скидкой",
			items:       []Item{{Price: 11, DiscountPercent: 101}},
			expectedErr: ErrInvalidPercent,
		},
	}

	for _, subTest := range tests {
		t.Run(subTest.name, func(t *testing.T) {
			t.Parallel()
			value, err := FinalPrice(subTest.items)
			assert.ErrorIs(t, err, subTest.expectedErr)
			assert.Equal(t, subTest.expectedValue, value)
		})
	}
}
