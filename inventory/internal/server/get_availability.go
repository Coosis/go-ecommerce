package server

import (
	"context"

	pb "github.com/Coosis/go-ecommerce/inventory/internal/pb/v1/inventory"
	sqlc "github.com/Coosis/go-ecommerce/inventory/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	log "github.com/sirupsen/logrus"
)

func(s *Server) GetAvailability(
	ctx context.Context,
	req *pb.GetAvailabilityRequest,
) (*pb.GetAvailabilityResponse, error) {
	queries := sqlc.New(s.Pool)
	var availabilities []*pb.Availability
	if req.GetWarehouseId() != "" {
		log.Infof(
			"Getting availability for product %s in warehouse %s",
			req.GetProductId(),
			req.GetWarehouseId(),
		)
		var warehouse_id_uuid pgtype.UUID
		if err := warehouse_id_uuid.Scan(req.GetWarehouseId()); err != nil {
			return nil, err
		}
		var product_id_uuid pgtype.UUID
		if err := product_id_uuid.Scan(req.GetProductId()); err != nil {
			return nil, err
		}
		var sku_ids []pgtype.UUID
		for _, sku_id := range req.GetSkuIds() {
			var sku_id_uuid pgtype.UUID
			if err := sku_id_uuid.Scan(sku_id); err != nil {
				log.Errorf("Error scanning sku_id: %v", err)
				return nil, err
			}
			sku_ids = append(sku_ids, sku_id_uuid)
		}
		resp, err := queries.GetStockLevelsByProductAndSkusAndWarehouse(
			ctx,
			sqlc.GetStockLevelsByProductAndSkusAndWarehouseParams{
				ProductID: product_id_uuid,
				SkuIds: sku_ids,
				WarehouseID: warehouse_id_uuid,
			},
		)
		if err != nil {
			return nil, err
		}
		for _, r := range resp {
			availabilities = append(availabilities, &pb.Availability{
				SkuId: r.SkuID.String(),
				WarehouseId: r.WarehouseID.String(),
				OnHand: r.OnHand,
				Reserved: r.Reserved,
				Available: r.OnHand - r.Reserved,
			})
		}
	} else {
		log.Infof(
			"Getting availability for product %s in all warehouses",
			req.GetProductId(),
		)
		var product_id_uuid pgtype.UUID
		if err := product_id_uuid.Scan(req.GetProductId()); err != nil {
			return nil, err
		}
		var sku_ids []pgtype.UUID
		for _, sku_id := range req.GetSkuIds() {
			var sku_id_uuid pgtype.UUID
			if err := sku_id_uuid.Scan(sku_id); err != nil {
				log.Errorf("Error scanning sku_id: %v", err)
				return nil, err
			}
			sku_ids = append(sku_ids, sku_id_uuid)
		}
		resp, err := queries.GetStockLevelsByProductAndSkus(
			ctx,
			sqlc.GetStockLevelsByProductAndSkusParams{
				ProductID: product_id_uuid,
				SkuIds: sku_ids,
			},
		)
		if err != nil {
			return nil, err
		}
		for _, r := range resp {
			availabilities = append(availabilities, &pb.Availability{
				SkuId: r.SkuID.String(),
				WarehouseId: r.WarehouseID.String(),
				OnHand: r.OnHand,
				Reserved: r.Reserved,
				Available: r.OnHand - r.Reserved,
			})
		}
	}
	return &pb.GetAvailabilityResponse{
		Availabilities: availabilities,
	}, nil
}
