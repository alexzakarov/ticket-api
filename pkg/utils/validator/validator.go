package validator

import (
	"context"
	"github.com/go-playground/validator/v10"
)

// Use a single instance of Validate, it caches struct info
var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

// Validate struct fields
func ValidateStruct(ctx context.Context, s interface{}) error {
	return Validate.StructCtx(ctx, s)
}
