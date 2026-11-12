package service_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"example.com/order-service/internal/order"
	"example.com/order-service/internal/service"
	"github.com/stretchr/testify/assert"
)

// func TestOrderService(t *testing.T) {
// 	t.Skip("TODO: replace this placeholder with table-driven unit tests")
// }
//
//
// type Repository interface {
// 	Create(ctx context.Context, value Order) (Order, error)
// 	GetByID(ctx context.Context, id int64) (Order, error)
// 	List(ctx context.Context) ([]Order, error)
// 	Update(ctx context.Context, value Order) (Order, error)
// }

type blankRepository struct {
	items  []order.Order
	nextID int64
}

func (r *blankRepository) Create(_ context.Context, value order.Order) (order.Order, error) {
	value.ID = r.nextID
	r.nextID++
	r.items = append(r.items, value)

	return value, nil
}

func (r *blankRepository) GetByID(_ context.Context, id int64) (order.Order, error) {
	for _, item := range r.items {
		if item.ID == id {
			return item, nil
		}
	}
	return order.Order{}, order.ErrOrderNotFound
}

func (r *blankRepository) List(_ context.Context) ([]order.Order, error) {
	itemsCopy := r.items
	return itemsCopy, nil
}

func (r *blankRepository) Update(_ context.Context, value order.Order) (order.Order, error) {
	for index, _ := range r.items {
		if r.items[index].ID == value.ID {
			r.items[index] = value
			return r.items[index], nil
		}
	}
	return order.Order{}, order.ErrOrderNotFound
}

func TestCreate(tt *testing.T) {
	tt.Parallel()

	type test struct {
		name        string
		context     context.Context
		input       order.CreateOrderInput
		expectedVal order.Order
		expectedErr error
	}

	valid := test{
		name:    "успешное создание заказа",
		context: context.Background(),
		input: order.CreateOrderInput{
			Items: []order.Item{
				{Name: "Monster White", Quantity: 999, Price: 2},
				{Name: "Sleeping hours", Quantity: 9, Price: 999999},
			},
		},
		expectedVal: order.Order{
			ID: 1,
			Items: []order.Item{
				{Name: "Monster White", Quantity: 999, Price: 2},
				{Name: "Sleeping hours", Quantity: 9, Price: 999999},
			},
			TotalPrice: 9001989,
			Status:     order.StatusNew,
			CreatedAt:  time.Time{},
			History: []order.StatusChange{
				{
					Status:    order.StatusNew,
					ChangedAt: time.Time{},
				},
			},
		},
		expectedErr: nil,
	}

	noItems := valid
	noItems.name = "заказ без позиций"
	noItems.input.Items = nil
	noItems.expectedVal = order.Order{}
	noItems.expectedErr = order.ErrInvalidOrder

	emptyName := valid
	emptyName.name = "заказ с пустым названием товара"
	emptyName.input.Items = slices.Clone(valid.input.Items)
	emptyName.input.Items[0].Name = ""
	emptyName.expectedVal = order.Order{}
	emptyName.expectedErr = order.ErrInvalidOrder

	zeroPrice := valid
	zeroPrice.name = "заказ с нулевой ценой"
	zeroPrice.input.Items = slices.Clone(valid.input.Items)
	zeroPrice.input.Items[0].Price = 0
	zeroPrice.expectedVal = order.Order{}
	zeroPrice.expectedErr = order.ErrInvalidOrder

	negativePrice := valid
	negativePrice.name = "заказ с отрицательной ценой"
	negativePrice.input.Items = slices.Clone(valid.input.Items)
	negativePrice.input.Items[0].Price = -1
	negativePrice.expectedVal = order.Order{}
	negativePrice.expectedErr = order.ErrInvalidOrder

	zeroQuantity := valid
	zeroQuantity.name = "заказ с нулевым количеством"
	zeroQuantity.input.Items = slices.Clone(valid.input.Items)
	zeroQuantity.input.Items[0].Quantity = 0
	zeroQuantity.expectedVal = order.Order{}
	zeroQuantity.expectedErr = order.ErrInvalidOrder

	negativeQuantity := valid
	negativeQuantity.name = "заказ с отрицательным количеством"
	negativeQuantity.input.Items = slices.Clone(valid.input.Items)
	negativeQuantity.input.Items[0].Quantity = -1
	negativeQuantity.expectedVal = order.Order{}
	negativeQuantity.expectedErr = order.ErrInvalidOrder

	tests := []test{
		valid,
		noItems,
		emptyName,
		zeroPrice,
		negativePrice,
		zeroQuantity,
		negativeQuantity,
	}

	for _, subTest := range tests {
		tt.Run(subTest.name, func(t *testing.T) {
			t.Parallel()
			srvc := service.NewOrderService(&blankRepository{nextID: 1})

			result, err := srvc.Create(subTest.context, subTest.input)

			result.CreatedAt = time.Time{}
			for index := range result.History {
				result.History[index].ChangedAt = time.Time{}
			}

			assert.ErrorIs(t, err, subTest.expectedErr)
			assert.Equal(t, subTest.expectedVal, result)
		})
	}
}

