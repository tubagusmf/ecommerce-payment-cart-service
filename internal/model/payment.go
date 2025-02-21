package model

import (
	"context"
	"time"
)

type PaymentStatus string

const (
	StatusPending PaymentStatus = "pending"
	StatusSuccess PaymentStatus = "success"
	StatusFailed  PaymentStatus = "failed"
)

type IPaymentMethodRepository interface {
	FindAll(ctx context.Context, paymentMethod PaymentMethod) ([]*PaymentMethod, error)
	FindByID(ctx context.Context, id int64) (*PaymentMethod, error)
	Create(ctx context.Context, paymentMethod PaymentMethod) error
	Update(ctx context.Context, paymentMethod PaymentMethod) error
	Delete(ctx context.Context, id int64) error
}

type IPaymentMethodUsecase interface {
	FindAll(ctx context.Context, paymentMethod PaymentMethod) ([]*PaymentMethod, error)
	FindByID(ctx context.Context, id int64) (*PaymentMethod, error)
	Create(ctx context.Context, in CreatePaymentMethod) error
	Update(ctx context.Context, id int64, in UpdatePaymentMethod) error
	Delete(ctx context.Context, id int64) error
}

type PaymentMethod struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	BankCode  string     `json:"bank_code"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type Payment struct {
	ID            int           `json:"id"`
	OrderID       string        `json:"order_id"`
	UserID        int           `json:"user_id"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Status        PaymentStatus `json:"status"`
	TransactionID string        `json:"transaction_id,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
}

type CreatePaymentMethod struct {
	Name     string `json:"name" validate:"required"`
	BankCode string `json:"bank_code" validate:"required"`
}

type UpdatePaymentMethod struct {
	Name     string `json:"name" validate:"required"`
	BankCode string `json:"bank_code" validate:"required"`
}
