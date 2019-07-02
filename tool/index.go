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

func IsTextType(typeName string) bool {
	return strings.HasPrefix(typeName, "text") ||
		strings.HasPrefix(typeName, "appli")
}
