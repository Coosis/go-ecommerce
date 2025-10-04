package server

import (
	"context"
	api "github.com/Coosis/go-ecommerce/payment/internal"
	pb "github.com/Coosis/go-ecommerce/payment/internal/pb/v1/payment"
	pborder "github.com/Coosis/go-ecommerce/payment/internal/pb/v1/order"
)

type Server struct{
    OrderClient pborder.OrderServiceClient
	
	pb.UnimplementedPaymentServiceServer
}

func(s *Server) GetOrder(
	ctx context.Context, 
	req *api.GetOrderRequest,
) (*api.GetOrderResponse, error) {
	resp, err := s.OrderClient.GetOrder(ctx, &pborder.GetOrderRequest{
		OrderId: req.OrderID,
	})
	if err != nil {
		return nil, err
	}

	var items []api.OrderItems
	for _, item := range resp.Items {
		name := ""
		if item.ProductName != nil {
			name = *item.ProductName
		}
		sku_name := ""
		if item.SkuName != nil {
			sku_name = *item.SkuName
		}
		items = append(items, api.OrderItems{
			ProductID: item.ProductId,
			ProductName: name,
			SkuID: item.SkuId,
			SkuName: sku_name,
			Qty: item.Qty,
			UnitPriceCents: *item.PriceCentsSnapshot,
		})
	}

	return &api.GetOrderResponse{
		OrderID: resp.OrderId,
		Items: items,
	}, nil
}
