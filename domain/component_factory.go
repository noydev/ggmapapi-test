package domain

import (
	util "github.com/noydev/ggmapapi-test/utils"
)

type ComponentFactory interface {
	CreateMapDomain(cfg util.Config) MapDomain
	CreateCacheClient(cfg util.Config) CacheDomain
}
