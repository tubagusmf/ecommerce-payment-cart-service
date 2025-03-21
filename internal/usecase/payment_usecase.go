package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/tubagusmf/ecommerce-payment-cart-service/internal/model"
	pbOrder "github.com/tubagusmf/ecommerce-user-product-service/pb/order"
	pbUser "github.com/tubagusmf/ecommerce-user-product-service/pb/user"
)

type PaymentUsecase struct {
	paymentRepo model.IPaymentRepository
	orderClient pbOrder.OrderServiceClient
	userClient  pbUser.UserServiceClient
}

func NewPaymentUsecase(
	paymentRepo model.IPaymentRepository,
	orderClient pbOrder.OrderServiceClient,
	userClient pbUser.UserServiceClient,
) model.IPaymentUsecase {
	return &PaymentUsecase{
		paymentRepo: paymentRepo,
		orderClient: orderClient,
		userClient:  userClient,
	}
}

func (u *PaymentUsecase) GetPaymentMethodByID(ctx context.Context, methodID int64) (*model.PaymentMethod, error) {
	paymentMethod, err := u.paymentRepo.FindPaymentMethodByID(ctx, methodID)
	if err != nil {
		log.Printf("[ERROR] Payment method not found: %v", err)
		return nil, errors.New("payment method not found")
	}
	return paymentMethod, nil
}

func (u *PaymentUsecase) ProcessPayment(ctx context.Context, orderID string, userID int64, paymentMethod model.PaymentMethod, paymentStatus model.PaymentStatus) (*model.Payment, error) {
	// check Order
	order, err := u.orderClient.GetOrder(ctx, &pbOrder.GetOrderRequest{OrderId: orderID})
	if err != nil || order == nil {
		log.Printf("[ERROR] Invalid order: %v", err)
		return nil, errors.New("invalid order")
	}

	// check User
	user, err := u.userClient.GetUser(ctx, &pbUser.GetUserRequest{UserId: userID})
	if err != nil || user == nil {
		log.Printf("[ERROR] Invalid user: %v", err)
		return nil, errors.New("invalid user")
	}

	if paymentMethod.ID == 0 {
		log.Printf("[ERROR] Invalid PaymentMethod ID: %d", paymentMethod.ID)
		return nil, errors.New("invalid payment method ID")
	}

	payment := &model.Payment{
		OrderID:         orderID,
		UserID:          userID,
		PaymentMethodID: paymentMethod.ID,
		PaymentMethod:   paymentMethod,
		Status:          paymentStatus,
	}

	err = u.paymentRepo.Create(ctx, payment)
	if err != nil {
		log.Printf("[ERROR] Failed to save payment: %v", err)
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	if paymentStatus == model.StatusSuccess {
		_, err := u.orderClient.MarkOrderPaid(ctx, &pbOrder.MarkOrderPaidRequest{OrderId: orderID})
		if err != nil {
			log.Printf("[ERROR] Failed to mark order as paid: %v", err)
			return nil, fmt.Errorf("failed to mark order as paid: %w", err)
		}
		log.Printf("[INFO] Order %s marked as PAID", orderID)
	}

	log.Printf("[INFO] Payment processed successfully for OrderID: %s", orderID)
	return payment, nil
}

func (u *PaymentUsecase) ConfirmPayment(ctx context.Context, orderID string) error {
	payment, err := u.paymentRepo.FindByOrderID(ctx, orderID)
	if err != nil || payment == nil {
		log.Printf("[ERROR] Payment not found for OrderID: %s", orderID)
		return errors.New("payment not found")
	}

	err = u.paymentRepo.UpdateStatus(ctx, orderID, model.StatusSuccess)
	if err != nil {
		log.Printf("[ERROR] Failed to update payment status: %v", err)
		return err
	}

	log.Printf("[INFO] Payment confirmed for OrderID: %s", orderID)
	return nil
}

func (u *PaymentUsecase) GetPaymentStatus(ctx context.Context, paymentID int64) (*model.Payment, error) {
	payment, err := u.paymentRepo.FindById(ctx, paymentID)
	if err != nil {
		log.Printf("[ERROR] Failed to get payment status: %v", err)
		return nil, err
	}
	return payment, nil
}

func (u *PaymentUsecase) GetPayments(ctx context.Context, payment model.Payment) ([]*model.Payment, error) {
	payments, err := u.paymentRepo.FindAll(ctx, payment)
	if err != nil {
		log.Printf("[ERROR] Failed to get payments: %v", err)
		return nil, err
	}
	return payments, nil
}

func (u *PaymentUsecase) GetPaymentByID(ctx context.Context, id int64) (*model.Payment, error) {
	payment, err := u.paymentRepo.FindById(ctx, id)
	if err != nil {
		log.Printf("[ERROR] Failed to get payment: %v", err)
		return nil, err
	}
	return payment, nil
}

func (u *PaymentUsecase) GetPaymentByOrderID(ctx context.Context, orderID string) (*model.Payment, error) {
	payment, err := u.paymentRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		log.Printf("[ERROR] Failed to get payment: %v", err)
		return nil, err
	}
	return payment, nil
}

func (u *PaymentUsecase) MarkPaymentPaid(ctx context.Context, id string) error {
	err := u.paymentRepo.UpdateStatus(ctx, id, model.StatusSuccess)
	if err != nil {
		log.Printf("[ERROR] Failed to mark payment as paid: %v", err)
		return err
	}
	return nil
}
