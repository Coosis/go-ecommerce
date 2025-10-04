package server

import (
	"context"
	"fmt"

	pb "github.com/Coosis/go-ecommerce/inventory/internal/pb/v1/inventory"
	sqlc "github.com/Coosis/go-ecommerce/inventory/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	log "github.com/sirupsen/logrus"
)

func(s *Server) ReserveStock(
	ctx context.Context, 
	req *pb.ReserveRequest,
) (*pb.ReserveResponse, error) {
	queries := sqlc.New(s.Pool)

	var warehouse_uuid pgtype.UUID
	if err := warehouse_uuid.Scan(req.GetWarehouseId()); err != nil {
		return nil, err
	}
	var order_uuid pgtype.UUID
	if err := order_uuid.Scan(req.GetOrderId()); err != nil {
		return nil, err
	}
	interval := pgtype.Interval{
		Microseconds: int64(0),
		Days:         int32(1),
		Months:       int32(0),
		Valid:        true,
	}

	var reservations []string
	var failed []*pb.FailedReserveRequest
	for _, item := range req.GetItems() {
		var sku_uuid pgtype.UUID
		if err := sku_uuid.Scan(item.GetSkuId()); err != nil {
			return nil, err
		}
		reservation, err := queries.ReserveStock(ctx, sqlc.ReserveStockParams{
			SkuID: sku_uuid,
			WarehouseID: warehouse_uuid,
			OrderID: order_uuid,
			Qty: item.GetQty(),
			ExpDuration: interval,
		})
		if err != nil {
			log.Errorf("Error reserving stock for sku %s: %v", item.GetSkuId(), err)
			failed = append(failed, &pb.FailedReserveRequest{
				SkuId: item.GetSkuId(),
				WarehouseId: req.GetWarehouseId(),
				OrderId: req.GetOrderId(),
				Reason: fmt.Errorf("error reserving stock: %v", err).Error(),
			})
			continue
		}
		reservations = append(reservations, reservation.ID.String())
	}
	return &pb.ReserveResponse{
		ReservationId: reservations,
		Failed: failed,
	}, nil
}
