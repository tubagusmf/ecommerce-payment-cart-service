package grpc

import (
	"context"
	"log"
	"strconv"

	"github.com/tubagusmf/ecommerce-payment-cart-service/internal/model"
	pb "github.com/tubagusmf/ecommerce-payment-cart-service/pb/payment_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentgRPCHandler struct {
	pb.UnimplementedPaymentServiceServer
	paymentUsecase model.IPaymentUsecase
}

func NewPaymentgRPCHandler(paymentUsecase model.IPaymentUsecase) *PaymentgRPCHandler {
	return &PaymentgRPCHandler{
		paymentUsecase: paymentUsecase,
	}
}

func (h *PaymentgRPCHandler) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
	log.Println("Processing payment for OrderID:", req.OrderId)

	if req.PaymentMethodId == 0 {
		log.Println("Error: PaymentMethod is missing")
		return nil, status.Errorf(codes.InvalidArgument, "Payment method is required")
	}

	paymentMethod, err := h.paymentUsecase.GetPaymentMethodByID(ctx, req.PaymentMethodId)
	if err != nil {
		log.Println("Error finding payment method:", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid payment method: %v", err)
	}

	var paymentStatus model.PaymentStatus
	switch req.Status {
	case pb.PaymentStatus_PAYMENT_STATUS_PENDING:
		paymentStatus = model.StatusPending
	case pb.PaymentStatus_PAYMENT_STATUS_SUCCESS:
		paymentStatus = model.StatusSuccess
	case pb.PaymentStatus_PAYMENT_STATUS_FAILED:
		paymentStatus = model.StatusFailed
	default:
		log.Println("Invalid payment status in request")
		return nil, status.Errorf(codes.InvalidArgument, "Invalid payment status")
	}

	createdPayment, err := h.paymentUsecase.ProcessPayment(ctx, req.OrderId, req.UserId, *paymentMethod, paymentStatus)
	if err != nil {
		log.Println("Error processing payment:", err)
		return nil, status.Errorf(codes.Internal, "Failed to process payment: %v", err)
	}

	return &pb.ProcessPaymentResponse{
		PaymentId:       strconv.FormatInt(createdPayment.ID, 10),
		OrderId:         createdPayment.OrderID,
		UserId:          createdPayment.UserID,
		PaymentMethodId: createdPayment.PaymentMethod.ID,
		Status:          model.ModelToProtoPaymentStatus(createdPayment.Status),
		TransactionId:   createdPayment.TransactionID,
	}, nil
}

func (h *PaymentgRPCHandler) GetPaymentStatus(ctx context.Context, req *pb.GetPaymentStatusRequest) (*pb.GetPaymentStatusResponse, error) {
	log.Println("Fetching payment status for PaymentID:", req.PaymentId)

	paymentID, err := strconv.ParseInt(req.PaymentId, 10, 64)
	if err != nil {
		log.Println("Error converting PaymentID to int64:", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid payment ID format")
	}

	payment, err := h.paymentUsecase.GetPaymentStatus(ctx, paymentID)
	if err != nil {
		log.Println("Error fetching payment status:", err)
		return nil, status.Errorf(codes.Internal, "Failed to get payment status: %v", err)
	}

	protoPaymentMethod := &pb.PaymentMethod{
		PaymentMethodId: payment.PaymentMethod.ID,
		Name:            payment.PaymentMethod.Name,
		BankCode:        payment.PaymentMethod.BankCode,
	}

	return &pb.GetPaymentStatusResponse{
		PaymentId:     strconv.FormatInt(payment.ID, 10),
		OrderId:       payment.OrderID,
		UserId:        payment.UserID,
		PaymentMethod: protoPaymentMethod,
		Status:        model.ModelToProtoPaymentStatus(payment.Status),
		TransactionId: payment.TransactionID,
	}, nil
}
