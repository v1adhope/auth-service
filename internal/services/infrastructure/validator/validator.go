package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/v1adhope/auth-service/internal/models"
)

type Validator struct {
	*validator.Validate
}

func New() *Validator {
	return &Validator{
		validator.New(validator.WithRequiredStructEnabled()),
	}
}

type guid struct {
	Value string `validate:"uuid"`
}

func (v *Validator) ValidateGuid(target string) error {
	guid := guid{target}

	if err := v.Struct(&guid); err != nil {
		return fmt.Errorf("validator: validator: ValidateGuid: Struct: %w", models.ErrNotValidGuid)
	}

	return nil
}
