package fuzz

import "testing"

func FuzzDontPanic(f *testing.F) {
	f.Fuzz(func(t *testing.T, input string) {
		DontPanic(input)
	})
}
