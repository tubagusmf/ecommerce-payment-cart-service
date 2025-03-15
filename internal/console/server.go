package console

import (
	"log"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/tubagusmf/ecommerce-payment-cart-service/db"
	"github.com/tubagusmf/ecommerce-payment-cart-service/internal/config"
	"github.com/tubagusmf/ecommerce-payment-cart-service/internal/repository"
	"github.com/tubagusmf/ecommerce-payment-cart-service/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpcHandler "github.com/tubagusmf/ecommerce-payment-cart-service/internal/delivery/grpc"
	httpHandler "github.com/tubagusmf/ecommerce-payment-cart-service/internal/delivery/http"

	pbPayment "github.com/tubagusmf/ecommerce-payment-cart-service/pb/payment_service"
	pbOrder "github.com/tubagusmf/ecommerce-user-product-service/pb/order"
	pbUser "github.com/tubagusmf/ecommerce-user-product-service/pb/user"
)

var serverCmd = &cobra.Command{
	Use:   "httpsrv",
	Short: "Run the Payment Service server",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadWithViper()
		postgresDB := db.NewPostgres()
		sqlDB, err := postgresDB.DB()
		if err != nil {
			log.Fatalf("Failed to get SQL DB from Gorm: %v", err)
		}
		defer sqlDB.Close()

		paymentRepo := repository.NewPaymentRepo(postgresDB)
		paymentMethodRepo := repository.NewPaymentMethodRepo(postgresDB)
		orderClient := newOrderClientGRPC()
		userClient := newUserClientGRPC()

		paymentUsecase := usecase.NewPaymentUsecase(paymentRepo, orderClient, userClient)
		paymentMethodUsecase := usecase.NewPaymentMethodUsecase(paymentMethodRepo)

		quitChannel := make(chan bool, 1)

		paymentUsecaseConcrete, ok := paymentUsecase.(*usecase.PaymentUsecase)
		if !ok {
			log.Fatal("Failed to assert paymentUsecase to *usecase.PaymentUsecase")
		}

		paymentMethodUsecaseConcrete, ok := paymentMethodUsecase.(*usecase.PaymentMethodUsecase)
		if !ok {
			log.Fatal("Failed to assert paymentMethodUsecase to *usecase.PaymentMethodUsecase")
		}

		go startHTTPServer(paymentUsecaseConcrete, paymentMethodUsecaseConcrete)
		go startGRPCServer(paymentUsecaseConcrete)

		<-quitChannel
	},
}

func startHTTPServer(paymentUsecase *usecase.PaymentUsecase, paymentMethodUsecase *usecase.PaymentMethodUsecase) {
	e := echo.New()

	httpHandler.NewPaymentMethodHandler(e, paymentMethodUsecase)

	httpHandler.NewPaymentHttpHandler(e, paymentUsecase, paymentMethodUsecase)

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong!")
	})

	log.Println("HTTP server running on port 3200")
	if err := e.Start(":3200"); err != nil && err != http.ErrServerClosed {
		log.Println("HTTP server error:", err)
	}
}

func startGRPCServer(paymentUsecase *usecase.PaymentUsecase) {
	grpcServer := grpc.NewServer()
	paymentgRPCHandler := grpcHandler.NewPaymentgRPCHandler(paymentUsecase)
	pbPayment.RegisterPaymentServiceServer(grpcServer, paymentgRPCHandler)

	listener, err := net.Listen("tcp", ":7000")
	if err != nil {
		log.Fatalf("Failed to create gRPC listener: %v", err)
	}

	log.Println("gRPC server running on port 7000")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}

func newUserClientGRPC() pbUser.UserServiceClient {
	conn, err := grpc.Dial("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to User Service: %v", err)
	}
	return pbUser.NewUserServiceClient(conn)
}

func newOrderClientGRPC() pbOrder.OrderServiceClient {
	conn, err := grpc.Dial("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Order Service: %v", err)
	}
	return pbOrder.NewOrderServiceClient(conn)
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
