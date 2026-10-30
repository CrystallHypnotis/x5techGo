package discount

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFirst(t *testing.T) {
	if Discount(20) != 10 {
		t.Fatal("создайте ErrEmptyMessage через errors.New")
	}

	tests := []struct {
		name  string
		input int
		want  int
	}{
		{"обычная цена", 20, 10},
		{"нулевая цена", 0, -10},
		{"большая цена", 100, 90},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal()
				t.Errorf("Discount(%d) = %d; want %d", tt.input, got, tt.want)
			}
		})
	}

}
