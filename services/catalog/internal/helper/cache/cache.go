package cache

import (
	"github.com/usegro/services/catalog/database"
	sharedcache "github.com/usegro/services/shared/pkg/cache"
)

// New returns a Cache using the service's Redis client and app prefix.
func New() *sharedcache.Cache {
	return sharedcache.New(database.GetRedisClient(), database.GetPrefix())
}
