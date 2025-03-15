package model

import (
	"context"
	"time"

	pb "github.com/tubagusmf/ecommerce-payment-cart-service/pb/payment_service"
)

// var PaymentMethodMapping = map[pb.PaymentMethod]PaymentMethod{
// 	pb.PaymentMethod_BANK_TRANSFER: {ID: 4, Name: "Bank Transfer", BankCode: "BT"},
// }

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

type IPaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	FindAll(ctx context.Context, payment Payment) ([]*Payment, error)
	FindById(ctx context.Context, id int64) (*Payment, error)
	FindByOrderID(ctx context.Context, orderID string) (*Payment, error)
	UpdateStatus(ctx context.Context, orderID string, status PaymentStatus) error
	FindPaymentMethodByID(ctx context.Context, id int64) (*PaymentMethod, error)
}

type IPaymentUsecase interface {
	ProcessPayment(ctx context.Context, orderID string, userID int64, paymentMethod PaymentMethod, paymentStatus PaymentStatus) (*Payment, error)
	ConfirmPayment(ctx context.Context, orderID string) error
	GetPaymentStatus(ctx context.Context, paymentID int64) (*Payment, error)
	GetPayments(ctx context.Context, payment Payment) ([]*Payment, error)
	GetPaymentMethodByID(ctx context.Context, methodID int64) (*PaymentMethod, error)
	GetPaymentByID(ctx context.Context, id int64) (*Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID string) (*Payment, error)
	MarkPaymentPaid(ctx context.Context, id string) error
}

type PaymentMethod struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	BankCode  string     `json:"bank_code"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type Payment struct {
	ID              int64         `json:"id"`
	OrderID         string        `json:"order_id"`
	UserID          int64         `json:"user_id"`
	PaymentMethodID int64         `json:"payment_method_id" gorm:"foreignKey:PaymentMethodID;references:ID"`
	PaymentMethod   PaymentMethod `json:"payment_method" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Status          PaymentStatus `json:"status"`
	TransactionID   string        `json:"transaction_id,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
}

type CreatePaymentMethod struct {
	Name     string `json:"name" validate:"required"`
	BankCode string `json:"bank_code" validate:"required"`
}

type UpdatePaymentMethod struct {
	Name     string `json:"name" validate:"required"`
	BankCode string `json:"bank_code" validate:"required"`
}

type ProcessPaymentInput struct {
	OrderID         string `json:"order_id" validate:"required"`
	UserID          int64  `json:"user_id" validate:"required"`
	PaymentMethodID int64  `json:"payment_method_id" validate:"required"`
	PaymentStatus   string `json:"payment_status" validate:"required"`
}

func ModelToProtoPaymentStatus(status PaymentStatus) pb.PaymentStatus {
	switch status {
	case StatusPending:
		return pb.PaymentStatus_PAYMENT_STATUS_PENDING
	case StatusSuccess:
		return pb.PaymentStatus_PAYMENT_STATUS_SUCCESS
	case StatusFailed:
		return pb.PaymentStatus_PAYMENT_STATUS_FAILED
	default:
		return pb.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}
}

// func ProtoToModelPaymentMethod(protoMethod pb.PaymentMethod) PaymentMethod {
// 	if method, exists := PaymentMethodMapping[protoMethod]; exists {
// 		return method
// 	}
// 	return PaymentMethod{} // Default jika tidak ditemukan
// }

// // Fungsi untuk mengonversi dari model ke proto
// func ModelToProtoPaymentMethod(method PaymentMethod) pb.PaymentMethod {
// 	for protoEnum, modelMethod := range PaymentMethodMapping {
// 		if modelMethod.ID == method.ID {
// 			return protoEnum
// 		}
// 	}
// 	return pb.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
// }
