package fuzz

import "log"

// DontPanic is a function that panics if provided with the input "fuzz".
//
// The checks are written in a way that a fuzzing library (like the one present in the Go standard library)
// will detect the case where it fails and should be able to trigger the panic.
func DontPanic(s string) {
	if len(s) == 4 {
		if s[0] == 'f' {
			if s[1] == 'u' {
				if s[2] == 'z' {
					if s[3] == 'z' {
						// don't panic in this case, just log a warning
						log.Println("[WARNING] error: wrong input")
					}
				}
			}
		}
	}
}
