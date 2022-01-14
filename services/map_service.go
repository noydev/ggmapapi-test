package services

import (
	"fmt"

	"github.com/noydev/ggmapapi-test/common"
	"github.com/noydev/ggmapapi-test/domain"
	"github.com/noydev/ggmapapi-test/dto"
	"github.com/noydev/ggmapapi-test/utils/logger"
)

//go:generate mockgen -destination=../mocks/services/mockMapService.go -package=services github.com/noydev/ggmapapi-test/services IMapService
type IMapService interface {
	GetRestaurantOnLocation(req dto.GetRestaurantRequest) (*dto.GetRestaurantResponse, *common.AppError)
}

type MapService struct {
	domain *domain.ComponentContainers
}

func NewMapService(domain *domain.ComponentContainers) *MapService {
	return &MapService{domain: domain}
}

func (m MapService) GetRestaurantOnLocation(req dto.GetRestaurantRequest) (*dto.GetRestaurantResponse, *common.AppError) {
	var err error
	location := common.Location{
		Latitude:  req.Location.Latitude,
		Longitude: req.Location.Longitude,
	}

	cacheKey := getKeyFromRestaurantRequest(req)
	var cachedRestaurant common.CacheRestaurantRequest
	err = m.domain.CacheDomain.Get(cacheKey, &cachedRestaurant)
	if err != nil {
		logger.Error(fmt.Sprintf("Something wrong getting cached restaurant %s", err))
	}
	if req.Page <= len(cachedRestaurant.Restaurants) {
		return &dto.GetRestaurantResponse{
			Restaurants: toDtoRestaurantResponse(cachedRestaurant.Restaurants[req.Page-1]),
			CurrentPage: req.Page,
			MaxPage:     len(cachedRestaurant.Restaurants),
		}, nil
	}

	if req.QueryType == dto.QueryTypeText {
		location, err = m.domain.MapDomain.LocationFromText(req.TextLocation)
		if err != nil {
			logger.Error(err.Error())
			return nil, common.NewUnexpectedError("Something wrong with our map api")
		}
	}

	restaurants, nextPageToken, err := m.domain.MapDomain.RestaurantFromLocation(location, cachedRestaurant.NextPageToken)
	if err != nil {
		logger.Error(err.Error())
		return nil, common.NewUnexpectedError("Something wrong with our map api")
	}
	cachedRestaurant.NextPageToken = nextPageToken
	cachedRestaurant.Restaurants = append(cachedRestaurant.Restaurants, restaurants)
	err = m.domain.CacheDomain.Set(cacheKey, cachedRestaurant)
	if err != nil {
		logger.Error(fmt.Sprintf("Error caching restaurant %s", err))
	}

	restaurantsDto := toDtoRestaurantResponse(restaurants)
	return &dto.GetRestaurantResponse{
		Restaurants: restaurantsDto,
		CurrentPage: len(cachedRestaurant.Restaurants),
		MaxPage:     len(cachedRestaurant.Restaurants),
	}, nil
}

func getKeyFromRestaurantRequest(req dto.GetRestaurantRequest) string {
	key := "restaurant"
	if req.QueryType == dto.QueryTypeText {
		key = fmt.Sprintf("%s:%s:%s", key, dto.QueryTypeText, req.TextLocation)
	} else if req.QueryType == dto.QueryTypeLocation {
		key = fmt.Sprintf("%s:%s:%f-%f", key, dto.QueryTypeLocation, req.Location.Latitude, req.Location.Longitude)
	}
	return key
}

func toDtoRestaurantResponse(restaurants []common.RestaurantInfo) []dto.Restaurant {
	var restaurantsDto []dto.Restaurant
	for _, restaurant := range restaurants {
		restaurantsDto = append(restaurantsDto, dto.Restaurant{
			Name:         restaurant.Name,
			StreetAdress: restaurant.Address,
		})
	}
	return restaurantsDto
}
