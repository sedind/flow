package dbe

import (
	"fmt"
	"strings"
)

// Column represents model column used for queries
type Column struct {
	Name string
}

// UpdateString prepares column statement for update statement
func (c Column) UpdateString() string {
	vName := c.Name
	if strings.Contains(vName, ".") {
		tmp := strings.Split(vName, ".")
		vName = tmp[1]
	}
	return fmt.Sprintf("%s = :%s", c.Name, vName)
}
