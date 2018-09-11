package validators

import (
	"fmt"

	"github.com/sedind/flow/validate"
)

// IntArrayIsPresent validator
type IntArrayIsPresent struct {
	Name  string
	Field []int
}

// IsValid validates if in array is present
func (v *IntArrayIsPresent) IsValid(errors *validate.Errors) {
	if len(v.Field) == 0 {
		errors.Add(GenerateKey(v.Name), fmt.Sprintf("%s can not be empty.", v.Name))
	}
}
