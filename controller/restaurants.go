package controller

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"myapp/model"

	"golang.org/x/exp/slices"

	"github.com/labstack/echo/v4"
)

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	var R = 6371 // Earth radius in kilometers
	var φ1 = lat1 * math.Pi / 180 // φ, λ in radians
	var φ2 = lat2 * math.Pi / 180
	var Δφ = (lat2 - lat1) * math.Pi / 180
	var Δλ = (lon2 - lon1) * math.Pi / 180

	var a = math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	var d = float64(R) * c // in kilometers

	return d
}

func GetRestaurant(c echo.Context) error {
	// restaurant_idを取得
	restaurant_id := c.QueryParam("restaurant_id")
	// monthを取得
	month := c.QueryParam("month")

	// restaurant_idが一致するrestaurantを取得
	restaurant := model.Restaurant{}
	if err := model.DB.Where("restaurant_id=?", restaurant_id).First(&restaurant).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Restaurant not found")
	}
	// restaurant_idとmonthが一致するseasonalDataを取得
	seasonalData := model.SeasonalData{}
	if err := model.DB.Where("restaurant_id=? AND month=?", restaurant_id, month).First(&seasonalData).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "SeasonalData not found")
	}
	// restaurant_idが一致するgenreを取得
	restaurant_genres := []model.RestaurantGenre{}
	if err := model.DB.Where("restaurant_id=?", restaurant_id).Find(&restaurant_genres).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "RestaurantGenre not found")
	}
	// restaurant_idが一致するrepresentativeReviewを取得
	representativeReviews := []model.RepresentativeReview{}
	if err := model.DB.Where("restaurant_id=?", restaurant_id).Find(&representativeReviews).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "RepresentativeReview not found")
	}

	// output_restaurantに変換
	output_restaurant := model.Restaurant_output{}
	output_restaurant.StoreName = restaurant.StoreName
	output_restaurant.StoreAddress = restaurant.StoreAddress
	output_restaurant.StoreDinnerPriceRange = restaurant.StoreDinnerPriceRange
	output_restaurant.StoreLunchPriceRange = restaurant.StoreLunchPriceRange
	output_restaurant.RestaurantLocalPopular = restaurant.RestaurantLocalPopular
	output_restaurant.RestaurantLocalCnt = restaurant.RestaurantLocalCnt
	output_restaurant.StoreReviewCnt = restaurant.StoreReviewCnt
	output_restaurant.Score = restaurant.Score
	output_restaurant.ImgURL = restaurant.ImgURL
	output_restaurant.StoreHomepage = restaurant.StoreHomepage
	output_restaurant.TabelogURL = restaurant.TabelogURL
	output_restaurant.LocalFoodName = restaurant.LocalFoodName
	output_restaurant.Longitude = restaurant.Longitude
	output_restaurant.Latitude = restaurant.Latitude
	output_restaurant.RestaurantID = restaurant.RestaurantID
	output_restaurant.RestaurantSeasonalPopular = seasonalData.Popular
	output_restaurant.RestaurantSeasonalCount = seasonalData.Count
	output_restaurant.RestaurantSeasonalShort = seasonalData.Short
	output_restaurant.RestaurantSeasonalFoodname = seasonalData.Foodname
	output_restaurant.NewRestaurantLocalPopularKakariuke = restaurant.NewRestaurantLocalPopularKakariuke
	output_restaurant.NewRestaurantLocalPopularBERT = restaurant.NewRestaurantLocalPopularBERT
	output_restaurant.RestaurantLocalRate = restaurant.RestaurantLocalRate
	output_restaurant.RestaurantZenkokuRate = restaurant.RestaurantZenkokuRate
	output_restaurant.RepresentativeReview = []string{}
	// output_restaurantにStoreGenreを追加
	for _, restaurant_genre := range restaurant_genres {
		genre := model.Genre{}
		if err := model.DB.Where("genre_id=?", restaurant_genre.GenreID).First(&genre).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Genre not found")
		}
		output_restaurant.StoreGenre = append(output_restaurant.StoreGenre, genre.GenreName)
	}

	// output_restaurantにRepresentativeReviewを追加
	// もしrepresentativeReviewがない場合は空配列を返す
	for _, representativeReview := range representativeReviews {
		output_restaurant.RepresentativeReview = append(output_restaurant.RepresentativeReview, representativeReview.Review)
	}

	return c.JSON(http.StatusOK, output_restaurant)
}

