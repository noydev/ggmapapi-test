package domain

import "github.com/noydev/ggmapapi-test/common"

//go:generate mockgen -destination=../mocks/domain/mockMapDomain.go -package=domain github.com/noydev/ggmapapi-test/domain MapDomain
type MapDomain interface {
	LocationFromText(text string) (common.Location, error)
	RestaurantFromLocation(location common.Location, nextPageToken string) ([]common.RestaurantInfo, string, error)
}
