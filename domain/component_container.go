package domain

import (
	"github.com/noydev/ggmapapi-test/domain/components"
	util "github.com/noydev/ggmapapi-test/utils"
)

type ComponentContainers struct {
	MapDomain   MapDomain
	CacheDomain CacheDomain
}

func InitDomains(cfg util.Config) *ComponentContainers {
	cf := components.NewComponentFactory(cfg)
	return &ComponentContainers{
		MapDomain:   cf.CreateMapDomain(),
		CacheDomain: cf.CreateCacheClient(cfg),
	}
}
