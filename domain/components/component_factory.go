package components

import (
	"os"

	util "github.com/noydev/ggmapapi-test/utils"
)

type ComponentFactory struct {
	cfg util.Config
}

func NewComponentFactory(cfg util.Config) ComponentFactory {
	return ComponentFactory{cfg: cfg}
}

func (cf ComponentFactory) CreateMapDomain() GoogleApiClient {
	apiKey := os.Getenv("GOOGLE_MAP_API_KEY")
	return NewGoogleApiClient(apiKey)
}

func (cf ComponentFactory) CreateCacheClient(cfg util.Config) RedisClient {
	redisCfgRaw := cfg.Sub("redis")
	var redisCfg RedisClientConfig
	redisCfgRaw.Unmarshal(&redisCfg)
	redisCfg.Password = os.Getenv("REDIS_PASSWORD")
	return NewRedisClient(redisCfg)
}
