package components

import (
	"context"
	"fmt"

	"github.com/noydev/ggmapapi-test/common"
	"github.com/noydev/ggmapapi-test/utils/logger"
	"googlemaps.github.io/maps"
)

type GoogleApiClient struct {
	mapClient *maps.Client
}

func NewGoogleApiClient(apiKey string) GoogleApiClient {
	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		logger.Fatal(fmt.Sprintf("fatal error: %s", err))
	}
	return GoogleApiClient{mapClient: c}
}

func (gac GoogleApiClient) LocationFromText(textQuery string) (common.Location, error) {
	r := &maps.TextSearchRequest{
		Query: textQuery,
	}
	locationRes, err := gac.mapClient.TextSearch(context.Background(), r)
	if err != nil {
		return common.Location{}, err
	}

	location := locationRes.Results[0].Geometry.Location
	return common.Location{
		Latitude:  location.Lat,
		Longitude: location.Lng,
	}, nil
}
func (gac GoogleApiClient) RestaurantFromLocation(location common.Location, nextPageToken string) ([]common.RestaurantInfo, string, error) {
	var r maps.NearbySearchRequest
	if len(nextPageToken) > 0 {
		r.PageToken = nextPageToken
	} else {
		r = maps.NearbySearchRequest{
			Location: &maps.LatLng{
				Lat: location.Latitude,
				Lng: location.Longitude,
			},
			Radius: 2000,
			Type:   "restaurant",
		}
	}

	placesRes, err := gac.mapClient.NearbySearch(context.Background(), &r)
	if err != nil {
		return []common.RestaurantInfo{}, "", err
	}
	var restaurantsRes []common.RestaurantInfo
	for _, place := range placesRes.Results {
		restaurant := common.RestaurantInfo{
			Name:    place.Name,
			Address: place.Vicinity,
		}
		restaurantsRes = append(restaurantsRes, restaurant)
	}
	return restaurantsRes, placesRes.NextPageToken, nil
}
