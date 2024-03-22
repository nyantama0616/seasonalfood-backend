package model

import (
	"log"
	"myapp/migrate"

	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"fmt"
	"github.com/joho/godotenv"
)

var DB *gorm.DB
var err error

func init(){
    _, err := os.Stat(".env")
	if err == nil {
		// .envファイルが存在する場合、読み込む
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else if !os.IsNotExist(err) {
		// .envファイルが存在しない以外のエラーが発生した場合、エラーをログに出力
		log.Fatal(err)
	}
	
	host := os.Getenv("DATABASE_HOST_NAME")
    user := os.Getenv("DATABASE_USER_NAME")
    password := os.Getenv("DATABASE_PASSWORD")
    dbname := os.Getenv("DATABASE_NAME")
    port := os.Getenv("DATABASE_PORT")

	// dsn := "host=postgresql user=root password=password dbname=test port=5432 sslmode=disable"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	for i :=0; i<5; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database. Attempt %d/5. Retrying in 5 seconds...", i+1)
		time.Sleep(5 * time.Second)
	}

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