func TestGetByID(tt *testing.T) {
	tt.Parallel()
	type test struct {
		name        string
		id          int64
		expectedVal order.Order
		expectedErr error
	}

	stored := order.Order{
		ID: 1,
		Items: []order.Item{{
			Name:     "Monster pink",
			Quantity: 99,
			Price:    100,
		}},
		TotalPrice: 9900,
		Status:     order.StatusNew,
	}

	tests := []test{
		{
			name:        "получение существующего заказа",
			id:          1,
			expectedVal: stored,
			expectedErr: nil,
		},
		{
			name:        "получение отсутствующего заказа",
			id:          999,
			expectedVal: order.Order{},
			expectedErr: order.ErrOrderNotFound,
		},
	}

	for _, subTest := range tests {
		tt.Run(subTest.name, func(t *testing.T) {
			t.Parallel()

			repository := &blankRepository{
				items:  []order.Order{stored},
				nextID: 2,
			}
			srvc := service.NewOrderService(repository)
			result, err := srvc.GetByID(context.Background(), subTest.id)

			assert.ErrorIs(t, err, subTest.expectedErr)
			assert.Equal(t, subTest.expectedVal, result)
		})
	}
}

func TestList(tt *testing.T) {
	tt.Parallel()

	type test struct {
		name        string
		items       []order.Order
		expectedVal []order.Order
		expectedErr error
	}

	stored := []order.Order{
		{
			ID: 1,
			Items: []order.Item{{
				Name:     "Monster pink",
				Quantity: 99,
				Price:    100,
			}},
			TotalPrice: 9900,
			Status:     order.StatusNew,
		},
		{
			ID: 2,
			Items: []order.Item{{
				Name:     "Sleeping hours",
				Quantity: 1,
				Price:    99999999,
			}},
			TotalPrice: 99999999,
			Status:     order.StatusNew,
		},
	}

	tests := []test{
		{
			name:        "получение списка заказов",
			items:       stored,
			expectedVal: stored,
			expectedErr: nil,
		},
	}

	for _, subTest := range tests {
		tt.Run(subTest.name, func(t *testing.T) {
			t.Parallel()

			repository := &blankRepository{
				items:  subTest.items,
				nextID: 1,
			}
			srvc := service.NewOrderService(repository)
			result, err := srvc.List(context.Background())

			assert.ErrorIs(t, err, subTest.expectedErr)
			assert.Equal(t, subTest.expectedVal, result)
		})
	}
}

