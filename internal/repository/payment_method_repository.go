package repository

import (
	"context"
	"errors"
	"log"

	"github.com/tubagusmf/ecommerce-payment-cart-service/internal/model"

	"gorm.io/gorm"
)

type PaymentMethodRepository struct {
	db *gorm.DB
}

func NewPaymentMethodRepo(db *gorm.DB) model.IPaymentMethodRepository {
	return &PaymentMethodRepository{db: db}
}

func (r *PaymentMethodRepository) FindAll(ctx context.Context, paymentMethod model.PaymentMethod) ([]*model.PaymentMethod, error) {
	var paymentMethods []*model.PaymentMethod
	err := r.db.WithContext(ctx).
		Where(&paymentMethod).
		Where("deleted_at IS NULL").
		Find(&paymentMethods).Error
	return paymentMethods, err
}

func (r *PaymentMethodRepository) FindByID(ctx context.Context, id int64) (*model.PaymentMethod, error) {
	var paymentMethod model.PaymentMethod
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&paymentMethod).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Payment method with ID", id, "not found in database")
			return nil, err
		}
		log.Println("Error retrieving payment method:", err)
		return nil, err
	}

	log.Println("Payment method found:", paymentMethod)
	return &paymentMethod, nil
}

func (r *PaymentMethodRepository) Create(ctx context.Context, paymentMethod model.PaymentMethod) error {
	err := r.db.WithContext(ctx).Create(&paymentMethod).Error
	if err != nil {
		log.Println("Error inserting payment method:", err)
		return err
	}
	log.Println("Successfully inserted payment method:", paymentMethod)
	return nil
}

func (r *PaymentMethodRepository) Update(ctx context.Context, paymentMethod model.PaymentMethod) error {
	err := r.db.WithContext(ctx).Model(&model.PaymentMethod{}).
		Where("id = ? AND deleted_at IS NULL", paymentMethod.ID).
		Updates(&paymentMethod).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PaymentMethodRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&model.PaymentMethod{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}
