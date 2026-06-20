package main

import (
	"errors"
	"os"

	"novel-logic/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		code := 1
		var ee *cli.ExitError
		if errors.As(err, &ee) && ee.Code != 0 {
			code = ee.Code
		}
		os.Exit(code)
	}
}