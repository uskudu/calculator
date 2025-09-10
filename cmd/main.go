package main

import (
	_ "backend/docs"
	"backend/internal/calculationService"
	"backend/internal/db"
	"backend/internal/handlers"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Calculation API
// @version 1.0
// @description Simple calculator API with Echo and Swagger
// @host localhost:8080
// @BasePath /
func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("couldnt connect to db: %v", err)
	}

	e := echo.New()

	calcRepo := calculationService.NewCalculationRepository(database)
	calcService := calculationService.NewCalculationService(calcRepo)
	calcHandlers := handlers.NewCalculationHandler(calcService)

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.GET("/swagger/*", echoSwagger.WrapHandler) // http://localhost:8080/swagger/index.html

	e.GET("/calculations", calcHandlers.GetCalculations)
	e.POST("/calculations", calcHandlers.PostCalculations)
	e.PATCH("/calculations/:id", calcHandlers.PatchCalculations)
	e.DELETE("/calculations/:id", calcHandlers.DeleteCalculations)

	if err := e.Start("localhost:8080"); err != nil {
		fmt.Println(err)
	}
}
