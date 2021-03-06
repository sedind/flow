package validators

import (
	"fmt"

	"github.com/gobuffalo/validate"
)

// IntIsLessThan validator
type IntIsLessThan struct {
	Name     string
	Field    int
	Compared int
}

// IsValid validates if Field is less than Compared
func (v *IntIsLessThan) IsValid(errors *validate.Errors) {
	if !(v.Field < v.Compared) {
		errors.Add(GenerateKey(v.Name), fmt.Sprintf("%d is not less than %d.", v.Field, v.Compared))
	}
}
