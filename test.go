package main

import (
	errors2 "errors"
	"fmt"

	"github.com/tech-nimble/go-tools/helpers/errors"
)

func main() {
	berr := errors2.New("base")
	err := errors.Wrap(berr, "Test")
	errCode := errors.AddErrorCode(err, 100)
	fmt.Println(errCode.Error())

	if errors2.Is(errCode, berr) {
		fmt.Println("Done")
	}
}
