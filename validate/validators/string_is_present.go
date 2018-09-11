package validators

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/validate"
)

// StringIsPresent validator
type StringIsPresent struct {
	Name  string
	Field string
}

// IsValid checks if Field is not empty string
func (v *StringIsPresent) IsValid(errors *validate.Errors) {
	if strings.TrimSpace(v.Field) == "" {
		errors.Add(GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
	}
}