func GetRestaurants(c echo.Context) error {

	genreName := c.QueryParam("genre")
	maxbudget := c.QueryParam("maxbudget")
	minbudget := c.QueryParam("minbudget")
	month := c.QueryParam("month")
	time := c.QueryParam("time")
	seasonalFoodName := c.QueryParam("seasonalfoodname")
	radius := c.QueryParam("radius")
	position_latitude := c.QueryParam("position_latitude")
	position_longitude := c.QueryParam("position_longitude")

	genre := model.Genre{}
	if err := model.DB.Where("genre_name=?", genreName).First(&genre).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Genre not found")
	}

	if genre.GenreID == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Genre not found")
	}

	// RestaurantGenreのrestaurant_idを取得
	restaurant_genres := []model.RestaurantGenre{}
	if err := model.DB.Where("genre_id=?", genre.GenreID).Find(&restaurant_genres).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// RestaurantGenreのrestaurant_idを取得
	restaurant_ids := []int{}
	for _, restaurant_genre := range restaurant_genres {
		restaurant_ids = append(restaurant_ids, restaurant_genre.RestaurantID)
	}

	// seasonalFoodNameのfoodidを取得
	seasonal_food_restaurant_ids := []int{}
	seasonalFoodNames := model.SeasonalFoodName{}
	//もし、seasonalFoodNameが空の場合は
	if seasonalFoodName != "" {
		if err := model.DB.Where("food_name=?", seasonalFoodName).First(&seasonalFoodNames).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "SeasonalFoodName not found")
		}

		// RestaurantSeasonalFoodのrestaurant_idを取得
		restaurant_seasonal_foods := []model.RestaurantSeasonalFood{}
		if err := model.DB.Where("food_id=?", seasonalFoodNames.FoodID).Find(&restaurant_seasonal_foods).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		// restaaurant_seasonal_foodsからrestaurant_idを取得
		// 重複を除く
		for _, restaurant_seasonal_food := range restaurant_seasonal_foods {
			seasonal_food_restaurant_ids = append(seasonal_food_restaurant_ids, restaurant_seasonal_food.RestaurantID)
		}

		// sliceを使用して重複を削除
		slices.Sort(seasonal_food_restaurant_ids)
		seasonal_food_restaurant_ids = slices.Compact(seasonal_food_restaurant_ids)
	}

	// 初期クエリ
	query := model.DB.Model(&model.Restaurant{})

	// restaurant_idsのクエリ
	query = query.Where("restaurant_id IN ?", restaurant_ids)
	// timeがlunchの場合,store_lunch_price_rangeのクエリ
	if time == "lunch" {
		query = query.Where("store_lunch_price_range<?", 60000)
		if minbudget != "" {
			query = query.Where("store_lunch_price_range>=? ", minbudget)
		}
		if maxbudget != "" {
			query = query.Where("store_lunch_price_range<=?", maxbudget)
		}
	}
	
	// timeがdinnerの場合,store_dinner_price_rangeのクエリ
	if time == "dinner" {
		query = query.Where("store_dinner_price_range<?", 60000)
		if minbudget != "" {
			query = query.Where("store_dinner_price_range>=?", minbudget)
		}
		if maxbudget != "" {
			query = query.Where("store_dinner_price_range<=?", maxbudget)
		}
	}

	if len(seasonal_food_restaurant_ids) != 0 {
		query = query.Where("restaurant_id IN ?", seasonal_food_restaurant_ids)
	} 
	if len(restaurant_ids) != 0 {
		query = query.Where("restaurant_id IN ?", restaurant_ids)
	}

	// 結果の取得
	restaurants := []model.Restaurant{}
	if err := query.Find(&restaurants).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	seasonalData := []model.SeasonalData{}

	query2 := model.DB.Model(&model.SeasonalData{})
	if month != "" {
		query2 = query2.Where("month=?", month)
	}

	if err := query2.Find(&seasonalData).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	
	
	// model.Restaurants_outputに変換
	Restaurants_output := []model.Restaurants_output{}
	for _, restaurant := range restaurants {
		// restaurant_idが一致するseasonalDataを取得
		restaurant_seasonalData := model.SeasonalData{}
		for _, seasonal_data := range seasonalData {
			if seasonal_data.RestaurantID == restaurant.RestaurantID {
				restaurant_seasonalData = seasonal_data
			}
		}
		// restaurant_seasonalDataをRestaurants_outputに追加
		Restaurants_output = append(Restaurants_output, model.Restaurants_output{
			StoreName:                  restaurant.StoreName,
			StoreAddress:               restaurant.StoreAddress,
			StoreDinnerPriceRange:      restaurant.StoreDinnerPriceRange,
			StoreLunchPriceRange:       restaurant.StoreLunchPriceRange,
			RestaurantLocalPopular:     restaurant.RestaurantLocalPopular,
			RestaurantLocalCnt:         restaurant.RestaurantLocalCnt,
			StoreReviewCnt:             restaurant.StoreReviewCnt,
			Score:                      restaurant.Score,
			ImgURL:                     restaurant.ImgURL,
			StoreHomepage:              restaurant.StoreHomepage,
			TabelogURL:                 restaurant.TabelogURL,
			LocalFoodName:              restaurant.LocalFoodName,
			Longitude:                  restaurant.Longitude,
			Latitude:                   restaurant.Latitude,
			RestaurantID:               restaurant.RestaurantID,
			RestaurantSeasonalPopular:  restaurant_seasonalData.Popular,
			RestaurantSeasonalCount:    restaurant_seasonalData.Count,
			RestaurantSeasonalShort:    restaurant_seasonalData.Short,
			RestaurantSeasonalFoodname: restaurant_seasonalData.Foodname,
			NewRestaurantLocalPopularKakariuke: restaurant.NewRestaurantLocalPopularKakariuke,
			NewRestaurantLocalPopularBERT: restaurant.NewRestaurantLocalPopularBERT,
			RestaurantLocalRate: restaurant.RestaurantLocalRate,
			RestaurantZenkokuRate: restaurant.RestaurantZenkokuRate,
		})

		// restaurant_outputの中でpositionとの距離がradius以内の飲食店を取得
		// undifinedの場合は、radius以内の飲食店を取得しない

		if radius != "" && position_latitude != "" && position_longitude != ""&& radius != "undefined" && position_latitude != "undefined" && position_longitude != "undefined" {
			var tmp_Restaurants_output []model.Restaurants_output

			radius, _ := strconv.ParseFloat(radius, 64)
			positionlatitude, _ := strconv.ParseFloat(position_latitude, 64)
			positionLongitude, _ := strconv.ParseFloat(position_longitude, 64)

			for _, restaurant_output := range Restaurants_output {
				// restaurantの緯度経度を取得
				restaurant_latitude, _ := strconv.ParseFloat(restaurant_output.Latitude, 64)
				restaurant_longitude, _ := strconv.ParseFloat(restaurant_output.Longitude, 64)
				// restaurantとpositionの距離を取得
				distance := haversine(positionlatitude, positionLongitude, restaurant_latitude, restaurant_longitude)
				// radius以内の飲食店を取得
				if distance > 0 && distance < radius{
					tmp_Restaurants_output = append(tmp_Restaurants_output, restaurant_output)
				}
			}
			
			Restaurants_output = tmp_Restaurants_output
		}
	}
	fmt.Printf("%v", Restaurants_output)
	if len(Restaurants_output) == 0 {
		return c.JSON(http.StatusOK, []model.Restaurants_output{})
	}

	return c.JSON(http.StatusOK, Restaurants_output)
}
