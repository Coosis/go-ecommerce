package sdk

import (
	"context"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"

	"encoding/json"

	pb "github.com/Coosis/go-ecommerce/payment/internal/pb/v1/payment"
	api "github.com/Coosis/go-ecommerce/payment/internal"
	internal "github.com/Coosis/go-ecommerce/payment/internal"
)

type Stripe struct{}

func NewStripeService() *Stripe {
	stripe.Key = internal.StripeApiKey
	return &Stripe{}
}

type CheckoutData struct {
	ClientSecret string `json:"clientSecret"`
}

func(cd *CheckoutData) String() string {
	b, _ := json.Marshal(cd)
	return string(b)
}

func(str *Stripe) CreatePaymentSession(
	ctx context.Context,
	info *pb.CreatePaymentSessionRequest,
	server api.Api,
) (*pb.CreatePaymentSessionResponse, error) {
	var lineItems []*stripe.CheckoutSessionLineItemParams
	order, err := server.GetOrder(ctx, &api.GetOrderRequest{
		OrderID: info.OrderId,
	})
	if err != nil {
		return nil, err
	}
	for _, item := range order.Items {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("usd"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(item.ProductName),
				},
				UnitAmount: stripe.Int64(int64(item.UnitPriceCents)),
			},
			Quantity: stripe.Int64(int64(item.Qty)),
		})
	}

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		UIMode: stripe.String("embedded"),
		ReturnURL: stripe.String("http://localhost:9765/protected"),
		LineItems: lineItems,
	}

	s, err := session.New(params)
	if err != nil {
		return nil, err
	}

	data := CheckoutData{
		ClientSecret: s.ClientSecret,
	}

	return &pb.CreatePaymentSessionResponse{
		Payment: data.String(),
		Status: "SUCCESS",
	}, nil
}

