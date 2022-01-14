package common

type RestaurantInfo struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type CacheRestaurantRequest struct {
	Restaurants   [][]RestaurantInfo `json:"restaurants"`
	NextPageToken string             `json:"next_page_token"`
}
