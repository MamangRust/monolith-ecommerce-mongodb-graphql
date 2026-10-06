package mencache

import "time"

const (
	ttlDefault = 5 * time.Minute

	cacheMerchantKey = "merchant_api_key:%s"
	cacheRoleKey     = "user_roles:%s"

	// CacheKeyPrefix namespaces every apigateway cache entry so the gateway can
	// share a Redis instance with the gRPC services (whose own caches use the
	// same logical keys but store different value shapes) without collisions.
	CacheKeyPrefix = "gql:"
)
