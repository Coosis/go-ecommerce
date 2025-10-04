package sdk

import (
	"context"
	pb "github.com/Coosis/go-ecommerce/payment/internal/pb/v1/payment"
	internal "github.com/Coosis/go-ecommerce/payment/internal"
)

const (
	PAYMENT_SDK_STRIPE = "stripe"
	PAYMENT_SDK_PAYPAL = "paypal"
)

type PaymentInfo struct {
	orderID string
	paymentSDK string // can only come from const PAYMENT_SDK_*
	metadata string
}

type PaymentResult struct {
	payment string
	status string
}

type PaymentService interface {
	CreatePaymentSession(
		ctx context.Context,
		info *pb.CreatePaymentSessionRequest,
		api internal.Api,
	) (*pb.CreatePaymentSessionResponse, error)
}
