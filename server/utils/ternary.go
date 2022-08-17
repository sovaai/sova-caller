package utils

func If[T any](condition bool, vtrue T, vfalse T) T {
	if condition {
		return vtrue
	}

	return vfalse
}
