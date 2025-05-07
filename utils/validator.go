package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func AnalyzeStructError(err error, source string) string {
	var validateErrs validator.ValidationErrors
	if !errors.As(err, &validateErrs) {
		return "Unknown error"
	}

	var errStringBuilder strings.Builder
	for _, e := range validateErrs {
		errStringBuilder.WriteString(fmt.Sprintf("%s, ", e.Error()))
	}

	errString := errStringBuilder.String()
	return errString[:len(errString)-2]
}
