package cartgraphqlmapper

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-ecommerce-pb/carts"
)

type CartGraphqlMapper interface {
	ToGraphqlResponseCartDelete(res *pb.ApiResponseCartDelete) *model.APIResponseCartDelete
	ToGraphqlResponseCartAll(res *pb.ApiResponseCartAll) *model.APIResponseCartAll
	ToGraphqlResponseCart(res *pb.ApiResponseCart) *model.APIResponseCart
	ToGraphqlResponsePaginationCart(res *pb.ApiResponsePaginationCart) *model.APIResponsePaginationCart
}
