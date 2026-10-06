package cart_cache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/model"
)

type CartQueryCache interface {
	GetCachedCarts(
		ctx context.Context,
		request *model.FindAllCartInput,
	) (*model.APIResponsePaginationCart, bool)

	SetCachedCarts(
		ctx context.Context,
		request *model.FindAllCartInput,
		response *model.APIResponsePaginationCart,
	)
}

type CartCommandCache interface {
	InvalidateCartsCache(ctx context.Context)
}
