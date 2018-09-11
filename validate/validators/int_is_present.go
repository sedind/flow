package validators

import (
	"fmt"

	"github.com/gobuffalo/validate"
)

// IntIsPresent validator
type IntIsPresent struct {
	Name  string
	Field int
}

// IsValid validates if Field is present -
// Note: Field Value 0 is considered to be blank
func (v *IntIsPresent) IsValid(errors *validate.Errors) {
	if v.Field == 0 {
		errors.Add(GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
	}
}
