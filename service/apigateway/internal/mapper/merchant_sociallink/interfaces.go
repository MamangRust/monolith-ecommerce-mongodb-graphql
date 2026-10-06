package merchant_sociallinkgraphqlmapper

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-ecommerce-pb/merchants_social_link"
)

type MerchantSocialLinkGraphqlMapper interface {
	ToGraphqlResponseMerchantSocialLink(res *pb.ApiResponseMerchantSocial) *model.APIResponseMerchantSocialMediaLink
}
