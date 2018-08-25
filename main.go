package main

import (
	"fmt"

	"github.com/sedind/flow/app"
)

func main() {
	ctx := app.New("config.yaml")
	fmt.Printf("%#v\n", ctx)
}
