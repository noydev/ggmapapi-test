package services

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/noydev/ggmapapi-test/common"
	realdomain "github.com/noydev/ggmapapi-test/domain"
	"github.com/noydev/ggmapapi-test/dto"
	"github.com/noydev/ggmapapi-test/mocks/domain"
)

var mockMapDomain *domain.MockMapDomain
var service *MapService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockMapDomain = domain.NewMockMapDomain(ctrl)
	cmpnContainer := &realdomain.ComponentContainers{
		MapDomain: mockMapDomain,
	}
	service = NewMapService(cmpnContainer)
	return func() {
		service = nil
		defer ctrl.Finish()
	}
}
func TestGetTransaction(t *testing.T) {
	tests := []struct {
		name                        string
		request                     dto.GetRestaurantRequest
		domainErr                   error
		locationFromText            common.Location
		findRestaurantlocationInput common.Location
		restaurantInfoRes           []common.RestaurantInfo
		nextPageToken               string
		expecterror                 error
	}{
		{
			name: "should_return_customers_with_status_code_200",
			request: dto.GetRestaurantRequest{
				Location: dto.Location{
					Longitude: 0,
					Latitude:  0,
				},
				QueryType:    dto.QueryTypeText,
				TextLocation: "yen lang, hanoi",
			},
			domainErr: nil,
			locationFromText: common.Location{
				Latitude:  21.0227387,
				Longitude: 105.8194541,
			},
			findRestaurantlocationInput: common.Location{
				Latitude:  21.0227387,
				Longitude: 105.8194541,
			},
			nextPageToken: "",
			restaurantInfoRes: []common.RestaurantInfo{
				{
					Name:    "Cafe A",
					Address: "Address A",
				},
				{
					Name:    "Cafe B",
					Address: "Address B",
				},
			},

			expecterror: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			teardown := setup(t)
			defer teardown()

			mockMapDomain.EXPECT().LocationFromText(test.request.TextLocation).Return(test.locationFromText, test.domainErr)
			mockMapDomain.EXPECT().RestaurantFromLocation(test.findRestaurantlocationInput, test.nextPageToken).Return(test.restaurantInfoRes, test.domainErr)
			// Act
			restaurants, appError := service.GetRestaurantOnLocation(test.request)
			// Assert
			if len(restaurants.Restaurants) != len(test.restaurantInfoRes) {
				t.Error("Test failed while creating new account")
			}
			if appError != nil {
				t.Error("Test failed while creating new account")
			}

		})
	}
}
