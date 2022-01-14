package app

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/noydev/ggmapapi-test/common"
	"github.com/noydev/ggmapapi-test/dto"
	"github.com/noydev/ggmapapi-test/services"
)

type MapHandler struct {
	service services.IMapService
}

func (h MapHandler) GetRestaurantByLocation(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	text := r.URL.Query().Get("text")
	queryType := r.URL.Query().Get("query_type")
	page := r.URL.Query().Get("page")

	var lat, long float64
	if len(location) > 0 {
		locationSlice := strings.Split(location, ",")

		if s, err := strconv.ParseFloat(locationSlice[0], 32); err == nil {
			lat = s
		}
		if s, err := strconv.ParseFloat(locationSlice[1], 32); err == nil {
			long = s
		}
	}
	if len(page) == 0 {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, &common.AppError{
			Message: "Invalid Page number",
		})
		return
	}

	request := dto.GetRestaurantRequest{
		Location: dto.Location{
			Latitude:  lat,
			Longitude: long,
		},
		TextLocation: text,
		QueryType:    dto.QueryType(queryType),
		Page:         pageInt,
	}
	if appError := request.Validate(); appError != nil {
		writeResponse(w, appError.Code, appError.AsMessage())
		return
	}

	restaurants, appError := h.service.GetRestaurantOnLocation(request)
	if appError != nil {
		writeResponse(w, appError.Code, appError.AsMessage())
	} else {
		writeResponse(w, http.StatusOK, restaurants)
	}
}
