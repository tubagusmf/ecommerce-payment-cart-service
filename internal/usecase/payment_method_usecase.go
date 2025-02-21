package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tubagusmf/ecommerce-payment-cart-service/internal/helper"
	"github.com/tubagusmf/ecommerce-payment-cart-service/internal/model"
)

type PaymentMethodUsecase struct {
	paymentMethodRepo model.IPaymentMethodRepository
}

func NewPaymentMethodUsecase(paymentMethodRepo model.IPaymentMethodRepository) model.IPaymentMethodUsecase {
	return &PaymentMethodUsecase{paymentMethodRepo: paymentMethodRepo}
}

func (u *PaymentMethodUsecase) FindAll(ctx context.Context, paymentMethod model.PaymentMethod) ([]*model.PaymentMethod, error) {
	log := logrus.WithFields(logrus.Fields{
		"paymentMethod": paymentMethod,
	})

	paymentMethods, err := u.paymentMethodRepo.FindAll(ctx, paymentMethod)
	if err != nil {
		log.Error("Failed to get payment methods: ", err)
		return nil, err
	}

	return paymentMethods, nil
}

func (u *PaymentMethodUsecase) FindByID(ctx context.Context, id int64) (*model.PaymentMethod, error) {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	paymentMethod, err := u.paymentMethodRepo.FindByID(ctx, id)
	if err != nil {
		log.Error("Failed to get payment method: ", err)
		return nil, err
	}

	return paymentMethod, nil
}

func (u *PaymentMethodUsecase) Create(ctx context.Context, in model.CreatePaymentMethod) error {
	log := logrus.WithFields(logrus.Fields{
		"in": in,
	})

	err := helper.Validator.Struct(in)
	if err != nil {
		log.Error("Validation error:", err)
		return err
	}

	paymentMethod := model.PaymentMethod{
		Name:      in.Name,
		BankCode:  in.BankCode,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.paymentMethodRepo.Create(ctx, paymentMethod); err != nil {
		log.Error("Failed to create payment method: ", err)
		return err
	}

	return nil
}

func (u *PaymentMethodUsecase) Update(ctx context.Context, id int64, in model.UpdatePaymentMethod) error {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
		"in": in,
	})

	err := helper.Validator.Struct(in)
	if err != nil {
		log.Error("Validation error:", err)
		return err
	}

	paymentMethod, err := u.paymentMethodRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("payment method not found")
	}

	paymentMethod.Name = in.Name
	paymentMethod.BankCode = in.BankCode
	paymentMethod.UpdatedAt = time.Now()

	if err := u.paymentMethodRepo.Update(ctx, *paymentMethod); err != nil {
		log.Error("Failed to update payment method: ", err)
		return err
	}

	return nil
}

func (u *PaymentMethodUsecase) Delete(ctx context.Context, id int64) error {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	paymentMethod, err := u.paymentMethodRepo.FindByID(ctx, id)
	if err != nil {
		log.Error("Failed to find payment method for deletion: ", err)
		return err
	}

	if paymentMethod == nil {
		log.Error("Payment method not found")
		return errors.New("payment method not found")
	}

	if paymentMethod.DeletedAt != nil {
		log.Error("Payment method is already deleted")
		return errors.New("payment method is already deleted")
	}

	if err := u.paymentMethodRepo.Delete(ctx, id); err != nil {
		log.Error("Failed to delete payment method: ", err)
		return err
	}

	log.Info("Successfully deleted payment method with ID: ", id)

	return nil
}
