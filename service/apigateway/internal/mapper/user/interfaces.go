package usergraphqlmapper

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-ecommerce-pb/users"
)

type UserGraphqlMapper interface {
	ToGraphqlResponseUser(res *pb.ApiResponseUser) *model.APIResponseUserResponse
	ToGraphqlResponseUserDeleteAt(res *pb.ApiResponseUserDeleteAt) *model.APIResponseUserResponseDeleteAt
	ToGraphqlResponseUsers(res *pb.ApiResponsesUser) *model.APIResponsesUser
	ToGraphqlResponseUserDelete(res *pb.ApiResponseUserDelete) *model.APIResponseUserDelete
	ToGraphqlResponseUserAll(res *pb.ApiResponseUserAll) *model.APIResponseUserAll
	ToGraphqlResponsePaginationUser(res *pb.ApiResponsePaginationUser) *model.APIResponsePaginationUser
	ToGraphqlResponsePaginationUserDeleteAt(res *pb.ApiResponsePaginationUserDeleteAt) *model.APIResponsePaginationUserDeleteAt
}
