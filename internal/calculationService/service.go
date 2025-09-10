package calculationService

import (
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
)

type CalculationService interface {
	CreateCalculation(expression string) (Calculation, error)
	GetCalculations() ([]Calculation, error)
	GetCalculationByID(id string) (Calculation, error)
	UpdateCalculation(id, expression string) (Calculation, error)
	DeleteCalculation(id string) error
}

type calcService struct {
	repo CalculationRepository
}

func validateUuid(id string) bool {
	if !(len(id) == 36) {
		return false
	}
	if !(strings.Count(id, "-") == 4) {
		return false
	}
	return true
}

func NewCalculationService(r CalculationRepository) CalculationService {
	return &calcService{repo: r}
}

func (s *calcService) calculateExpression(expression string) (string, error) {
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

func (s *calcService) CreateCalculation(expression string) (Calculation, error) {
	result, err := s.calculateExpression(expression)
	if err != nil {
		return Calculation{}, err
	}
	calc := Calculation{
		ID:         uuid.NewString(),
		Expression: expression,
		Result:     result,
	}

	if err := s.repo.CreateCalculation(calc); err != nil {
		return Calculation{}, err
	}
	return calc, err
}

func (s *calcService) GetCalculations() ([]Calculation, error) {
	return s.repo.GetCalculations()
}

func (s *calcService) GetCalculationByID(id string) (Calculation, error) {
	if !validateUuid(id) {
		return Calculation{}, fmt.Errorf("invalid id %v", id)
	}
	return s.repo.GetCalculationByID(id)
}

func (s *calcService) UpdateCalculation(id, expression string) (Calculation, error) {
	if !validateUuid(id) {
		return Calculation{}, fmt.Errorf("invalid id %v", id)
	}
	calc, err := s.repo.GetCalculationByID(id)
	if err != nil {
		return Calculation{}, err
	}
	result, err := s.calculateExpression(expression)
	if err != nil {
		return Calculation{}, err
	}

	calc.Expression = expression
	calc.Result = result

	if err := s.repo.UpdateCalculation(calc); err != nil {
		return Calculation{}, err
	}
	return calc, err
}

func (s *calcService) DeleteCalculation(id string) error {
	return s.repo.DeleteCalculation(id)
}
