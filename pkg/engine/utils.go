package engine

import (
	"os"

	cueerrors "cuelang.org/go/cue/errors"
)

func PrintCueErrorIfNeed(err error) {
	for {
		switch x := err.(type) {
		case cueerrors.Error:
			cueerrors.Print(os.Stderr, x, nil)
			return
		case interface{ Unwrap() error }:
			err = x.Unwrap()
		default:
			return
		}
	}
}
