package validators

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/validate"
)

// StringInclusion validator
type StringInclusion struct {
	Name  string
	Field string
	List  []string
}

// IsValid performs validation on List array
func (v *StringInclusion) IsValid(errors *validate.Errors) {
	found := false
	for _, l := range v.List {
		if l == v.Field {
			found = true
			break
		}
	}
	if !found {
		errors.Add(GenerateKey(v.Name), fmt.Sprintf("%s is not in the list [%s].", v.Name, strings.Join(v.List, ", ")))
	}
}
