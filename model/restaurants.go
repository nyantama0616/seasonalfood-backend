package model

type Restaurant struct {
	StoreName              string `gorm:"primaryKey"`
	StoreAddress           string
	StoreDinnerPriceRange  int
	StoreLunchPriceRange   int
	RestaurantLocalPopular float64
	RestaurantLocalCnt     float64
	StoreReviewCnt         int
	Score                  int
	ImgURL                 string
	StoreHomepage          string
	TabelogURL             string
	LocalFoodName          string
	Longitude              string
	Latitude               string
	RestaurantID           int 
	NewRestaurantLocalPopularKakariuke int 
	NewRestaurantLocalPopularBERT int 
	RestaurantLocalRate int 
	RestaurantZenkokuRate int 
}

type Genre struct {
	GenreID  int  
	GenreName string `gorm:"primaryKey"`
}

type RestaurantGenre struct {
	RestaurantID int
	GenreID      int
}

type SeasonalData struct {
	RestaurantID int
	Month         int
	Popular       float64
	Count         float64
	Short         float64
	Foodname     string
}

type SeasonalFoodName struct {
	FoodID int
	FoodName string
}

type RestaurantSeasonalFood struct {
	RestaurantID int
	FoodID int
	Month int
}

type Restaurants_output struct {
	StoreName              string `gorm:"primaryKey"`
	StoreAddress           string
	StoreDinnerPriceRange  int
	StoreLunchPriceRange   int
	RestaurantLocalPopular float64
	RestaurantLocalCnt     float64
	StoreReviewCnt         int
	Score                  int
	ImgURL                 string
	StoreHomepage          string
	TabelogURL             string
	LocalFoodName          string
	Longitude              string
	Latitude               string
	RestaurantID           int
	RestaurantSeasonalPopular       float64
	RestaurantSeasonalCount         float64
	RestaurantSeasonalShort         float64
	RestaurantSeasonalFoodname     string
	NewRestaurantLocalPopularKakariuke int 
	NewRestaurantLocalPopularBERT int 
	RestaurantLocalRate int 
	RestaurantZenkokuRate int 
}

type RepresentativeReview struct {
	RestaurantID int
	Review       string
}

type Restaurant_output struct {
	StoreName              string `gorm:"primaryKey"`
	StoreAddress           string
	StoreDinnerPriceRange  int
	StoreLunchPriceRange   int
	StoreGenre			 []string
	RestaurantLocalPopular float64
	RestaurantLocalCnt     float64
	StoreReviewCnt         int
	Score                  int
	ImgURL                 string
	StoreHomepage          string
	TabelogURL             string
	LocalFoodName          string
	Longitude              string
	Latitude               string
	RestaurantID           int
	RestaurantSeasonalPopular       float64
	RestaurantSeasonalCount         float64
	RestaurantSeasonalShort         float64
	RestaurantSeasonalFoodname     string
	RepresentativeReview []string
	NewRestaurantLocalPopularKakariuke int 
	NewRestaurantLocalPopularBERT int 
	RestaurantLocalRate int 
	RestaurantZenkokuRate int 
}