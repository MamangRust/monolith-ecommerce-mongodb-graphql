package merchant_policygraphqlmapper

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/model"
	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	pb "github.com/MamangRust/monolith-ecommerce-pb/merchants_policy"
)

type MerchantPolicyGraphqlMapper interface {
	ToGraphqlResponseMerchantPolicyDelete(res *pbmerchant.ApiResponseMerchantDelete) *model.APIResponseMerchantPolicyDelete
	ToGraphqlResponseMerchantPolicyAll(res *pbmerchant.ApiResponseMerchantAll) *model.APIResponseMerchantPolicyAll
	ToGraphqlResponseMerchantPolicy(res *pb.ApiResponseMerchantPolicies) *model.APIResponseMerchantPolicy
	ToGraphqlResponseMerchantPolicyDeleteAt(res *pb.ApiResponseMerchantPoliciesDeleteAt) *model.APIResponseMerchantPolicyDeleteAt
	ToGraphqlResponsesMerchantPolicy(res *pb.ApiResponsesMerchantPolicies) *model.APIResponsesMerchantPolicy
	ToGraphqlResponsePaginationMerchantPolicyDeleteAt(res *pb.ApiResponsePaginationMerchantPoliciesDeleteAt) *model.APIResponsePaginationMerchantPolicyDeleteAt
	ToGraphqlResponsePaginationMerchantPolicy(res *pb.ApiResponsePaginationMerchantPolicies) *model.APIResponsePaginationMerchantPolicy
}
