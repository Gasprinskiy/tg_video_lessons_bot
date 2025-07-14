package dump

import "github.com/davecgh/go-spew/spew"

// Struct выводит структуру в подробном виде
func Struct[T any](s T) string {
	return spew.Sdump(s)
}
