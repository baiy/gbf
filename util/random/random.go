package random

import "github.com/duke-git/lancet/v2/random"

func String(length int) string {
	return random.RandString(length)
}

func NumeralOrLetter(length int) string {
	return random.RandNumeralOrLetter(length)
}
func Int(min, max int) int {
	return random.RandInt(min, max)
}
