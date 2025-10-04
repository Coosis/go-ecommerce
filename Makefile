grpc:
	protoc \
		-I. \
		-I proto \
		--go_out=api-gateway/internal/pb/ \
		--go_out=auth/internal/pb/ \
		--go-grpc_out=api-gateway/internal/pb/ \
		--go-grpc_out=auth/internal/pb/ \
		--grpc-gateway_out=api-gateway/internal/pb/ \
		--grpc-gateway_out=auth/internal/pb/ \
		proto/v1/auth.proto

	protoc \
		-I. \
		-I proto \
		--go_out=api-gateway/internal/pb/ \
		--go_out=catalog/internal/pb/ \
		--go_out=cart/internal/pb/ \
		--go_out=order/internal/pb/ \
		--go-grpc_out=api-gateway/internal/pb/ \
		--go-grpc_out=catalog/internal/pb/ \
		--go-grpc_out=cart/internal/pb/ \
		--go-grpc_out=order/internal/pb/ \
		--grpc-gateway_out=api-gateway/internal/pb/ \
		--grpc-gateway_out=catalog/internal/pb/ \
		--grpc-gateway_out=cart/internal/pb/ \
		--grpc-gateway_out=order/internal/pb/ \
		proto/v1/catalog.proto

	protoc \
		-I. \
		-I proto \
		--go_out=api-gateway/internal/pb/ \
		--go_out=cart/internal/pb/ \
		--go-grpc_out=api-gateway/internal/pb/ \
		--go-grpc_out=cart/internal/pb/ \
		--grpc-gateway_out=api-gateway/internal/pb/ \
		--grpc-gateway_out=cart/internal/pb/ \
		proto/v1/cart.proto

	protoc \
		-I. \
		-I proto \
		--go_out=api-gateway/internal/pb/ \
		--go_out=inventory/internal/pb/ \
		--go_out=order/internal/pb/ \
		--go-grpc_out=api-gateway/internal/pb/ \
		--go-grpc_out=inventory/internal/pb/ \
		--go-grpc_out=order/internal/pb/ \
		--grpc-gateway_out=api-gateway/internal/pb/ \
		--grpc-gateway_out=inventory/internal/pb/ \
		--grpc-gateway_out=order/internal/pb/ \
		proto/v1/inventory.proto

	protoc \
		-I. \
		-I proto \
		--go_out=api-gateway/internal/pb/ \
		--go_out=order/internal/pb/ \
		--go_out=cart/internal/pb/ \
		--go_out=payment/internal/pb/ \
		--go-grpc_out=api-gateway/internal/pb/ \
		--go-grpc_out=order/internal/pb/ \
		--go-grpc_out=cart/internal/pb/ \
		--go-grpc_out=payment/internal/pb/ \
		--grpc-gateway_out=api-gateway/internal/pb/ \
		--grpc-gateway_out=order/internal/pb/ \
		--grpc-gateway_out=cart/internal/pb/ \
		--grpc-gateway_out=payment/internal/pb/ \
		proto/v1/order.proto

	protoc \
		-I. \
		-I proto \
		--go_out=api-gateway/internal/pb/ \
		--go_out=payment/internal/pb/ \
		--go-grpc_out=api-gateway/internal/pb/ \
		--go-grpc_out=payment/internal/pb/ \
		--grpc-gateway_out=api-gateway/internal/pb/ \
		--grpc-gateway_out=payment/internal/pb/ \
		proto/v1/payment.proto
