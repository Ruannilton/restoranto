package dependencies

import "time"

type IDistributedCacheRepository interface {
	Set(key string, value string, duration time.Duration) error
	Get(key string) (string, error)
}
