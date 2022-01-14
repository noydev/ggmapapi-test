package dto

type Restaurant struct {
	Name         string `json:"name"`
	StreetAdress string `json:"street_adress"`
}

type GetRestaurantResponse struct {
	Restaurants []Restaurant `json:"restaurants"`
	CurrentPage int          `json:"current_page"`
	MaxPage     int          `json:"max_page"`
}
