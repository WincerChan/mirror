package tool

import (
	"strings"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
func HasGziped(coding string) bool {
	return strings.HasPrefix(coding, "gz")
}

func HasBrotli(coding string) bool {
	return strings.HasPrefix(coding, "br")
}

func IsTextType(typeName string) bool {
	return strings.HasPrefix(typeName, "text") ||
		strings.HasPrefix(typeName, "appli")
}
