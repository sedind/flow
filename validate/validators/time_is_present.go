package validators

import (
	"fmt"
	"time"

	"github.com/gobuffalo/validate"
)

// TimeIsPresent validator
type TimeIsPresent struct {
	Name  string
	Field time.Time
}

// IsValid check if Filed time value is Present
func (v *TimeIsPresent) IsValid(errors *validate.Errors) {
	t := time.Time{}
	if v.Field.UnixNano() == t.UnixNano() {
		errors.Add(GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
	}
}
