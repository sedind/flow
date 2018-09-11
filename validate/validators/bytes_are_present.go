package validators

import (
	"fmt"

	"github.com/sedind/flow/validate"
)

// BytesArePresent validates if Fiels has Byte Array
type BytesArePresent struct {
	Name  string
	Field []byte
}

// IsValid validates if Fiels has Byte Array
func (v *BytesArePresent) IsValid(errors *validate.Errors) {
	if len(v.Field) == 0 {
		errors.Add(GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
	}
}