func TestUpdateStatus(tt *testing.T) {
	tt.Parallel()

	type test struct {
		name            string
		id              int64
		status          order.Status
		expectedVal     order.Order
		expectedHistory []order.Status
		expectedErr     error
	}

	newOrder := order.Order{
		ID:     1,
		Status: order.StatusNew,
		History: []order.StatusChange{
			{Status: order.StatusNew},
		},
	}
	processingOrder := order.Order{
		ID:     2,
		Status: order.StatusProcessing,
		History: []order.StatusChange{
			{Status: order.StatusNew},
			{Status: order.StatusProcessing},
		},
	}
	completedOrder := order.Order{
		ID:     3,
		Status: order.StatusCompleted,
		History: []order.StatusChange{
			{Status: order.StatusNew},
			{Status: order.StatusProcessing},
			{Status: order.StatusCompleted},
		},
	}
	cancelledOrder := order.Order{
		ID:     4,
		Status: order.StatusCancelled,
		History: []order.StatusChange{
			{Status: order.StatusNew},
			{Status: order.StatusCancelled},
		},
	}

	tests := []test{
		{
			name:            "допустимый переход new к processing",
			id:              1,
			status:          order.StatusProcessing,
			expectedVal:     order.Order{ID: 1, Status: order.StatusProcessing},
			expectedHistory: []order.Status{order.StatusNew, order.StatusProcessing},
			expectedErr:     nil,
		},
		{
			name:            "допустимый переход new к cancelled",
			id:              1,
			status:          order.StatusCancelled,
			expectedVal:     order.Order{ID: 1, Status: order.StatusCancelled},
			expectedHistory: []order.Status{order.StatusNew, order.StatusCancelled},
			expectedErr:     nil,
		},
		{
			name:            "допустимый переход processing к completed",
			id:              2,
			status:          order.StatusCompleted,
			expectedVal:     order.Order{ID: 2, Status: order.StatusCompleted},
			expectedHistory: []order.Status{order.StatusNew, order.StatusProcessing, order.StatusCompleted},
			expectedErr:     nil,
		},
		{
			name:            "допустимый переход processing к cancelled",
			id:              2,
			status:          order.StatusCancelled,
			expectedVal:     order.Order{ID: 2, Status: order.StatusCancelled},
			expectedHistory: []order.Status{order.StatusNew, order.StatusProcessing, order.StatusCancelled},
			expectedErr:     nil,
		},
		{
			name:        "недопустимый переход new к completed",
			id:          1,
			status:      order.StatusCompleted,
			expectedVal: order.Order{},
			expectedErr: order.ErrInvalidStatusTransition,
		},
		{
			name:        "запрет изменения завершённого заказа",
			id:          3,
			status:      order.StatusCancelled,
			expectedVal: order.Order{},
			expectedErr: order.ErrInvalidStatusTransition,
		},
		{
			name:        "запрет изменения отменённого заказа",
			id:          4,
			status:      order.StatusCompleted,
			expectedVal: order.Order{},
			expectedErr: order.ErrInvalidStatusTransition,
		},
		{
			name:        "заказ не найден",
			id:          999,
			status:      order.StatusProcessing,
			expectedVal: order.Order{},
			expectedErr: order.ErrOrderNotFound,
		},
	}

	for _, subTest := range tests {
		tt.Run(subTest.name, func(t *testing.T) {
			t.Parallel()
			repository := &blankRepository{
				items: []order.Order{
					newOrder,
					processingOrder,
					completedOrder,
					cancelledOrder,
				},
				nextID: 5,
			}
			srvc := service.NewOrderService(repository)
			result, err := srvc.UpdateStatus(context.Background(), subTest.id, subTest.status)

			assert.ErrorIs(t, err, subTest.expectedErr)
			assert.Equal(t, subTest.expectedVal.ID, result.ID)
			assert.Equal(t, subTest.expectedVal.Status, result.Status)

			var actualHistory []order.Status
			for _, change := range result.History {
				actualHistory = append(actualHistory, change.Status)
			}
			assert.Equal(t, subTest.expectedHistory, actualHistory)
		})
	}
}

func TestHistory(tt *testing.T) {
	tt.Parallel()

	type test struct {
		name        string
		id          int64
		expectedVal []order.StatusChange
		expectedErr error
	}

	stored := order.Order{
		ID:     1,
		Status: order.StatusCompleted,
		History: []order.StatusChange{
			{Status: order.StatusNew, ChangedAt: time.Unix(67, 0)},
			{Status: order.StatusProcessing, ChangedAt: time.Unix(667, 0)},
			{Status: order.StatusCompleted, ChangedAt: time.Unix(6668, 0)},
		},
	}

	tests := []test{
		{
			name:        "получение истории заказа",
			id:          1,
			expectedVal: stored.History,
			expectedErr: nil,
		},
		{
			name:        "история отсутствующего заказа",
			id:          999,
			expectedVal: []order.StatusChange{},
			expectedErr: order.ErrOrderNotFound,
		},
	}

	for _, subTest := range tests {
		tt.Run(subTest.name, func(t *testing.T) {
			t.Parallel()
			repository := &blankRepository{
				items:  []order.Order{stored},
				nextID: 2,
			}
			srvc := service.NewOrderService(repository)
			result, err := srvc.History(context.Background(), subTest.id)

			assert.ErrorIs(t, err, subTest.expectedErr)
			assert.Equal(t, subTest.expectedVal, result)

			for i := 1; i < len(result); i++ {
				assert.True(t, result[i-1].ChangedAt.Before(result[i].ChangedAt))
			}
		})
	}
}
