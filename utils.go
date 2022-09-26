package cvault

import (
	"fmt"
	"strings"
)

const FileType = ".cvault"

func ToEncryptedName(name string) string {
	return fmt.Sprintf("%s%s", name, FileType)
}

func ToDecryptedName(name string) string {
	return strings.TrimSuffix(name, FileType)
}

func IsEncryptedName(name string) bool {
	return strings.HasSuffix(name, FileType)
}

func SliceFilter[T any](arr []T, fn func(T) bool) []T {
	ret := []T{}
	for _, it := range arr {
		if fn(it) {
			ret = append(ret, it)
		}
	}
	return ret
}
