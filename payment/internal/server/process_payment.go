package server

import (
	"context"
	pb "github.com/Coosis/go-ecommerce/payment/internal/pb/v1/payment"
	sdk "github.com/Coosis/go-ecommerce/payment/sdk"
)

func(s *Server) CreatePaymentSession(
	ctx context.Context, 
	req *pb.CreatePaymentSessionRequest,
) (*pb.CreatePaymentSessionResponse, error) {
	var service sdk.PaymentService
	switch req.PaymentSdk {
	case sdk.PAYMENT_SDK_STRIPE:
		service = sdk.NewStripeService()
	case sdk.PAYMENT_SDK_PAYPAL:
		panic("not implemented")
	default:
		panic("unknown payment ")
	}
	return service.CreatePaymentSession(ctx, req, s)
}
