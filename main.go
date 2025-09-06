package main

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"

	_ "backend/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Calculation response structure
type Calculation struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

// CalculationRequest request structure
type CalculationRequest struct {
	Expression string `json:"expression"`
}

var calculations = []Calculation{}

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
	return c.JSON(http.StatusOK, calculations)
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
	calculations = append(calculations, calc)
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	result, err := calculateExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}
	for i, calc := range calculations {
		if calc.ID == id {
			calculations[i].Expression = req.Expression
			calculations[i].Result = result
			return c.JSON(http.StatusOK, calculations[i])
		}
	}
	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Calculation not found"})
}

// deleteCalculations godoc
// @Summary Delete calculation
// @Tags calculations
// @Param id path string true "Calculation ID"
// @Success 204
// @Router /calculations/{id} [delete]
func deleteCalculations(c echo.Context) error {
	id := c.Param("id")
	for i, calc := range calculations {
		if calc.ID == id {
			calculations = append(calculations[:i], calculations[i+1:]...)
			return c.NoContent(http.StatusNoContent)
		}
	}
	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Calculation not found"})
}

// @title Calculation API
// @version 1.0
// @description Simple calculator API with Echo and Swagger
// @host localhost:8080
// @BasePath /
func main() {
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
