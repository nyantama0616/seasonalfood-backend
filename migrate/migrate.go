package migrate

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// jsonの構造体
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
	NewRestaurantLocalPopularKakariuke int `json:"new_restaurant_local_popular_kakariuke"`
	NewRestaurantLocalPopularBERT int `json:"new_restaurant_local_popular_bert"`
	RestaurantLocalRate int `json:"restaurant_local_rate"`
	RestaurantZenkokuRate int `json:"restaurant_zenkoku_rate"`
}

// output用の構造体
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

type RepresentativeReview struct {
	RestaurantID int
	Review       string
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

func removeDuplicate(genre_data []Genre) []Genre {
	// genre_dataの重複を削除する
	var output_genres []Genre
	for _, genre := range genre_data {
		if len(output_genres) == 0 {
			output_genres = append(output_genres, genre)
		} else {
			flag := 0
			for _, output_genre := range output_genres {
				if genre.GenreName == output_genre.GenreName {
					flag = 1
				}
			}
			if flag == 0 {
				output_genres = append(output_genres, genre)
			}
		}

	}
	return output_genres
}

func Add_data_db(db *gorm.DB) {
	// ファイルを読み込む
	log.Println("migrate start")
	file, err := os.Open("new_restaurant_info.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := os.ReadFile("new_restaurant_info.json")
	if err != nil {
		panic(err)
	}

	var restaurants []Json_Restaurant
	if err := json.Unmarshal(b, &restaurants); err != nil {
		panic(err)
	}

	// output用の構造体に変換
	var output_restaurants []Restaurant
	for _, restaurant := range restaurants {
		output_restaurants = append(output_restaurants, Restaurant{
			StoreName:              restaurant.StoreName,
			StoreAddress:           restaurant.StoreAddress,
			StoreDinnerPriceRange:  restaurant.StoreDinnerPriceRange,
			StoreLunchPriceRange:   restaurant.StoreLunchPriceRange,
			RestaurantLocalPopular: restaurant.RestaurantLocalPopular,
			RestaurantLocalCnt:     restaurant.RestaurantLocalCnt,
			StoreReviewCnt:         restaurant.StoreReviewCnt,
			Score:                  restaurant.Score,
			ImgURL:                 restaurant.ImgURL,
			StoreHomepage:          restaurant.StoreHomepage,
			TabelogURL:             restaurant.TabelogURL,
			LocalFoodName:          restaurant.LocalFoodName,
			Longitude:              restaurant.Longitude,
			Latitude:               restaurant.Latitude,
			RestaurantID:           restaurant.ID,
			NewRestaurantLocalPopularKakariuke: restaurant.NewRestaurantLocalPopularKakariuke,
			NewRestaurantLocalPopularBERT: restaurant.NewRestaurantLocalPopularBERT,
			RestaurantLocalRate: restaurant.RestaurantLocalRate,
			RestaurantZenkokuRate: restaurant.RestaurantZenkokuRate,
		})
	}

	// output_genresにユニークなジャンルを追加
	var output_genres []Genre

	for _, restaurant := range restaurants {
		for _, genre := range restaurant.StoreGenre {
			if genre != "" {
				output_genres = append(output_genres, Genre{
					GenreName: genre,
				})
			}
		}
	}
	// output_genreの重複を削除する
	output_genres = removeDuplicate(output_genres)
	// GenreIDを追加
	for i:= range output_genres {
		output_genres[i].GenreID = i + 1
	}
	

	// output_restaurant_genresにrestaurant_idとgenre_idを追加
	var output_restaurant_genres []RestaurantGenre
	for _, restaurant := range restaurants {
		for _, genre := range restaurant.StoreGenre {
			if genre != "" {
				for _, output_genre := range output_genres {
					if genre == output_genre.GenreName {
						output_restaurant_genres = append(output_restaurant_genres, RestaurantGenre{
							GenreID: 	int(output_genre.GenreID),
							RestaurantID: 	int(output_restaurants[restaurant.ID-1].RestaurantID),
						})
					}
				}
			}
		}
	
	}

	// output_seasonal_dataにデータを追加
	var output_seasonal_data []SeasonalData
	for _, restaurant := range restaurants {
		for i := range restaurant.SeasonalPopularArray {
			output_seasonal_data = append(output_seasonal_data, SeasonalData{
				RestaurantID:  int(output_restaurants[restaurant.ID-1].RestaurantID),
				Month:         i + 1,
				Popular:       restaurant.SeasonalPopularArray[i],
				Count:         restaurant.SeasonalCountArray[i],
				Short:         restaurant.SeasonalShortArray[i],
				Foodname:      restaurant.SeasonalFoodName[i],
			})
		}
	}

	var output_food_data []SeasonalFoodName
	var output_restaurant_seasonal_food []RestaurantSeasonalFood
	var food_id = 1

	for _, seasonal_data :=  range output_seasonal_data {
		var foodstring = seasonal_data.Foodname
		if foodstring != "" {
			// 末尾の空白を削除
			foodstring = strings.TrimRight(foodstring, " ")
			// 空白で区切る
			foodarray := strings.Split(foodstring, "　")
			for _, food := range foodarray {
				// もし、output_food_dataのfood_nameにfoodがない場合は追加する
				flag := 0
				for _, food_data := range output_food_data {
					if food_data.FoodName == food {
						flag = 1
					}
				}
				if flag == 0 {
					output_food_data = append(output_food_data, SeasonalFoodName{
						FoodID: food_id,
						FoodName: food,
					})
					food_id += 1
				}
				// output_restaurant_seasonal_foodにrestaurant_id, food_id, monthを追加
				for _, food_data := range output_food_data {
					if food_data.FoodName == food {
						output_restaurant_seasonal_food = append(output_restaurant_seasonal_food, RestaurantSeasonalFood{
							RestaurantID: seasonal_data.RestaurantID,
							FoodID: food_data.FoodID,
							Month: seasonal_data.Month,
						})
					}
				}
			}
		}
	}

	// csvファイルを読み込む
	csv_file, err := os.Open("representative_revew.csv")
	if err != nil {
    log.Fatal(err)
  }
  defer csv_file.Close()
  representative_review := []RepresentativeReview{}
  csv_reader := csv.NewReader(csv_file)
  for {
	record, err := csv_reader.Read()
	if err != nil {
		break
	}
	// record[0]をintに変換
	var rid, _ = strconv.Atoi(record[1])
	representative_review = append(representative_review, RepresentativeReview{
		Review: record[0],
		RestaurantID: rid,
	})
	  }
	// テーブルを削除
	if err := db.Migrator().DropTable(&Restaurant{}); err != nil {
		log.Println("restaurant table drop error")
		log.Println(err)
	} else {
		log.Println("restaurant table drop success")
	}

	// テーブルを作成
	if err := db.AutoMigrate(&Restaurant{}); err != nil {
		log.Println("restaurant table create error")
		log.Println(err)
	} else {
		log.Println("restaurant table create success")
	}

	if err := db.AutoMigrate(&Genre{}); err != nil {
		log.Println("genre table create error")
		log.Println(err)
	}else {
		log.Println("genre table create success")
	}

	if err := db.AutoMigrate(&RestaurantGenre{}); err != nil {
		log.Println("restaurantgenre table create error")
		log.Println(err)
	} else {
		log.Println("restaurantgenre table create success")
	}

	if err := db.AutoMigrate(&SeasonalData{}); err != nil {
		log.Println("SeasonalData table create error")
		log.Println(err)
	} else {
		log.Println("SeasonalData table create success")
	}

	if err := db.AutoMigrate(&RepresentativeReview{}); err != nil {
		log.Println("RepresentativeReview table create error")
		log.Println(err)
	} else {
		log.Println("RepresentativeReview table create success")
	}
	if err := db.AutoMigrate(&SeasonalFoodName{}); err != nil {
		log.Println("SeasonalFoodName table create error")
		log.Println(err)
	} else {
		log.Println("SeasonalFoodName table create success")
	}
	if err := db.AutoMigrate(&RestaurantSeasonalFood{}); err != nil {
		log.Println("RestaurantSeasonalFood table create error")
		log.Println(err)
	} else {
		log.Println("RestaurantSeasonalFood table create success")
	}

	// データベースに追加
	if err := db.Create(&output_restaurants).Error; err != nil {
		log.Println(" restaurant data add error")
		log.Println(err)
	}else {
		log.Println("restaurant data add success")
	}
	if err := db.Create(&output_genres).Error; err != nil {
		log.Println("genre data add error")
		log.Println(err)
	}else {
		log.Println("genre data add success")
	}
	if err := db.Create(&output_restaurant_genres).Error; err != nil {
		log.Println("restaurantgenre data add error")
		log.Println(err)
	}else {
		log.Println("restaurantgenre data add success")
	}
	if err := db.Create(&representative_review).Error; err != nil {
		log.Println("representative_review data add error")
		log.Println(err)
	}else {
		log.Println("representative_review data add success")
	}
	if err := db.Create(&output_food_data).Error; err != nil {
		log.Println("food data add error")
		log.Println(err)
	}else {
		log.Println("food data add success")
	}

	// クエリにおいて許容されるプレースホルダーの最大数を超えるためバッチインサートの分割をする
	// 1000個ずつに分割
	var output_seasonal_data_1000 [][]SeasonalData
	for i := 0; i < len(output_seasonal_data); i += 1000 {
		end := i + 1000
		if end > len(output_seasonal_data) {
			end = len(output_seasonal_data)
		}
		output_seasonal_data_1000 = append(output_seasonal_data_1000, output_seasonal_data[i:end])
	}

	// バッチインサート
	for _, output_seasonal_data_1000 := range output_seasonal_data_1000 {
		if err := db.Create(&output_seasonal_data_1000).Error; err != nil {
			log.Println("seasonaldata data add error")
			log.Println(err)
		}else {
			log.Println("seasonaldata data add success")
		}
	}

	// output_restaurant_seasonal_foodが1000個以上の場合は分割する
	var output_restaurant_seasonal_food_1000 [][]RestaurantSeasonalFood
	for i := 0; i < len(output_restaurant_seasonal_food); i += 1000 {
		end := i + 1000
		if end > len(output_restaurant_seasonal_food) {
			end = len(output_restaurant_seasonal_food)
		}
		output_restaurant_seasonal_food_1000 = append(output_restaurant_seasonal_food_1000, output_restaurant_seasonal_food[i:end])
	}

	// バッチインサート
	for _, output_restaurant_seasonal_food_1000 := range output_restaurant_seasonal_food_1000 {
		if err := db.Create(&output_restaurant_seasonal_food_1000).Error; err != nil {
			log.Println("restaurantseasonalfood data add error")
			log.Println(err)
		}else {
			log.Println("restaurantseasonalfood data add success")
		}
	}



	log.Println("migrate end")

}
