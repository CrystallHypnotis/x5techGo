package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	_ "example.com/order-service/docs"
	"example.com/order-service/internal/order"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type OrderService interface {
	Create(ctx context.Context, input order.CreateOrderInput) (order.Order, error)
	GetByID(ctx context.Context, id int64) (order.Order, error)
	List(ctx context.Context) ([]order.Order, error)
	UpdateStatus(ctx context.Context, id int64, status order.Status) (order.Order, error)
	History(ctx context.Context, id int64) ([]order.StatusChange, error)
}

type Handler struct {
	service OrderService
	logger  *slog.Logger
}

func NewRouter(service OrderService, logger *slog.Logger) http.Handler {
	h := &Handler{service: service, logger: logger}
	router := chi.NewRouter()
	router.Use(recoveryMiddleware(logger))
	router.Use(requestLoggingMiddleware(logger))

	router.Get("/health", h.health)
	router.Route("/api/v1/orders", func(r chi.Router) {
		r.Post("/", h.createOrder)
		r.Get("/", h.getOrdersList)
		r.Get("/{id}", h.getOrder)
		r.Patch("/{id}/status", h.updateOrder)
		r.Get("/{id}/history", h.getOrderHistory)
	})

	router.Get("/swagger/*", httpSwagger.Handler())

	return router
}

// @Summary		Проверка статуса сервиса
// @Description	Возвращает статус сервиса
// @Tags		health
// @Produce		json
// @Success		200	{object} map[string]string
// @Router		/health [get]
func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// @Summary		Создание заказа в сервисе
// @Description	Принимает список товаров, валидирует их, после передаёт данные в repository
// @Tags		orders
// @Accept		json
// @Produce		json
// @Param		items body order.CreateOrderInput true "Товары в заказе"
// @Success		201	{object} order.Order
// @Failure		400	{object} errorResponse
// @Failure		500	{object} errorResponse
// @Router		/api/v1/orders [post]
func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	var input order.CreateOrderInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "invalid JSON",
			Code:  "INVALID_JSON",
		})
		return
	}

	newOrder, err := h.service.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, order.ErrInvalidOrder) {
			writeJSON(w, http.StatusBadRequest, errorResponse{
				Error: "invalid order",
				Code:  "INVALID_ORDER",
			})
			return
		}

		h.logger.Error("failed to create order", slog.String("error", err.Error()))
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Error: "internal error",
			Code:  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	writeJSON(w, http.StatusCreated, newOrder)
}

// @Summary		Получение списка заказов
// @Description	Возвращает все заказы в repository
// @Tags		orders
// @Produce		json
// @Success		200	{array}	order.Order
// @Failure		500	{object} errorResponse
// @Router		/api/v1/orders [get]
func (h *Handler) getOrdersList(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.List(r.Context())
	if err != nil {
		h.logger.Error("failed to get orders list", slog.String("error", err.Error()))
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Error: "internal error",
			Code:  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	writeJSON(w, http.StatusOK, orders)
}

// @Summary		Получение заказа
// @Description	Возвращает заказ по его ID
// @Tags		orders
// @Produce		json
// @Param		id path int64 true "ID заказа"
// @Success		200	{object} order.Order
// @Failure		400	{object} errorResponse
// @Failure		404	{object} errorResponse
// @Failure		500	{object} errorResponse
// @Router		/api/v1/orders/{id} [get]
func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "invalid order id",
			Code:  "INVALID_ORDER_ID",
		})
		return
	}

	result, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, order.ErrOrderNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{
				Error: "order not found",
				Code:  "ORDER_NOT_FOUND",
			})
			return
		}

		h.logger.Error("failed to get order", slog.Int64("order_id", id), slog.String("error", err.Error()))
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Error: "internal error",
			Code:  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// @Summary		Изменение статуса заказа
// @Description	Измененяет статус заказа по ID с валидацией допустимости транзита статуса
// @Tags		orders
// @Accept		json
// @Produce		json
// @Param		id path	int64 true "ID заказа"
// @Param		status body	order.UpdateStatusInput	true "Новый статус"
// @Success		200	{object} order.Order
// @Failure		400	{object} errorResponse
// @Failure		404	{object} errorResponse
// @Failure		409	{object} errorResponse
// @Failure		500	{object} errorResponse
// @Router		/api/v1/orders/{id}/status [patch]
func (h *Handler) updateOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "invalid order id",
			Code:  "INVALID_ORDER_ID",
		})
		return
	}

	var input order.UpdateStatusInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "invalid order",
			Code:  "INVALID_ORDER",
		})
		return
	}

	newOrder, err := h.service.UpdateStatus(r.Context(), id, input.Status)
	if err != nil {
		if errors.Is(err, order.ErrOrderNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{
				Error: "order not found",
				Code:  "ORDER_NOT_FOUND",
			})
			return
		}

		if errors.Is(err, order.ErrInvalidStatus) {
			writeJSON(w, http.StatusBadRequest, errorResponse{
				Error: "invalid status",
				Code:  "INVALID_STATUS",
			})
			return
		}
		if errors.Is(err, order.ErrInvalidStatusTransition) {
			writeJSON(w, http.StatusConflict, errorResponse{
				Error: "invalid status transition",
				Code:  "INVALID_STATUS_TRANSITION",
			})
			return
		}

		h.logger.Error("failed to update order", slog.Int64("order_id", id), slog.String("error", err.Error()))
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Error: "internal error",
			Code:  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	writeJSON(w, http.StatusOK, newOrder)
}

// @Summary		Получение истории статусов заказа
// @Description	Возвращает историю изменения статусов заказа с временными метками
// @Tags		orders
// @Produce		json
// @Param		id	path int64 true	"ID заказа"
// @Success		200	{array}	order.StatusChange
// @Failure		400	{object} errorResponse
// @Failure		404	{object} errorResponse
// @Failure		500	{object} errorResponse
// @Router		/api/v1/orders/{id}/history [get]
func (h *Handler) getOrderHistory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "invalid order id",
			Code:  "INVALID_ORDER_ID",
		})
		return
	}

	result, err := h.service.History(r.Context(), id)
	if err != nil {
		if errors.Is(err, order.ErrOrderNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{
				Error: "order not found",
				Code:  "ORDER_NOT_FOUND",
			})
			return
		}

		h.logger.Error("failed to get order's history", slog.Int64("order_id", id), slog.String("error", err.Error()))
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Error: "internal error",
			Code:  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) notImplemented(w http.ResponseWriter, _ *http.Request) {
	// TODO: replace this placeholder with separate handlers for every endpoint.
	writeJSON(w, http.StatusNotImplemented, errorResponse{
		Error: "endpoint is not implemented",
		Code:  "NOT_IMPLEMENTED",
	})
}

type errorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
