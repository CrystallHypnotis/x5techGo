package service

import (
	"context"

	"slices"
	"strings"
	"time"

	"example.com/order-service/internal/order"
)

var statusTransitions = map[order.Status][]order.Status{
	order.StatusNew:        {order.StatusProcessing, order.StatusCancelled},
	order.StatusProcessing: {order.StatusCompleted, order.StatusCancelled},
}

func checkOrderStatusTransition(target order.Order, status order.Status) bool {
	allowed, ok := statusTransitions[target.Status]
	if !ok {
		return false
	}

	return slices.Contains(allowed, status)
}

func validateCreateOrderInput(input order.CreateOrderInput) error {
	if len(input.Items) == 0 {
		return order.ErrInvalidOrder
	}
	for index, _ := range input.Items {
		if input.Items[index].Price <= 0 || input.Items[index].Quantity <= 0 ||
			strings.TrimSpace(input.Items[index].Name) == "" {
			return order.ErrInvalidOrder
		}
	}
	return nil
}

func calcTotalSum(items []order.Item) int64 {
	var sum int64
	for index, _ := range items {
		sum += items[index].Price * int64(items[index].Quantity)
	}

	return sum
}

type OrderService struct {
	repository order.Repository
}

func NewOrderService(repository order.Repository) *OrderService {
	return &OrderService{repository: repository}
}

func (s *OrderService) Create(ctx context.Context, input order.CreateOrderInput) (order.Order, error) {
	// TODO: validate input, calculate total, set initial status and history,
	// then save the order through the repository interface.
	moment := time.Now()

	err := validateCreateOrderInput(input)
	if err != nil {
		return order.Order{}, order.ErrInvalidOrder
	}

	var preNewOrder = order.Order{
		Items:      input.Items,
		TotalPrice: calcTotalSum(input.Items),
		Status:     order.StatusNew,
		CreatedAt:  moment,
		History: []order.StatusChange{
			{
				Status:    order.StatusNew,
				ChangedAt: moment,
			},
		},
	}

	newOrder, err := s.repository.Create(ctx, preNewOrder)
	if err != nil {
		return order.Order{}, err
	}
	return newOrder, err
}

func (s *OrderService) GetByID(ctx context.Context, id int64) (order.Order, error) {
	// TODO: implement the scenario through the repository interface.
	result, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return order.Order{}, err
	}
	return result, nil
}

func (s *OrderService) List(ctx context.Context) ([]order.Order, error) {
	// TODO: implement the scenario through the repository interface.
	result, err := s.repository.List(ctx)
	if err != nil {
		return []order.Order{}, err
	}
	return result, nil
}

func (s *OrderService) UpdateStatus(ctx context.Context, id int64, status order.Status) (order.Order, error) {
	// TODO: validate the requested transition, update current status and history,
	// then persist the changed order.
	moment := time.Now()

	oldOrder, err := s.GetByID(ctx, id)
	if err != nil {
		return order.Order{}, err
	}

	if allowed := checkOrderStatusTransition(oldOrder, status); allowed == false {
		return order.Order{}, order.ErrInvalidStatusTransition
	}

	oldOrder.Status = status
	oldOrder.History = append(oldOrder.History,
		order.StatusChange{
			Status:    status,
			ChangedAt: moment,
		})

	newOrder, err := s.repository.Update(ctx, oldOrder)
	if err != nil {
		return order.Order{}, err
	}
	return newOrder, nil
}

func (s *OrderService) History(ctx context.Context, id int64) ([]order.StatusChange, error) {
	// TODO: return status changes in chronological order.
	target, err := s.GetByID(ctx, id)
	if err != nil {
		return []order.StatusChange{}, err
	}

	copyHistory := slices.Clone(target.History)

	return copyHistory, nil
}
