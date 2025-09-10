package handlers

import (
	"backend/internal/calculationService"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CalculationHandler struct {
	service calculationService.CalculationService
}

func NewCalculationHandler(s calculationService.CalculationService) *CalculationHandler {
	return &CalculationHandler{service: s}
}

// GetCalculations godoc
// @Summary Get all calculations
// @Tags calculations
// @Produce json
// @Success 200 {array} Calculation
// @Router /calculations [get]
func (h *CalculationHandler) GetCalculations(c echo.Context) error {
	calculations, err := h.service.GetCalculations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not get calculations"})
	}
	return c.JSON(http.StatusOK, calculations)
}

// PostCalculations godoc
// @Summary Create calculation
// @Tags calculations
// @Accept json
// @Produce json
// @Param input body CalculationRequest true "Expression"
// @Success 201 {object} Calculation
// @Router /calculations [post]
func (h *CalculationHandler) PostCalculations(c echo.Context) error {
	var req calculationService.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	calc, err := h.service.CreateCalculation(req.Expression)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not create calculation"})
	}
	return c.JSON(http.StatusCreated, calc)
}

// PatchCalculations godoc
// @Summary Update calculation
// @Tags calculations
// @Accept json
// @Produce json
// @Param id path string true "Calculation ID"
// @Param input body CalculationRequest true "Expression"
// @Success 200 {object} Calculation
// @Router /calculations/{id} [patch]
func (h *CalculationHandler) PatchCalculations(c echo.Context) error {
	id := c.Param("id")
	var req calculationService.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	updatedCalc, err := h.service.UpdateCalculation(id, req.Expression)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not update calculation"})
	}
	return c.JSON(http.StatusOK, updatedCalc)
}

// DeleteCalculations godoc
// @Summary Delete calculation
// @Tags calculations
// @Param id path string true "Calculation ID"
// @Success 204
// @Router /calculations/{id} [delete]
func (h *CalculationHandler) DeleteCalculations(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.DeleteCalculation(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not delete calculation"})
	}
	return c.NoContent(http.StatusNoContent)
}
