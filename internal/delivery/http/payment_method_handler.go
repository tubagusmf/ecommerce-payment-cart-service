package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tubagusmf/ecommerce-payment-cart-service/internal/model"
)

type PaymentMethodHandler struct {
	paymentMethodUsecase model.IPaymentMethodUsecase
}

func NewPaymentMethodHandler(e *echo.Echo, paymentMethodUsecase model.IPaymentMethodUsecase) {
	handler := &PaymentMethodHandler{
		paymentMethodUsecase: paymentMethodUsecase,
	}

	route := e.Group("/v1/payment-methods")
	route.GET("", handler.FindAll)
	route.GET("/:id", handler.FindByID)
	route.POST("/create", handler.Create)
	route.PUT("/update/:id", handler.Update)
	route.DELETE("/delete/:id", handler.Delete)
}

func (h *PaymentMethodHandler) FindAll(c echo.Context) error {
	paymentMethods, err := h.paymentMethodUsecase.FindAll(c.Request().Context(), model.PaymentMethod{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   paymentMethods,
	})
}

func (h *PaymentMethodHandler) FindByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	paymentMethod, err := h.paymentMethodUsecase.FindByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if paymentMethod == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Payment method not found")
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   paymentMethod,
	})
}

func (h *PaymentMethodHandler) Create(c echo.Context) error {
	var body model.CreatePaymentMethod
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	err := h.paymentMethodUsecase.Create(c.Request().Context(), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, Response{
		Status:  http.StatusCreated,
		Message: "Payment method created successfully",
	})
}

func (h *PaymentMethodHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	var body model.UpdatePaymentMethod
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	err = h.paymentMethodUsecase.Update(c.Request().Context(), id, body)
	if err != nil {
		if err.Error() == "payment method not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Payment method not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Payment method updated successfully",
	})
}

func (h *PaymentMethodHandler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = h.paymentMethodUsecase.Delete(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "payment method not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Payment method not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Payment method deleted successfully",
	})
}
