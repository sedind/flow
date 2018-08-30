package migration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

// Generate writes contents for a given migration in normalized files
func Generate(path, name, ext string, up, down []byte) error {
	n := time.Now().UTC()
	s := n.Format("20060102150405")

	err := os.MkdirAll(path, 0766)
	if err != nil {
		return errors.Wrapf(err, "couldn't create migrations path %s", path)
	}

	upf := filepath.Join(path, (fmt.Sprintf("%s_%s.up.%s", s, name, ext)))
	err = ioutil.WriteFile(upf, up, 0666)
	if err != nil {
		return errors.Wrapf(err, "couldn't write up migration %s", upf)
	}
	fmt.Printf("> %s\n", upf)

	downf := filepath.Join(path, (fmt.Sprintf("%s_%s.down.%s", s, name, ext)))
	err = ioutil.WriteFile(downf, down, 0666)
	if err != nil {
		return errors.Wrapf(err, "couldn't write up migration %s", downf)
	}

	fmt.Printf("> %s\n", downf)
	return nil
}
