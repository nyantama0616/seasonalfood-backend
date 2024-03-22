package main

import (
	"myapp/controller"
	"myapp/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func connect(c echo.Context) error {
	db, _ := model.DB.DB()
	defer db.Close()
	err := db.Ping()
	if err != nil {
		return c.String(http.StatusInternalServerError, "DB接続失敗しました")
	} else {
		return c.String(http.StatusOK, "DB接続しました")
	}
}

func main() {
	e := echo.New()

	// CORSの設定
     e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"http://localhost:8064","http://localhost:3000","http://180.43.174.138:8064", "http://localhost:10000", "https://seasonalfood-front.onrender.com"},
        AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
    }))
	
	e.GET("/", connect)
	e.GET("/restaurants", controller.GetRestaurants)
	e.GET("/restaurant", controller.GetRestaurant)
	e.Logger.Fatal(e.Start(":10000"))
	// e.Logger.Fatal(e.Start(":8023"))
}
