package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "backend/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=pw dbname=postgres port=5433 sslmode=disable"
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("couldnt connect to db: %v", err)
	}
	if err := db.AutoMigrate(&Calculation{}); err != nil {
		log.Fatalf("couldnt migrate: %v", err)
	}
}

// Calculation response structure
type Calculation struct {
	ID         string `gorm:"primaryKey" json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

// CalculationRequest request structure
type CalculationRequest struct {
	Expression string `json:"expression"`
}

func calculateExpression(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return "", err
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), nil
}

// getCalculations godoc
// @Summary Get all calculations
// @Tags calculations
// @Produce json
// @Success 200 {array} Calculation
// @Router /calculations [get]
func getCalculations(c echo.Context) error {
	var calclations []Calculation
	if err := db.Find(&calclations).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "couldnt find calculations"})
	}
	return c.JSON(http.StatusOK, calclations)
}

// postCalculations godoc
// @Summary Create calculation
// @Tags calculations
// @Accept json
// @Produce json
// @Param input body CalculationRequest true "Expression"
// @Success 201 {object} Calculation
// @Router /calculations [post]
func postCalculations(c echo.Context) error {
	var req CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	result, err := calculateExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}
	calc := Calculation{
		ID:         uuid.NewString(),
		Expression: req.Expression,
		Result:     result,
	}
	if err := db.Create(&calc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "couldnt add calculation"})
	}
	return c.JSON(http.StatusCreated, calc)
}

// patchCalculations godoc
// @Summary Update calculation
// @Tags calculations
// @Accept json
// @Produce json
// @Param id path string true "Calculation ID"
// @Param input body CalculationRequest true "Expression"
// @Success 200 {object} Calculation
// @Router /calculations/{id} [patch]
func patchCalculations(c echo.Context) error {
	id := c.Param("id")
	var req CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	result, err := calculateExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid expression"})
	}
	var calc Calculation
	if err := db.Find(&calc, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "couldnt find expression"})
	}
	calc.Expression = req.Expression
	calc.Result = result

	if err := db.Save(&calc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "couldnt update calculations"})
	}
	return c.JSON(http.StatusOK, calc)
}

// deleteCalculations godoc
// @Summary Delete calculation
// @Tags calculations
// @Param id path string true "Calculation ID"
// @Success 204
// @Router /calculations/{id} [delete]
func deleteCalculations(c echo.Context) error {
	id := c.Param("id")
	if err := db.Delete(&Calculation{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "couldnt delete calculation"})
	}
	return c.NoContent(http.StatusNoContent)
}

// @title Calculation API
// @version 1.0
// @description Simple calculator API with Echo and Swagger
// @host localhost:8080
// @BasePath /
func main() {
	initDB()
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.GET("/swagger/*", echoSwagger.WrapHandler) // http://localhost:8080/swagger/index.html

	e.GET("/calculations", getCalculations)
	e.POST("/calculations", postCalculations)
	e.PATCH("/calculations/:id", patchCalculations)
	e.DELETE("/calculations/:id", deleteCalculations)

	if err := e.Start("localhost:8080"); err != nil {
		fmt.Println(err)
	}
}
