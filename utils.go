package cvault

import (
	"fmt"
	"os"
	"strings"
)

const (
	FileType    = ".cvault"
	TinkHcVault = "hcvault://"
)

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

func SliceMap[T any, V any](arr []T, fn func(T) V) []V {
	ret := make([]V, len(arr))
	for i, it := range arr {
		ret[i] = fn(it)
	}
	return ret
}

func LookupEditor() string {
	if editor, ok := os.LookupEnv("EDITOR"); ok {
		return editor
	}
	if visual, ok := os.LookupEnv("VISUAL"); ok {
		return visual
	}
	return ""
}
