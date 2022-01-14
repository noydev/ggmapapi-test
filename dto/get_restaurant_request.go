package dto

import (
	"github.com/noydev/ggmapapi-test/common"
)

type QueryType string

const (
	QueryTypeLocation = "location"
	QueryTypeText     = "text"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type GetRestaurantRequest struct {
	Location     Location  `json:"location"`
	TextLocation string    `json:"text_location"`
	QueryType    QueryType `json:"query_type"`
	Page         int       `json:"page"`
}

func (r GetRestaurantRequest) Validate() *common.AppError {
	if r.Location.Latitude == 0 && r.Location.Longitude == 0 && len(r.TextLocation) == 0 {
		return common.NewValidationError("Empty query")
	}
	if !r.QueryType.isValid() {
		return common.NewValidationError("Invalid Query Type")
	}
	if r.Page < 1 {
		return common.NewValidationError("Invalid Page number")
	}
	return nil
}

func (q QueryType) isValid() bool {
	switch q {
	case QueryTypeLocation, QueryTypeText:
		return true
	}
	return false
}
