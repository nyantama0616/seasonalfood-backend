package model


type Json_Restaurant struct {
	StoreName              string    `json:"store_name"`
	StoreAddress           string    `json:"store_address"`
	StoreGenre             []string  `json:"store_genre"`
	StoreDinnerPriceRange  int       `json:"store_dinner_price_range"`
	StoreLunchPriceRange   int       `json:"store_lunch_price_range"`
	RestaurantLocalPopular float64   `json:"restaurant_local_popular"`
	RestaurantLocalCnt     float64   `json:"restaurant_local_cnt"`
	StoreReviewCnt         int       `json:"store_review_cnt"`
	Score                  int       `json:"score"`
	SeasonalPopularArray   []float64 `json:"seasonal_popular_array"`
	SeasonalCountArray     []float64 `json:"seasonal_count_array"`
	SeasonalShortArray     []float64 `json:"seasonal_short_array"`
	SeasonalFoodName       []string  `json:"seasonal_food_name"`
	ImgURL                 string    `json:"img_url"`
	StoreHomepage          string    `json:"store_homepage"`
	TabelogURL             string    `json:"tabelog_url"`
	LocalFoodName          string    `json:"local_food_name"`
	Longitude              string    `json:"longitude"`
	Latitude               string    `json:"latitude"`
	ID                     int       `json:"id"`
}