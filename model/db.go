package model

import (
	"log"
	"myapp/migrate"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func init(){
	dsn :="tester:password@tcp(db:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(dsn+"database can't connect")
	}
	// テーブルをすべて削除する
	DB.Migrator().DropTable("genres")
	DB.Migrator().DropTable("restaurants")
	DB.Migrator().DropTable("restaurant_genres")
	DB.Migrator().DropTable("seasonal_data")
	DB.Migrator().DropTable("representative_reviews")
	DB.Migrator().DropTable("seasonal_food_names")
	DB.Migrator().DropTable("restaurant_seasonal_foods")
	// データベースにテーブルがない場合はmigrate.addDataDB()を実行する
	migrate.Add_data_db(DB)
}