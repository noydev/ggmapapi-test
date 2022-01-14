package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/noydev/ggmapapi-test/common"
	"github.com/noydev/ggmapapi-test/dto"
	"github.com/noydev/ggmapapi-test/mocks/services"
)

var router *mux.Router
var mh MapHandler
var mockService *services.MockIMapService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockService = services.NewMockIMapService(ctrl)
	mh = MapHandler{mockService}
	router = mux.NewRouter()
	router.HandleFunc("/map/restaurant", mh.GetRestaurantByLocation)
	return func() {
		router = nil
		defer ctrl.Finish()
	}
}

func TestGetNearByRestaurant(t *testing.T) {
	tests := []struct {
		name             string
		restaurants      *dto.GetRestaurantResponse
		serviceMapError  *common.AppError
		queryType        dto.QueryType
		textLocation     string
		latitude         float64
		longtitude       float64
		expecterror      error
		expectedHttpCode int
	}{
		{
			name: "should_return_customers_with_status_code_200",
			restaurants: &dto.GetRestaurantResponse{
				Restaurants: []dto.Restaurant{
					{
						Name:         "restaurant1",
						StreetAdress: "street address 1",
					},
				},
			},
			queryType:        dto.QueryTypeLocation,
			textLocation:     "yen lang, hanoi",
			expecterror:      nil,
			expectedHttpCode: http.StatusOK,
		},
		{
			name: "should_return_return_status_code_422_with_error_message",
			restaurants: &dto.GetRestaurantResponse{
				Restaurants: []dto.Restaurant{
					{
						Name:         "restaurant1",
						StreetAdress: "street address 1",
					},
				},
			},
			expecterror:      nil,
			expectedHttpCode: http.StatusUnprocessableEntity,
		},
		{
			name:        "should_return_return_status_code_500_with_error_message",
			restaurants: &dto.GetRestaurantResponse{},
			serviceMapError: &common.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Internal error",
			},
			queryType:        dto.QueryTypeLocation,
			textLocation:     "yen lang, hanoi",
			expecterror:      nil,
			expectedHttpCode: http.StatusInternalServerError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			teardown := setup(t)
			defer teardown()

			expectedMapServiceInput := dto.GetRestaurantRequest{
				Location: dto.Location{
					Longitude: test.longtitude,
					Latitude:  test.latitude,
				},
				QueryType:    test.queryType,
				TextLocation: test.textLocation,
			}

			mockService.EXPECT().GetRestaurantOnLocation(expectedMapServiceInput).Return(test.restaurants, test.serviceMapError)
			req, _ := http.NewRequest(http.MethodGet, "/map/restaurant", nil)
			q := req.URL.Query()
			q.Add("query_type", string(test.queryType))
			q.Add("text", test.textLocation)
			req.URL.RawQuery = q.Encode()

			fmt.Println(req.URL.String())
			// Act
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			// Assert
			if recorder.Code != test.expectedHttpCode {
				t.Error("Failed while testing the status code")
			}
		})
	}
}
