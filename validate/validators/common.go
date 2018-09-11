package validators

import (
	"strings"
)

// CustomKeys holds custom validator keys
var CustomKeys = map[string]string{}

// GenerateKey for validator
func GenerateKey(s string) string {
	key := CustomKeys[s]
	if key != "" {
		return key
	}
	key = strings.Replace(s, " ", "", -1)
	key = strings.Replace(key, "-", "", -1)
	return key
}
