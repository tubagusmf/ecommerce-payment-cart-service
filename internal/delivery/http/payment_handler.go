package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tubagusmf/ecommerce-payment-cart-service/internal/model"
)

type PaymentHttpHandler struct {
	paymentUsecase       model.IPaymentUsecase
	paymentMethodUsecase model.IPaymentMethodUsecase
}

func NewPaymentHttpHandler(e *echo.Echo, paymentUsecase model.IPaymentUsecase, paymentMethodUsecase model.IPaymentMethodUsecase) {
	handler := &PaymentHttpHandler{
		paymentUsecase:       paymentUsecase,
		paymentMethodUsecase: paymentMethodUsecase,
	}

	routePayment := e.Group("v1/payments")
	routePayment.POST("/create", handler.ProcessPayment)
	routePayment.GET("/", handler.GetPayments)
	routePayment.GET("/:id", handler.GetPaymentByID)
	routePayment.GET("/order/:id", handler.GetPaymentByOrderID)
}

func (h *PaymentHttpHandler) ProcessPayment(c echo.Context) error {
	var req model.ProcessPaymentInput

	if err := c.Bind(&req); err != nil {
		log.Println("Error binding request:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	paymentMethod, err := h.paymentMethodUsecase.FindByID(c.Request().Context(), req.PaymentMethodID)
	if err != nil || paymentMethod == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payment method"})
	}

	var paymentStatus model.PaymentStatus
	switch req.PaymentStatus {
	case "success":
		paymentStatus = model.StatusSuccess
	case "failed":
		paymentStatus = model.StatusFailed
	default:
		paymentStatus = model.StatusPending
	}

	createdPayment, err := h.paymentUsecase.ProcessPayment(
		c.Request().Context(),
		req.OrderID,
		req.UserID,
		*paymentMethod,
		paymentStatus,
	)
	if err != nil {
		log.Println("Error processing payment:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process payment"})
	}

	return c.JSON(http.StatusOK, createdPayment)
}

func (h *PaymentHttpHandler) GetPayments(c echo.Context) error {
	payments, err := h.paymentUsecase.GetPayments(c.Request().Context(), model.Payment{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get payments"})
	}

	return c.JSON(http.StatusOK, payments)
}

func (h *PaymentHttpHandler) GetPaymentByID(c echo.Context) error {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payment ID"})
	}

	payment, err := h.paymentUsecase.GetPaymentByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Payment not found"})
	}

	return c.JSON(http.StatusOK, payment)
}

func (h *PaymentHttpHandler) GetPaymentByOrderID(c echo.Context) error {
	idStr := c.Param("id")

	payment, err := h.paymentUsecase.GetPaymentByOrderID(c.Request().Context(), idStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Payment not found"})
	}

	return c.JSON(http.StatusOK, payment)
}
