package repository

import (
	"context"
	"errors"

	"github.com/tubagusmf/ecommerce-payment-cart-service/internal/model"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) model.IPaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	if payment.PaymentMethodID == 0 {
		return errors.New("payment method ID is required")
	}
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *PaymentRepository) FindAll(ctx context.Context, payment model.Payment) ([]*model.Payment, error) {
	var payments []*model.Payment

	err := r.db.WithContext(ctx).
		Preload("PaymentMethod").
		Where(&payment).
		Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepository) FindById(ctx context.Context, id int64) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.WithContext(ctx).
		Preload("PaymentMethod").
		Where("id = ?", id).
		First(&payment).Error

	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) FindPaymentMethodByID(ctx context.Context, id int64) (*model.PaymentMethod, error) {
	var paymentMethod model.PaymentMethod
	err := r.db.Where("id = ?", id).First(&paymentMethod).Error
	if err != nil {
		return nil, err
	}
	return &paymentMethod, nil
}

func (r *PaymentRepository) FindByOrderID(ctx context.Context, orderID string) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.WithContext(ctx).
		Preload("PaymentMethod").
		Where("order_id = ?", orderID).
		First(&payment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, orderID string, status model.PaymentStatus) error {
	return r.db.WithContext(ctx).Model(&model.Payment{}).Where("order_id = ?", orderID).Update("status", status).Error
}
