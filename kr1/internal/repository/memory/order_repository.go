package memory

import (
	"context"
	"slices"

	"example.com/order-service/internal/order"
)

type OrderRepository struct {
	items  []order.Order
	nextID int64
}

func copyOrder(value order.Order) order.Order {
	value.Items = slices.Clone(value.Items)
	value.History = slices.Clone(value.History)
	return value
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		items:  make([]order.Order, 0),
		nextID: 1,
	}
}

func (r *OrderRepository) Create(_ context.Context, value order.Order) (order.Order, error) {
	// TODO: assign an ID and store the order.
	value.ID = r.nextID
	r.nextID++
	r.items = append(r.items, copyOrder(value))

	return copyOrder(value), nil
}

func (r *OrderRepository) GetByID(_ context.Context, id int64) (order.Order, error) {
	// TODO: return ErrOrderNotFound when the order does not exist.
	for _, item := range r.items {
		if item.ID == id {
			return copyOrder(item), nil
		}
	}
	return order.Order{}, order.ErrOrderNotFound
}

func (r *OrderRepository) List(_ context.Context) ([]order.Order, error) {
	// TODO: return all stored orders in creation order without exposing mutable internal state.
	itemsCopy := make([]order.Order, len(r.items))
	for index := range r.items {
		itemsCopy[index] = copyOrder(r.items[index])
	}
	return itemsCopy, nil
}

func (r *OrderRepository) Update(_ context.Context, value order.Order) (order.Order, error) {
	// TODO: replace an existing order or return ErrOrderNotFound.
	for index, _ := range r.items {
		if r.items[index].ID == value.ID {
			r.items[index] = copyOrder(value)
			return copyOrder(r.items[index]), nil
		}
	}
	return order.Order{}, order.ErrOrderNotFound
}
